package xcache

import (
	"context"
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
	"time"
)

type (
	basis struct {
		data      map[string]node
		loadGroup *singleflight.Group
		ctx       context.Context
	}

	node struct {
		data      []byte
		expire    time.Duration
		creatTime time.Time
	}

	do func() (interface{}, error)
)

func (b *basis) GetContext() context.Context {
	if b.ctx == nil {
		b.ctx = context.Background()
	}
	return b.ctx
}

func (b *basis) WithContext(ctx context.Context) Cache {
	b.ctx = ctx
	return b
}

func NewBasis() *basis {
	return &basis{
		data:      make(map[string]node),
		loadGroup: &singleflight.Group{},
	}
}

func (b *basis) Del(keys ...string) error {
	_, err := done(b.GetContext(), func() (interface{}, error) {
		for _, key := range keys {
			b.doDel(key)
		}
		return nil, nil
	})
	return err
}

func done(ctx context.Context, df do) (interface{}, error) {
	var (
		c   = make(chan struct{})
		ret interface{}
		err error
	)
	go func() {
		ret, err = df()
		close(c)
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c:
	}
	return ret, err
}

func (b basis) GetE(key string) ([]byte, error) {
	ret, err := done(b.GetContext(), func() (interface{}, error) {
		node, ok := b.data[key]
		if !ok {
			return nil, nilError
		}
		if b.checkExpire(node) {
			b.doDel(key)
			return nil, nilError
		}
		return node.data, nil
	})
	return toByte(ret, err)
}

func (b basis) Get(key string) []byte {
	data, _ := b.GetE(key)
	return data
}

func (b basis) GetWithCreateE(key string, h Handle) ([]byte, error) {
	ret, err := done(b.GetContext(), func() (interface{}, error) {
		data, err := b.GetE(key)
		if err == nil {
			return data, err
		}
		if !b.IsNilError(err) {
			return nil, err
		}
		doData, err, _ := b.loadGroup.Do(key, func() (interface{}, error) {
			data, err := h.Create(b.GetContext())
			if err == nil {
				b.doSetWithData(key, data, h.Expire())
			}
			return data, err
		})
		if err != nil {
			return nil, errors.Wrap(err, "x cache create data error")
		}
		return doData, err
	})
	return toByte(ret, err)
}

func (b basis) GetWithCreate(key string, h Handle) []byte {
	data, _ := b.GetWithCreateE(key, h)
	return data
}

func (b basis) doSetWithData(key string, data []byte, expire time.Duration) {
	b.data[key] = node{
		creatTime: time.Now(),
		data:      data,
		expire:    expire,
	}
}

func (b basis) Set(key string, h Handle) error {
	_, err := done(b.GetContext(), func() (interface{}, error) {
		_, err, _ := b.loadGroup.Do(key, func() (interface{}, error) {
			data, err := h.Create(b.GetContext())
			if err == nil {
				b.doSetWithData(key, data, h.Expire())
			}
			return data, err
		})
		if err != nil {
			return nil, errors.Wrap(err, "x cache create data error")
		}
		return nil, nil
	})
	return err
}

func (b basis) IsNilError(err error) bool {
	return err == nilError
}

func (b basis) IsExist(key string) bool {
	_, err := done(b.GetContext(), func() (interface{}, error) {
		return nil, b.doIsExist(key)
	})
	return err == nilError
}

func (b basis) checkExpire(n node) bool {
	if n.expire == 0 {
		return false
	}
	return time.Now().After(n.creatTime.Add(n.expire))
}

func (b basis) doIsExist(key string) error {
	if _, ok := b.data[key]; !ok {
		return nilError
	}
	return nil
}

func (b basis) doDel(key string) {
	delete(b.data, key)
}

func toByte(ret interface{}, err error) ([]byte, error) {
	if data, ok := ret.([]byte); ok {
		return data, err
	}
	return nil, err
}

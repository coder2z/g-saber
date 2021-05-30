package xcache

import (
	"github.com/pkg/errors"
	"golang.org/x/sync/singleflight"
	"time"
)

type (
	basis struct {
		data      map[string]node
		loadGroup *singleflight.Group
	}

	node struct {
		data      []byte
		expire    time.Duration
		creatTime time.Time
	}
)

func NewBasis() *basis {
	return &basis{
		data:      make(map[string]node),
		loadGroup: &singleflight.Group{},
	}
}

func (b basis) Del(keys ...string) error {
	for _, key := range keys {
		b.doDel(key)
	}
	return nil
}

func (b basis) GetE(key string) ([]byte, error) {
	node, ok := b.data[key]
	if !ok {
		return nil, nilError
	}
	if b.checkExpire(node) {
		b.doDel(key)
		return nil, nilError
	}
	return node.data, nil
}

func (b basis) Get(key string) []byte {
	data, _ := b.GetE(key)
	return data
}

func (b basis) GetWithCreateE(key string, h Handle) ([]byte, error) {
	data, err := b.GetE(key)
	if err == nil {
		return data, err
	}
	if !b.IsNilError(err) {
		return nil, err
	}
	do, err, _ := b.loadGroup.Do(key, func() (interface{}, error) {
		return h.Create()
	})
	if err != nil {
		return nil, errors.Wrap(err, "x cache create data error")
	}
	data, _ = do.([]byte)
	b.doSetWithData(key, data, h.Expire())
	return data, nil
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
	do, err, _ := b.loadGroup.Do(key, func() (interface{}, error) {
		return h.Create()
	})
	if err != nil {
		return errors.Wrap(err, "x cache create data error")
	}
	data, _ := do.([]byte)
	b.doSetWithData(key, data, h.Expire())
	return nil
}

func (b basis) IsNilError(err error) bool {
	return err == nilError
}

func (b basis) IsExist(key string) bool {
	return b.doIsExist(key) == nilError
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
package xcache

import (
	"testing"
	"time"
)

type testdata struct {
}

func (t testdata) Create() ([]byte, error) {
	return []byte("静安里看风景阿卡丽大家啊就完蛋了卡"), nil
}

func (t testdata) Expire() time.Duration {
	return 0
}

func TestName(t *testing.T) {
	b := NewBasis()

	err := b.Set("key", new(testdata))
	if err != nil {
		return
	}
	data, err := b.Get("key")
	if err != nil {
		return
	}
	t.Log(string(data))
}

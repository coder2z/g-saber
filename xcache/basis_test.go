package xcache

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type testdata struct{}

func (t testdata) Create(ctx context.Context) ([]byte, error) {
	time.Sleep(1900 * time.Millisecond)
	return []byte("test Create"), nil
}

func (t testdata) Expire() time.Duration {
	return time.Second * 3
}

func BenchmarkBasis(b *testing.B) {
	b.ReportAllocs()
	b.StopTimer()
	c := NewBasis()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		_ = c.Set(fmt.Sprintf("key_%d", i), new(testdata))
		if string(c.Get(fmt.Sprintf("key_%d", i))) != "test" {
			b.Error("data error ")
		}
	}
}

func TestBasis(t *testing.T) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelFunc()

	c := NewBasis().WithContext(ctx)
	if err := c.Set("abc", new(testdata)); err != nil {
		t.Error(err)
		return
	}
	data, err := c.GetE("abc")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(string(data))
}

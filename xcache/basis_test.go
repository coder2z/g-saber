package xcache

import (
	"fmt"
	"testing"
	"time"
)

type testdata struct{}

func (t testdata) Create() ([]byte, error) {
	return []byte("test"), nil
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

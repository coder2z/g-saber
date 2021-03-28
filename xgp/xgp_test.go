package xgp

import (
	"github.com/coder2z/g-saber/xconsole"
	"testing"
	"time"
)

func TestPG(t *testing.T) {
	pool := NewWaitPool(2)
	pool.Submit(func() {
		xconsole.Red("test1")
		time.Sleep(2 * time.Second)
		xconsole.Red("test1 Ok")
	},nil)

	pool.Submit(func() {
		xconsole.Red("test2")
		time.Sleep(3 * time.Second)
		xconsole.Red("test2 Ok")
	},nil)

	pool.Submit(func() {
		xconsole.Red("test3")
		time.Sleep(4 * time.Second)
		xconsole.Red("test3 Ok")
	},nil)

	pool.Submit(func() {
		time.Sleep(8 * time.Second)
		panic("test4")
	}, func(err error) {
		xconsole.Red(err.Error())
	})

	pool.Wait()
}

package file

import (
	"github.com/coder2z/g-saber/xcfg"
	"testing"
)

func TestNewDataSource(t *testing.T) {
	err := xcfg.LoadFromConfigAddr("config.toml")
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(xcfg.Get("app"))
}


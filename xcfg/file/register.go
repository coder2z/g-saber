package file

import (
	"github.com/coder2z/g-saber/xcfg"
)

func init() {
	xcfg.Register(new(file))
}

// DataSourceFile defines file scheme
const DataSourceFile = "file"

type file struct{}

func (f file) Register() (string, xcfg.DataSourceCreatorFunc) {
	return DataSourceFile, func(configAddr string) xcfg.DataSource {
		if configAddr == "" {
			return nil
		}
		return NewDataSource(configAddr, false)
	}
}

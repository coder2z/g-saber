package xcfg

import (
	"errors"
	"fmt"
	"github.com/coder2z/g-saber/xconsole"
	"net/url"
	"sync"
)

var (
	//ErrConfigAddr not xcfg
	ErrConfigAddr = errors.New("no xcfg... ")
	// ErrInvalidDataSource defines an error that the scheme has been registered
	ErrInvalidDataSource = errors.New("invalid data source, please make sure the scheme has been registered")
	registry             = make(map[string]DataSourceCreatorFunc)
	syncM                sync.RWMutex
	//DefaultScheme ..
	DefaultScheme = `file`
)

// DataSourceCreatorFunc represents a dataSource creator function
type DataSourceCreatorFunc func(addr string) DataSource

type DataSourceImp interface {
	Register() (string, DataSourceCreatorFunc)
}

// Register registers a dataSource creator function to the registry
func Register(data ...DataSourceImp) {
	syncM.Lock()
	defer syncM.Unlock()
	for _, datum := range data {
		name, f := datum.Register()
		registry[name] = f
	}
}

//NewDataSource ..
func NewDataSource(configAddr string) (DataSource, error) {
	var (
		scheme = DefaultScheme
	)
	if configAddr == "" {
		return nil, ErrConfigAddr
	}
	urlObj, err := url.Parse(configAddr)
	if err == nil && len(urlObj.Scheme) > 1 {
		scheme = urlObj.Scheme
	}
	syncM.RLock()
	creatorFunc, exist := registry[scheme]
	syncM.RUnlock()
	if !exist {
		return nil, ErrInvalidDataSource
	}
	xconsole.Green(fmt.Sprintf("Get xcfg from:%s", configAddr))
	source := creatorFunc(configAddr)
	if source == nil {
		return nil, ErrInvalidDataSource
	}
	return source, nil
}

package etcdv3

import (
	"github.com/coder2z/g-saber/xcfg"
	"github.com/spf13/cast"
	"net/url"
	"time"

	"go.etcd.io/etcd/clientv3"
)

func init() {
	xcfg.Register(new(etcd))
}

// DataSourceEtcd defines etcd scheme
const DataSourceEtcd = "etcd"

type etcd struct{}

func (e etcd) Register() (string, xcfg.DataSourceCreatorFunc) {
	return DataSourceEtcd, func(configAddr string) xcfg.DataSource {
		if configAddr == "" {
			return nil
		}
		var (
			watch bool
		)
		// configAddr is a string in this format:
		// etcd://ip:port?username=XXX&password=XXX&key=key&watch=false
		urlObj, err := url.Parse(configAddr)
		if err != nil {
			return nil
		}
		watch = cast.ToBool(urlObj.Query().Get("watch"))
		etcdConf := clientv3.Config{
			DialKeepAliveTime:    10 * time.Second,
			DialKeepAliveTimeout: 3 * time.Second,
			Endpoints:[]string{urlObj.Host},
			Username: urlObj.Query().Get("username"),
			Password:urlObj.Query().Get("password"),
		}
		client, err := clientv3.New(etcdConf)
		if err != nil {
			return nil
		}
		return NewDataSource(client, urlObj.Query().Get("key"), watch)
	}
}

package apollo

import (
	"github.com/coder2z/g-saber/xcfg"
	"github.com/philchia/agollo/v4"
	"github.com/spf13/cast"
	"net/url"
)

func init() {
	xcfg.Register(new(apollo))
}

// DataSourceApollo defines apollo scheme
const DataSourceApollo = "apollo"

type apollo struct{}

func (e apollo) Register() (string, xcfg.DataSourceCreatorFunc) {
	return DataSourceApollo, func(addr string) xcfg.DataSource {
		var watch bool
		if addr == "" {
			return nil
		}
		// configAddr is a string in this format:
		// apollo://ip:port?appId=XXX&cluster=XXX&namespaceName=XXX&key=XXX&accesskeySecret=XXX&insecureSkipVerify=XXX&cacheDir=XXX&watch=true
		urlObj, err := url.Parse(addr)
		if err != nil {
			return nil
		}
		watch = cast.ToBool(urlObj.Query().Get("watch"))
		apolloConf := agollo.Conf{
			AppID:              urlObj.Query().Get("appId"),
			Cluster:            urlObj.Query().Get("cluster"),
			NameSpaceNames:     []string{urlObj.Query().Get("namespaceName")},
			MetaAddr:           urlObj.Host,
			InsecureSkipVerify: cast.ToBool(urlObj.Query().Get("insecureSkipVerify")),
			AccesskeySecret:    urlObj.Query().Get("accesskeySecret"),
			CacheDir:           ".",
		}
		if urlObj.Query().Get("cacheDir") != "" {
			apolloConf.CacheDir = urlObj.Query().Get("cacheDir")
		}
		return NewDataSource(&apolloConf, urlObj.Query().Get("namespaceName"), urlObj.Query().Get("key"), watch)
	}
}

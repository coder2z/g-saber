module github.com/coder2m/g-saber

go 1.16

require (
	github.com/coreos/etcd v3.3.13+incompatible
	github.com/davecgh/go-spew v1.1.1
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.4.1
	github.com/golang/protobuf v1.4.3
	github.com/json-iterator/go v1.1.10
	github.com/mitchellh/mapstructure v1.1.2
	github.com/modern-go/reflect2 v1.0.1
	github.com/philchia/agollo/v4 v4.1.3
	github.com/pkg/errors v0.8.1
	github.com/robfig/cron/v3 v3.0.1
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.3
	go.uber.org/zap v1.10.0
	google.golang.org/genproto v0.0.0-20210303154014-9728d6b83eeb
	google.golang.org/grpc v1.36.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

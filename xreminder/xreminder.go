/**
* @Author: myxy99 <myxy99@foxmail.com>
* @Date: 2020/11/4 15:27
 */
package xreminder

import (
	"github.com/robfig/cron/v3"
)

type config struct {
	Spec []string
}

type OptionFunc func(o *config)

func SetSpec(s []string) OptionFunc {
	return func(o *config) {
		o.Spec = s
	}
}

func NewReminderCfg(of ...OptionFunc) *config {
	o := new(config)
	for _, optionFunc := range of {
		optionFunc(o)
	}
	return o
}

type cronServer struct {
	o *config
	c *cron.Cron
}

func (r *cronServer) Run(stopCh <-chan struct{}, f cron.Job) {
	for _, v := range r.o.Spec {
		_, _ = r.c.AddJob(v, f)
	}
	r.c.Start()
	<-stopCh
	r.c.Stop()
}

func NewReminderClient(o *config) *cronServer {
	return &cronServer{c: cron.New(), o: o}
}

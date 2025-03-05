package kqx

import (
	"github.com/zeromicro/go-queue/kq"
	"sync"
)

var once sync.Once
var gc *GroupConsumer

type GroupConsumer struct {
	group *ServiceGroup
}

func NewGroupConsumerInstance() *GroupConsumer {
	once.Do(func() {
		gc = &GroupConsumer{
			group: NewServiceGroup(),
		}
	})
	return gc
}

func (gc *GroupConsumer) Register(conf KqConf, queue kq.ConsumeHandler) *GroupConsumer {
	gc.group.Add(MustNewQueue(conf, queue))
	return gc
}

func (gc *GroupConsumer) Start() {
	gc.group.Start()
}

func (gc *GroupConsumer) StopSignal() {
	gc.group.Stop()
}

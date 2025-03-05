package svc

import "github.com/simance-ai/smdx/rpcx/ws/internal/config"

type ServiceContext struct {
	Config           config.Config
	WebsocketContext *WebsocketContext
}

func NewServiceContext(c config.Config, wsc *WebsocketContext) *ServiceContext {
	return &ServiceContext{
		Config:           c,
		WebsocketContext: wsc,
	}
}

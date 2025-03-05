package svc

import (
	"github.com/simance-ai/smdx/rpcx/chains/eth/internal/config"
)

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}

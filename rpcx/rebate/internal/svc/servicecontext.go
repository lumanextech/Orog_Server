package svc

import (
	"github.com/simance-ai/smdx/rpcx/account/accountclient"
	"github.com/simance-ai/smdx/rpcx/order/orderclient"
	"github.com/simance-ai/smdx/rpcx/rebate/internal/config"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	AccountClient accountclient.Account
	OrderClient   orderclient.Order
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:        c,
		AccountClient: accountclient.NewAccount(zrpc.MustNewClient(c.AccountClientConf)),
		OrderClient:   orderclient.NewOrder(zrpc.MustNewClient(c.OrderClientConf)),
	}
}

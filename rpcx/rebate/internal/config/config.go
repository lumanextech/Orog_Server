package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf

	AccountClientConf zrpc.RpcClientConf
	OrderClientConf   zrpc.RpcClientConf
	Corn              Cron
}

type Cron struct {
	Rebate string
}

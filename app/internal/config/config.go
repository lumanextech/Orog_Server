package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf

	SolClientConf     zrpc.RpcClientConf
	AccountClientConf zrpc.RpcClientConf
	OrderClientConf   zrpc.RpcClientConf
	Auth              struct {
		AccessSecret string
		AccessExpire int64
	}
}

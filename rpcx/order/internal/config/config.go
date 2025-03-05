package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Cache    cache.CacheConf
	PkEncode struct {
		Key string
		Iv  string
	}

	DBConf struct {
		DSN string
	}
}

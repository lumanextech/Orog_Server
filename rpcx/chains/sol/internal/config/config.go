package config

import (
	"github.com/SpectatorNan/gorm-zero/gormc/config/pg"
	okxos "github.com/simance-ai/smdx/pkg/dexapi/okx_os"
	solscan "github.com/simance-ai/smdx/pkg/dexapi/sol_scan"
	"github.com/simance-ai/smdx/pkg/kqx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	PgSqlConf     pg.PgSql
	PgSqlReadConf pg.PgSql

	Cache struct {
		Redis cache.CacheConf
	}

	Chains []struct {
		Name string
		Rpc  string
		Ws   string
	}

	PkEncode struct {
		Key string
		Iv  string
	}

	KafkaConf kqx.KqConf

	AccountClientConf zrpc.RpcClientConf
	WsClientConf      zrpc.RpcClientConf

	OkxOsConf   okxos.Config
	SolScanConf solscan.Config
}

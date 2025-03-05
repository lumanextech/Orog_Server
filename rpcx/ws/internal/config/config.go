package config

import (
	"github.com/simance-ai/smdx/pkg/kqx"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type WebSocketConf struct {
	ReadBufferSize  int
	WriteBufferSize int

	MessageFormat             int
	WebsocketHandshakeTimeout int
	EnableCompression         int
}

type AppConfig struct {
	Secret string
}

type Config struct {
	zrpc.RpcServerConf

	WebSocket WebSocketConf
	Rest      rest.RestConf
	App       AppConfig

	Cache struct {
		Redis cache.CacheConf
	}

	MarketKlineTopicKafkaConf kqx.KqConf

	MarketSwapTopicKafkaConf kqx.KqConf
}

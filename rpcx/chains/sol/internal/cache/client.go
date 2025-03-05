package cache

import (
	"context"

	"github.com/alibaba/tair-go/tair"
	"github.com/bsm/redislock"
	solrpc "github.com/gagliardetto/solana-go/rpc"
	"github.com/simance-ai/smdx/pkg/dexapi/sol_scan"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type DB struct {
	redisClient   *tair.TairClusterClient
	pgDB          *gorm.DB
	solscanClient *sol_scan.Client
	solClient     *solrpc.Client
	logx.Logger
}

func NewDB(
	redisClient *tair.TairClusterClient,
	pgDB *gorm.DB,
	solscanClient *sol_scan.Client,
	solClient *solrpc.Client,
) *DB {
	return &DB{
		redisClient:   redisClient,
		pgDB:          pgDB,
		solscanClient: solscanClient,
		solClient:     solClient,
		Logger:        logx.WithContext(context.Background()),
	}
}

func (c *DB) CreateLock() *redislock.Client {
	return redislock.New(c.redisClient)
}

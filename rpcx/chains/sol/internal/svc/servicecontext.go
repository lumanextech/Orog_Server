package svc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"time"

	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/cache"
	tsmodule "github.com/simance-ai/smdx/rpcx/chains/sol/internal/cache/ts_module"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/dao"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"

	"github.com/alibaba/tair-go/tair"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/simance-ai/smdx/rpcx/account/accountclient"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/config"
	"github.com/simance-ai/smdx/rpcx/ws/wsclient"
	"github.com/zeromicro/go-zero/zrpc"

	_ "github.com/lib/pq"
	solscan "github.com/simance-ai/smdx/pkg/dexapi/sol_scan"
)

type ServiceContext struct {
	Config    config.Config
	CacheDB   *cache.DB
	SolClient *rpc.Client
	PgDB      *gorm.DB
	ReadPgDB  *gorm.DB

	AccountClient accountclient.Account
	WsClient      wsclient.Ws
	SolScanClient *solscan.Client

	TsModule tsmodule.Ts
}

func NewServiceContext(c config.Config) *ServiceContext {

	if len(c.Cache.Redis) <= 0 {
		log.Fatalf("redis config is empty")
	}

	if len(c.Chains) <= 0 {
		log.Fatalf("chains config is empty")
	}

	redisClusterConf := c.Cache.Redis[0]
	var tlsConfig *tls.Config // Explicitly define the type
	if redisClusterConf.Tls {
		tlsConfig = &tls.Config{
			InsecureSkipVerify: redisClusterConf.Tls,
		}
	}
	redisClient := tair.NewTairClusterClient(&tair.TairClusterOptions{
		ClusterOptions: &redis.ClusterOptions{
			Addrs:     []string{redisClusterConf.Host},
			Password:  redisClusterConf.Pass,
			TLSConfig: tlsConfig,
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	successStatusErr := redisClient.Ping(ctx).Err()
	if successStatusErr != nil {
		log.Fatalf("redis connect error: %v", successStatusErr)
	}

	defaultNode := c.Chains[0]
	cluster0 := rpc.Cluster{
		Name: defaultNode.Name,
		WS:   defaultNode.Ws,
		RPC:  defaultNode.Rpc,
	}
	rpcClient := rpc.New(cluster0.RPC)

	defaultConfig := &gorm.Config{
		Logger: dao.NewZeroLog(gormLogger.Config{
			SlowThreshold:             c.PgSqlConf.GetSlowThreshold(),
			Colorful:                  c.PgSqlConf.GetColorful(),
			LogLevel:                  c.PgSqlConf.GetGormLogMode(),
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      true,
		}),
	}

	//ssl mode
	c = rewriteSslMode(c)

	if c.PgSqlReadConf.SslMode == "enable" {
		log.Fatalf("sslmode is require, please check your config")
	}

	dbClient, err := gorm.Open(postgres.Open(c.PgSqlConf.Dsn()), defaultConfig)
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}

	//use dbresolver plugin
	err = dbClient.Use(dbresolver.Register(dbresolver.Config{
		Replicas:          []gorm.Dialector{postgres.Open(c.PgSqlReadConf.Dsn())}, //only one
		Policy:            dbresolver.RandomPolicy{},
		TraceResolverMode: true,
	}))
	if err != nil {
		log.Fatalf("failed opening dbresolver connection to postgres: %v", err)
	}

	//init solscan client
	solscanClient := solscan.NewClient(c.SolScanConf)

	return &ServiceContext{
		Config:        c,
		SolClient:     rpcClient,
		PgDB:          dbClient,
		ReadPgDB:      dbClient,
		AccountClient: accountclient.NewAccount(zrpc.MustNewClient(c.AccountClientConf)),
		WsClient:      wsclient.NewWs(zrpc.MustNewClient(c.WsClientConf)),
		TsModule:      tsmodule.NewRedisTs(redisClient),
		CacheDB:       cache.NewDB(redisClient, dbClient, solscanClient, rpcClient),
		SolScanClient: solscanClient,
	}
}

func rewriteSslMode(c config.Config) config.Config {
	if c.PgSqlConf.SslMode == "enable" {
		rootCertPool := x509.NewCertPool()
		pem, err := os.ReadFile(c.KafkaConf.CaFile)
		if err != nil {
			log.Fatalf("failed to read pem file: %v", err)
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			log.Fatalf("failed to append pem")
		}
		stdlib.RegisterConnConfig(&pgx.ConnConfig{
			Config: pgconn.Config{
				TLSConfig: &tls.Config{RootCAs: rootCertPool},
			},
		})

		c.PgSqlConf.SslMode = "require"
	}

	if c.PgSqlReadConf.SslMode == "enable" {
		rootCertPool := x509.NewCertPool()
		pem, err := os.ReadFile(c.KafkaConf.CaFile)
		if err != nil {
			log.Fatalf("failed to read pem file: %v", err)
		}
		if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
			log.Fatalf("failed to append pem")
		}
		stdlib.RegisterConnConfig(&pgx.ConnConfig{
			Config: pgconn.Config{
				TLSConfig: &tls.Config{RootCAs: rootCertPool},
			},
		})

		c.PgSqlReadConf.SslMode = "require"
	}

	return c
}

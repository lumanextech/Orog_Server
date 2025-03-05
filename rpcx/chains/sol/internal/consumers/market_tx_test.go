package consumers_test

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"

	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/config"
	"github.com/simance-ai/smdx/rpcx/chains/sol/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
)

var ErrKeyExists = errors.New("ERR TSDB: key already exists")

var TestSolContext *SolContext

type SolContext struct {
	SrvCtx *svc.ServiceContext
	Ctx    context.Context
}

func Init() {
	dir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("error: %v \n", err)
	}

	configPath := dir + "/Documents/GitHub/smdx/rpcx/chains/sol/etc/sol.yaml"

	var configFile = flag.String("f", configPath, "the config file")

	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)
	ctxB := context.Background()

	TestSolContext = &SolContext{
		SrvCtx: ctx,
		Ctx:    ctxB,
	}
}

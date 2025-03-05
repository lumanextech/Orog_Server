package main

import (
	"flag"
	"fmt"
	timejob "github.com/simance-ai/smdx/rpcx/rebate/internal/logic/timejobs"
	"log"

	"github.com/simance-ai/smdx/rpcx/rebate/internal/config"
	"github.com/simance-ai/smdx/rpcx/rebate/internal/server"
	"github.com/simance-ai/smdx/rpcx/rebate/internal/svc"
	"github.com/simance-ai/smdx/rpcx/rebate/rebate"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/rebate.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		rebate.RegisterRebateServer(grpcServer, server.NewRebateServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	cron, err := timejob.Run(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer cron.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

package main

import (
	"flag"
	"fmt"

	"github.com/simance-ai/smdx/rpcx/chains/bsc/bsc"
	"github.com/simance-ai/smdx/rpcx/chains/bsc/internal/config"
	"github.com/simance-ai/smdx/rpcx/chains/bsc/internal/server"
	"github.com/simance-ai/smdx/rpcx/chains/bsc/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/bsc.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		bsc.RegisterBscServer(grpcServer, server.NewBscServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

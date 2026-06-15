package main

import (
	"flag"
	"fmt"

	"smartcommunity-microservices/app/workorder/rpc/internal/config"
	"smartcommunity-microservices/app/workorder/rpc/internal/server"
	"smartcommunity-microservices/app/workorder/rpc/internal/svc"
	"smartcommunity-microservices/app/workorder/rpc/types/workorder"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/workorder.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		workorder.RegisterWorkorderRpcServer(grpcServer, server.NewWorkorderRpcServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

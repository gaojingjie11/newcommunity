package main

import (
	"flag"
	"fmt"

	"smartcommunity-microservices/app/mall/rpc/internal/config"
	"smartcommunity-microservices/app/mall/rpc/internal/server"
	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/mall-rpc.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())
	ctx := svc.NewServiceContext(c)
	defer ctx.TimeoutSvc.Stop()
	defer ctx.PaymentReconcileSvc.Stop()

	// Start RabbitMQ Notification Consumer in background
	go svc.StartNotificationConsumer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		mall.RegisterMallRpcServer(grpcServer, server.NewMallRpcServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}

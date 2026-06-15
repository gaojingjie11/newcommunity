// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"context"
	"log"

	"smartcommunity-microservices/app/agent/rpc/agentrpc"
	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/config"
	"smartcommunity-microservices/app/mall/rpc/mallrpc"
	"smartcommunity-microservices/app/stats/rpc/statsrpc"
	"smartcommunity-microservices/app/user/rpc/userrpc"
	"smartcommunity-microservices/app/workorder/rpc/workorderrpc"
	"smartcommunity-microservices/common/redis"
	"smartcommunity-microservices/common/storage"

	"github.com/minio/minio-go/v7"
	goredis "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ServiceContext struct {
	Config       config.Config
	UserRpc      userrpc.UserRpc
	MallRpc      mallrpc.MallRpc
	CommunityRpc communityrpc.CommunityRpc
	WorkorderRpc workorderrpc.WorkorderRpc
	StatsRpc     statsrpc.StatsRpc
	AgentRpc     agentrpc.AgentRpc
	RedisClient  *goredis.Client
	MinioClient  *minio.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	rdb, err := redis.Init(c.BizRedis)
	if err != nil {
		log.Fatalf("failed to init redis in gateway-api: %v", err)
	}

	minioClient, err := storage.Init(c.MinIO)
	if err != nil {
		log.Printf("warning: minio init failed in gateway-api, upload will be unavailable: %v", err)
	}

	return &ServiceContext{
		Config:       c,
		UserRpc:      userrpc.NewUserRpc(zrpc.MustNewClient(c.UserRpc)),
		MallRpc:      mallrpc.NewMallRpc(zrpc.MustNewClient(c.MallRpc, zrpc.WithUnaryClientInterceptor(gatewayGrpcInterceptor))),
		CommunityRpc: communityrpc.NewCommunityRpc(zrpc.MustNewClient(c.CommunityRpc)),
		WorkorderRpc: workorderrpc.NewWorkorderRpc(zrpc.MustNewClient(c.WorkorderRpc)),
		StatsRpc:     statsrpc.NewStatsRpc(zrpc.MustNewClient(c.StatsRpc)),
		AgentRpc:     agentrpc.NewAgentRpc(zrpc.MustNewClient(c.AgentRpc)),
		RedisClient:  rdb,
		MinioClient:  minioClient,
	}
}

func gatewayGrpcInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	ids, _ := ctx.Value("x-store-ids").(string)
	isAdmin, _ := ctx.Value("x-is-admin").(bool)
	log.Printf("[DEBUG] gatewayGrpcInterceptor: method=%s, x-store-ids=%s, x-is-admin=%t\n", method, ids, isAdmin)

	if ids != "" {
		ctx = metadata.AppendToOutgoingContext(ctx, "x-store-ids", ids)
	}
	if isAdmin {
		ctx = metadata.AppendToOutgoingContext(ctx, "x-is-admin", "true")
	}
	return invoker(ctx, method, req, reply, cc, opts...)
}

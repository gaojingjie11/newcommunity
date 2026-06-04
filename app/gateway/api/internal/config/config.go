// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package config

import (
	"smartcommunity-microservices/common/redis"
	"smartcommunity-microservices/common/storage"

	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	UserRpc      zrpc.RpcClientConf
	MallRpc      zrpc.RpcClientConf
	CommunityRpc zrpc.RpcClientConf
	WorkorderRpc zrpc.RpcClientConf
	StatsRpc     zrpc.RpcClientConf
	BizRedis     redis.RedisConfig
	MinIO        storage.MinIOConfig
	Gateway struct {
		InternalToken string
		Services      map[string]string
	}
}

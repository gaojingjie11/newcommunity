package config

import (
	"smartcommunity-microservices/common/db"
	"smartcommunity-microservices/common/redis"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	MySQL    db.MySQLConfig
	BizRedis redis.RedisConfig
	MallRpc  zrpc.RpcClientConf
}

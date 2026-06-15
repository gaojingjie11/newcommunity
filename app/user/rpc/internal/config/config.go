package config

import (
	"smartcommunity-microservices/common/db"
	"smartcommunity-microservices/common/mq"
	"smartcommunity-microservices/common/redis"
	"smartcommunity-microservices/common/storage"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Postgres db.PostgresConfig
	BizRedis redis.RedisConfig
	RabbitMQ mq.RabbitMQConfig
	MinIO    storage.MinIOConfig
	Jwt      struct {
		Secret string
		TTL    string
	}
}


package config

import (
	"smartcommunity-microservices/common/db"
	"smartcommunity-microservices/common/mq"
	"smartcommunity-microservices/common/redis"
	"smartcommunity-microservices/common/storage"

	"github.com/zeromicro/go-zero/zrpc"
)

type AgentConfig struct {
	BaseUrl    string `json:",optional"`
	LlmApiKey  string `json:",optional"`
	LlmBaseUrl string `json:",optional"`
	LlmModel   string `json:",optional"`
}

type Config struct {
	zrpc.RpcServerConf
	Postgres db.PostgresConfig
	BizRedis redis.RedisConfig
	Agent    AgentConfig       `json:",optional"`
	RabbitMQ mq.RabbitMQConfig `json:",optional"`
	MinIO    storage.MinIOConfig `json:",optional"`
}


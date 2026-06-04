package config

import (
	"smartcommunity-microservices/common/db"
	"smartcommunity-microservices/common/redis"

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
	MySQL    db.MySQLConfig
	BizRedis redis.RedisConfig
	Agent    AgentConfig `json:",optional"`
}

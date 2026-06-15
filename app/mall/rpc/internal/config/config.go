package config

import (
	"smartcommunity-microservices/common/db"
	"smartcommunity-microservices/common/mail"
	"smartcommunity-microservices/common/mq"
	"smartcommunity-microservices/common/redis"

	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Postgres db.PostgresConfig
	BizRedis redis.RedisConfig
	RabbitMQ mq.RabbitMQConfig
	Mail     mail.MailConfig
	Alipay   struct {
		AppID           string
		PrivateKey      string
		AlipayPublicKey string
		NotifyURL       string
		ReturnURL       string
		Sandbox         bool
	}
}

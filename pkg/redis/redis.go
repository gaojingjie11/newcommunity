package redis

import (
	"context"
	"time"

	"smartcommunity-microservices/pkg/config"

	goredis "github.com/redis/go-redis/v9"
)

func Init(cfg config.RedisConfig) (*goredis.Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     cfg.Addr(),
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return client, client.Ping(ctx).Err()
}

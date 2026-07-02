package main

import (
	"context"
	"log"
	"time"

	commonredis "smartcommunity-microservices/common/redis"
)

func main() {
	client, err := commonredis.Init(commonredis.RedisConfig{
		Host:     "101.42.34.232",
		Port:     6379,
		Password: "dsw123456",
		DB:       15,
	})
	if err != nil {
		log.Fatalf("connect redis failed: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.FlushDB(ctx).Err(); err != nil {
		log.Fatalf("flush redis db 15 failed: %v", err)
	}

	log.Println("redis db 15 flushed")
}

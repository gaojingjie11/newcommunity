package svc

import (
	"log"

	"smartcommunity-microservices/app/workorder/rpc/internal/config"
	"smartcommunity-microservices/app/workorder/rpc/internal/repository"
	"smartcommunity-microservices/app/workorder/rpc/internal/service"
	"smartcommunity-microservices/common/db"
	"smartcommunity-microservices/common/logger"
	"smartcommunity-microservices/common/mq"
	"smartcommunity-microservices/common/redis"

	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config        config.Config
	DB            *gorm.DB
	Redis         *goredis.Client
	WorkorderRepo *repository.WorkorderRepo
	EventBus      *service.EventBus
}

func NewServiceContext(c config.Config) *ServiceContext {
	logr := logger.New(c.Name)

	// Initialize MySQL
	database, err := db.InitMySQL(c.MySQL)
	if err != nil {
		log.Fatalf("failed to init mysql in workorder-rpc: %v", err)
	}

	// Initialize Redis
	rdb, err := redis.Init(c.BizRedis)
	if err != nil {
		log.Fatalf("failed to init redis in workorder-rpc: %v", err)
	}

	// Initialize MQ
	var mqClient *mq.Client
	if c.RabbitMQ.URL() != "" {
		mqClient, err = mq.Connect(c.RabbitMQ)
		if err != nil {
			logr.Warn("rabbitmq connect failed, events disabled", "error", err)
		} else {
			logr.Info("rabbitmq connected")
		}
	}

	eventBus := service.NewEventBus(mqClient, logr)
	workorderRepo := repository.NewWorkorderRepo(database)

	return &ServiceContext{
		Config:        c,
		DB:            database,
		Redis:         rdb,
		WorkorderRepo: workorderRepo,
		EventBus:      eventBus,
	}
}

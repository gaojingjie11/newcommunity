package svc

import (
	"log"

	"smartcommunity-microservices/app/community/rpc/internal/config"
	"smartcommunity-microservices/app/community/rpc/internal/model"
	"smartcommunity-microservices/app/community/rpc/internal/repository"
	"smartcommunity-microservices/app/mall/rpc/mallrpc"
	"smartcommunity-microservices/common/db"
	"smartcommunity-microservices/common/mq"
	"smartcommunity-microservices/common/redis"

	goredis "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config          config.Config
	DB              *gorm.DB
	Redis           *goredis.Client
	MQ              *mq.Client
	NoticeRepo      *repository.NoticeRepo
	VisitorRepo     *repository.VisitorRepo
	ParkingRepo     *repository.ParkingRepo
	PropertyFeeRepo *repository.PropertyFeeRepo
	MessageRepo     *repository.MessageRepo
	MallRpc         mallrpc.MallRpc
}

func NewServiceContext(c config.Config) *ServiceContext {
	// Initialize Postgres
	database, err := db.InitPostgres(c.Postgres)
	if err != nil {
		log.Fatalf("failed to init postgres in community-rpc: %v", err)
	}

	// Run AutoMigrate for PostgreSQL tables
	_ = database.AutoMigrate(
		&model.Notice{},
		&model.Visitor{},
		&model.ParkingSpace{},
		&model.UserParkingBinding{},
		&model.PropertyFee{},
		&model.PropertyFeePayment{},
		&model.CommunityMessage{},
	)

	if err := database.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_user_parking_bindings_active_plate_unique
		ON user_parking_bindings (car_plate)
		WHERE status = 1 AND car_plate <> ''
	`).Error; err != nil {
		log.Printf("failed to ensure active car plate unique index: %v", err)
	}

	// Initialize Redis
	rdb, err := redis.Init(c.BizRedis)
	if err != nil {
		log.Fatalf("failed to init redis in community-rpc: %v", err)
	}

	// Initialize MQ
	var mqClient *mq.Client
	if c.RabbitMQ.URL() != "" {
		var mqErr error
		mqClient, mqErr = mq.Connect(c.RabbitMQ)
		if mqErr != nil {
			log.Printf("rabbitmq connect failed in community-rpc: %v", mqErr)
		} else {
			log.Printf("rabbitmq connected in community-rpc")
		}
	}

	// Initialize Repositories
	noticeRepo := repository.NewNoticeRepo(database)
	visitorRepo := repository.NewVisitorRepo(database)
	parkingRepo := repository.NewParkingRepo(database)
	feeRepo := repository.NewPropertyFeeRepo(database)
	messageRepo := repository.NewMessageRepo(database)

	// Initialize MallRpc Client
	mallClient := mallrpc.NewMallRpc(zrpc.MustNewClient(c.MallRpc))

	return &ServiceContext{
		Config:          c,
		DB:              database,
		Redis:           rdb,
		MQ:              mqClient,
		NoticeRepo:      noticeRepo,
		VisitorRepo:     visitorRepo,
		ParkingRepo:     parkingRepo,
		PropertyFeeRepo: feeRepo,
		MessageRepo:     messageRepo,
		MallRpc:         mallClient,
	}
}

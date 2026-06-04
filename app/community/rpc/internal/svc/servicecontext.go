package svc

import (
	"log"

	"smartcommunity-microservices/app/community/rpc/internal/config"
	"smartcommunity-microservices/app/community/rpc/internal/repository"
	"smartcommunity-microservices/app/mall/rpc/mallrpc"
	"smartcommunity-microservices/common/db"
	"smartcommunity-microservices/common/redis"

	goredis "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config          config.Config
	DB              *gorm.DB
	Redis           *goredis.Client
	NoticeRepo      *repository.NoticeRepo
	VisitorRepo     *repository.VisitorRepo
	ParkingRepo     *repository.ParkingRepo
	PropertyFeeRepo *repository.PropertyFeeRepo
	MessageRepo     *repository.MessageRepo
	MallRpc         mallrpc.MallRpc
}

func NewServiceContext(c config.Config) *ServiceContext {
	// Initialize MySQL
	database, err := db.InitMySQL(c.MySQL)
	if err != nil {
		log.Fatalf("failed to init mysql in community-rpc: %v", err)
	}

	// Initialize Redis
	rdb, err := redis.Init(c.BizRedis)
	if err != nil {
		log.Fatalf("failed to init redis in community-rpc: %v", err)
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
		NoticeRepo:      noticeRepo,
		VisitorRepo:     visitorRepo,
		ParkingRepo:     parkingRepo,
		PropertyFeeRepo: feeRepo,
		MessageRepo:     messageRepo,
		MallRpc:         mallClient,
	}
}

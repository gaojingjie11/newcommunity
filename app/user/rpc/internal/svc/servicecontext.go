package svc

import (
	"log"
	"time"

	"smartcommunity-microservices/common/db"
	"smartcommunity-microservices/common/mq"
	"smartcommunity-microservices/common/redis"
	"smartcommunity-microservices/app/user/rpc/internal/config"
	"smartcommunity-microservices/app/user/rpc/internal/repository"
	"smartcommunity-microservices/app/user/rpc/internal/service"

	goredis "github.com/redis/go-redis/v9"
)

type ServiceContext struct {
	Config          config.Config
	AuthService     *service.AuthService
	UserService     *service.UserService
	AdminService    *service.AdminService
	LoginLogService *service.LoginLogService
	RedisClient     *goredis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	database, err := db.InitMySQL(c.MySQL)
	if err != nil {
		log.Fatalf("failed to init mysql in user-rpc: %v", err)
	}

	rdb, err := redis.Init(c.BizRedis)
	if err != nil {
		log.Fatalf("failed to init redis in user-rpc: %v", err)
	}

	var mqClient *mq.Client
	if c.RabbitMQ.Host != "" {
		mqClient, err = mq.Connect(c.RabbitMQ)
		if err != nil {
			log.Printf("warning: rabbitmq connection failed in user-rpc: %v", err)
		}
	}

	jwtTTL, err := time.ParseDuration(c.Jwt.TTL)
	if err != nil {
		jwtTTL = 24 * time.Hour
	}

	userRepo := repository.NewUserRepo(database)
	roleRepo := repository.NewRoleRepo(database)
	logRepo := repository.NewLoginLogRepo(database)
	resetRepo := repository.NewPasswordResetRepo(database)

	return &ServiceContext{
		Config:          c,
		AuthService:     service.NewAuthService(userRepo, roleRepo, logRepo, resetRepo, rdb, c.Jwt.Secret, jwtTTL),
		UserService:     service.NewUserService(userRepo, rdb, mqClient),
		AdminService:    service.NewAdminService(userRepo, roleRepo, rdb),
		LoginLogService: service.NewLoginLogService(logRepo),
		RedisClient:     rdb,
	}
}

package svc

import (
	"log"
	"time"

	"smartcommunity-microservices/app/user/rpc/internal/config"
	"smartcommunity-microservices/app/user/rpc/internal/model"
	"smartcommunity-microservices/app/user/rpc/internal/repository"
	"smartcommunity-microservices/app/user/rpc/internal/service"
	"smartcommunity-microservices/common/db"
	"smartcommunity-microservices/common/mq"
	"smartcommunity-microservices/common/redis"
	"smartcommunity-microservices/common/storage"

	"github.com/minio/minio-go/v7"
	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config          config.Config
	DB              *gorm.DB
	AuthService     *service.AuthService
	UserService     *service.UserService
	AdminService    *service.AdminService
	LoginLogService *service.LoginLogService
	RedisClient     *goredis.Client
	MqClient        *mq.Client
	MinioClient     *minio.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	database, err := db.InitPostgres(c.Postgres)

	if err != nil {
		log.Fatalf("failed to init postgres in user-rpc: %v", err)
	}

	// Run AutoMigrate for PostgreSQL tables
	_ = database.AutoMigrate(
		&model.SysUser{},
		&model.SysRole{},
		&model.SysMenu{},
		&model.SysUserRole{},
		&model.SysRoleMenu{},
		&model.SysPermission{},
		&model.SysRolePermission{},
		&model.AdminLoginLog{},
		&model.UserLoginLog{},
		&model.PasswordResetCode{},
	)

	if err := database.Exec(`
		UPDATE sys_user
		SET green_points_total_earned = green_points
		WHERE green_points_total_earned < green_points
	`).Error; err != nil {
		log.Printf("warning: failed to backfill green_points_total_earned in user-rpc: %v", err)
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

	var minioClient *minio.Client
	if c.MinIO.Endpoint != "" {
		minioClient, err = storage.Init(c.MinIO)
		if err != nil {
			log.Printf("warning: minio initialization failed in user-rpc: %v", err)
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
		DB:              database,
		AuthService:     service.NewAuthService(userRepo, roleRepo, logRepo, resetRepo, rdb, c.Jwt.Secret, jwtTTL),
		UserService:     service.NewUserService(userRepo, rdb, mqClient),
		AdminService:    service.NewAdminService(userRepo, roleRepo, rdb),
		LoginLogService: service.NewLoginLogService(logRepo),
		RedisClient:     rdb,
		MqClient:        mqClient,
		MinioClient:     minioClient,
	}
}

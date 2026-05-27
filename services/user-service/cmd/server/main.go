package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"smartcommunity-microservices/pkg/config"
	"smartcommunity-microservices/pkg/db"
	"smartcommunity-microservices/pkg/logger"
	"smartcommunity-microservices/pkg/middleware"
	"smartcommunity-microservices/pkg/minio"
	"smartcommunity-microservices/pkg/nacos"
	"smartcommunity-microservices/pkg/rabbitmq"
	"smartcommunity-microservices/pkg/redis"
	"smartcommunity-microservices/pkg/response"

	"smartcommunity-microservices/services/user-service/internal/handler"
	"smartcommunity-microservices/services/user-service/internal/model"
	"smartcommunity-microservices/services/user-service/internal/repository"
	"smartcommunity-microservices/services/user-service/internal/router"
	"smartcommunity-microservices/services/user-service/internal/service"

	"github.com/gin-gonic/gin"
	miniogo "github.com/minio/minio-go/v7"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}
	if cfg.Service.Name == "" {
		cfg.Service.Name = "user-service"
	}
	if cfg.Service.Port == 0 {
		cfg.Service.Port = 8001
	}

	logr := logger.New(cfg.Service.Name)

	// MySQL
	database, err := db.InitMySQL(cfg.MySQL)
	if err != nil {
		log.Fatalf("init mysql failed: %v", err)
	}

	// AutoMigrate
	if err := database.AutoMigrate(
		&model.SysUser{},
		&model.SysRole{},
		&model.SysMenu{},
		&model.SysPermission{},
		&model.SysUserRole{},
		&model.SysRoleMenu{},
		&model.SysRolePermission{},
		&model.UserLoginLog{},
		&model.AdminLoginLog{},
		&model.PasswordResetCode{},
	); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	// Redis
	rdb, err := redis.Init(cfg.Redis)
	if err != nil {
		log.Fatalf("init redis failed: %v", err)
	}

	// JWT config from YAML (via Raw map)
	jwtSecret := "dev-secret-change-me"
	jwtTTL := 24 * time.Hour
	if v, ok := cfg.Raw["jwt"].(map[string]interface{}); ok {
		if s, ok := v["secret"].(string); ok && s != "" {
			jwtSecret = s
		}
		if t, ok := v["ttl"].(string); ok {
			if d, err := time.ParseDuration(t); err == nil {
				jwtTTL = d
			}
		}
	}

	// Internal token for service-to-service auth
	internalToken := ""
	if v, ok := cfg.Raw["internal_token"].(string); ok {
		internalToken = v
	}

	// Nacos registration (best-effort)
	if err := nacos.RegisterService(context.Background(), cfg.Nacos, cfg.Service.Name, cfg.Service.Host, cfg.Service.Port, map[string]string{"domain": "user"}, cfg.Service.RegisterIP); err != nil {
		logr.Warn("nacos registration skipped", "error", err)
	}

	// RabbitMQ
	var mqClient *rabbitmq.Client
	if cfg.RabbitMQ.URL() != "" {
		var err error
		mqClient, err = rabbitmq.Connect(cfg.RabbitMQ)
		if err != nil {
			logr.Warn("rabbitmq connect failed, events disabled", "error", err)
		} else {
			logr.Info("rabbitmq connected")
			defer mqClient.Close()
		}
	}

	// MinIO (for file cleanup consumer)
	var minioClient *miniogo.Client
	if cfg.MinIO.Endpoint != "" {
		var err error
		minioClient, err = minio.Init(cfg.MinIO)
		if err != nil {
			logr.Warn("minio client init failed, cleanup consumer disabled", "error", err)
		} else {
			logr.Info("minio client initialized for cleanup")
		}
	}

	// Start RabbitMQ cleanup consumer
	if mqClient != nil && minioClient != nil {
		err = mqClient.ConsumeEvents("file.cleanup", func(d amqp.Delivery) {
			var msg struct {
				URL string `json:"url"`
			}
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				logr.Error("failed to unmarshal cleanup message", "error", err)
				_ = d.Nack(false, false) // discard malformed message
				return
			}
			if msg.URL == "" {
				_ = d.Ack(false)
				return
			}

			// Extract key from URL
			u, err := url.Parse(msg.URL)
			if err != nil {
				logr.Error("failed to parse cleanup URL", "url", msg.URL, "error", err)
				_ = d.Ack(false)
				return
			}

			path := strings.TrimPrefix(u.Path, "/")
			bucketPrefix := cfg.MinIO.Bucket + "/"
			if !strings.HasPrefix(path, bucketPrefix) {
				logr.Warn("cleanup URL path does not match expected bucket", "url", msg.URL, "bucket", cfg.MinIO.Bucket)
				_ = d.Ack(false)
				return
			}

			key := strings.TrimPrefix(path, bucketPrefix)
			if key == "" {
				_ = d.Ack(false)
				return
			}

			logr.Info("deleting object from MinIO", "bucket", cfg.MinIO.Bucket, "key", key)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			err = minioClient.RemoveObject(ctx, cfg.MinIO.Bucket, key, miniogo.RemoveObjectOptions{})
			if err != nil {
				logr.Error("failed to remove object from minio", "bucket", cfg.MinIO.Bucket, "key", key, "error", err)
				_ = d.Nack(false, true) // requeue
				return
			}

			_ = d.Ack(false)
		})
		if err != nil {
			logr.Error("failed to start cleanup consumer", "error", err)
		} else {
			logr.Info("started cleanup consumer for file.cleanup queue")
		}
	}

	// Dependency injection
	userRepo := repository.NewUserRepo(database)
	roleRepo := repository.NewRoleRepo(database)
	logRepo := repository.NewLoginLogRepo(database)
	resetRepo := repository.NewPasswordResetRepo(database)

	authService := service.NewAuthService(userRepo, roleRepo, logRepo, resetRepo, rdb, jwtSecret, jwtTTL)
	userService := service.NewUserService(userRepo, rdb, mqClient)
	adminService := service.NewAdminService(userRepo, roleRepo, rdb)
	loginLogService := service.NewLoginLogService(logRepo)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	adminHandler := handler.NewAdminHandler(adminService, loginLogService)

	// Router
	r := gin.New()
	r.Use(middleware.RequestID(), middleware.Logger(logr), middleware.Recovery(logr), middleware.CORS())
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"service": cfg.Service.Name, "status": "ok"})
	})

	router.SetupRoutes(r, router.RouterConfig{
		AuthHandler:   authHandler,
		UserHandler:   userHandler,
		AdminHandler:  adminHandler,
		JWTSecret:     jwtSecret,
		RedisClient:   rdb,
		InternalToken: internalToken,
	})

	addr := fmt.Sprintf("%s:%d", cfg.Service.Host, cfg.Service.Port)
	logr.Info("starting service", "addr", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("run server failed: %v", err)
	}
}

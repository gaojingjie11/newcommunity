package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"smartcommunity-microservices/pkg/config"
	"smartcommunity-microservices/pkg/db"
	"smartcommunity-microservices/pkg/logger"
	"smartcommunity-microservices/pkg/middleware"
	"smartcommunity-microservices/pkg/nacos"
	"smartcommunity-microservices/pkg/rabbitmq"
	"smartcommunity-microservices/pkg/redis"
	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/community-service/internal/client"
	"smartcommunity-microservices/services/community-service/internal/handler"
	"smartcommunity-microservices/services/community-service/internal/model"
	"smartcommunity-microservices/services/community-service/internal/repository"
	"smartcommunity-microservices/services/community-service/internal/router"
	"smartcommunity-microservices/services/community-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}
	if cfg.Service.Name == "" {
		cfg.Service.Name = "community-service"
	}
	if cfg.Service.Port == 0 {
		cfg.Service.Port = 8003
	}

	logr := logger.New(cfg.Service.Name)

	database, err := db.InitMySQL(cfg.MySQL)
	if err != nil {
		log.Fatalf("init mysql failed: %v", err)
	}
	if err := database.AutoMigrate(
		&model.Notice{},
		&model.NoticeViewLog{},
		&model.Visitor{},
		&model.ParkingSpace{},
		&model.UserParkingBinding{},
		&model.PropertyFee{},
		&model.PropertyFeePayment{},
		&model.WorkOrder{},
		&model.WorkorderLog{},
		&model.AIReport{},
		&model.CommunityMessage{},
	); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	rdb, err := redis.Init(cfg.Redis)
	if err != nil {
		log.Fatalf("init redis failed: %v", err)
	}

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
	_ = jwtTTL

	// Mall-service internal API client (for wallet operations)
	var mallClient *client.MallClient
	if v, ok := cfg.Raw["mall_internal"].(map[string]interface{}); ok {
		baseURL, _ := v["base_url"].(string)
		token, _ := v["internal_token"].(string)
		if baseURL != "" && token != "" {
			mallClient = client.NewMallClient(baseURL, token)
			logr.Info("mall-service internal client configured", "base_url", baseURL)
		}
	}

	if err := nacos.RegisterService(context.Background(), cfg.Nacos, cfg.Service.Name, cfg.Service.Host, cfg.Service.Port, map[string]string{"domain": "community"}, cfg.Service.RegisterIP); err != nil {
		logr.Warn("nacos registration skipped", "error", err)
	}

	// RabbitMQ (best-effort, for workorder event bus)
	mq, err := rabbitmq.Connect(cfg.RabbitMQ)
	if err != nil {
		logr.Warn("rabbitmq unavailable; events will be skipped", "error", err)
	}
	defer func() { _ = mq.Close() }()

	noticeRepo := repository.NewNoticeRepo(database)
	visitorRepo := repository.NewVisitorRepo(database)
	parkingRepo := repository.NewParkingRepo(database)
	propertyFeeRepo := repository.NewPropertyFeeRepo(database)
	permissionRepo := repository.NewPermissionRepo(database)

	noticeSvc := service.NewNoticeService(noticeRepo)
	visitorSvc := service.NewVisitorService(visitorRepo)
	parkingSvc := service.NewParkingService(parkingRepo)
	propertyFeeSvc := service.NewPropertyFeeService(propertyFeeRepo, mallClient)

	noticeHandler := handler.NewNoticeHandler(noticeSvc)
	visitorHandler := handler.NewVisitorHandler(visitorSvc)
	parkingHandler := handler.NewParkingHandler(parkingSvc)
	propertyFeeHandler := handler.NewPropertyFeeHandler(propertyFeeSvc)
	permProvider := handler.NewPermissionProvider(permissionRepo)

	// Workorder DI
	workorderRepo := repository.NewWorkorderRepo(database)
	eventBus := service.NewEventBus(mq, logr)
	workorderSvc := service.NewWorkorderService(workorderRepo, eventBus)
	workorderHandler := handler.NewWorkorderHandler(workorderSvc)

	// Statistics DI
	statsRepo := repository.NewStatsRepo(database)
	statsSvc := service.NewStatsService(statsRepo, rdb, logr)
	statsHandler := handler.NewStatsHandler(statsSvc)

	// AI Report DI
	reportRepo := repository.NewReportRepo(database)
	aiSvc := service.NewAIService(cfg.Agent.LLMAPIKey, cfg.Agent.LLMBaseURL, cfg.Agent.LLMModel)
	reportSvc := service.NewReportService(reportRepo, aiSvc)
	reportHandler := handler.NewReportHandler(reportSvc)

	// Community Chat DI
	messageRepo := repository.NewMessageRepo(database)
	messageSvc := service.NewMessageService(messageRepo)
	messageHandler := handler.NewMessageHandler(messageSvc)

	r := gin.New()
	r.Use(middleware.RequestID(), middleware.Logger(logr), middleware.Recovery(logr), middleware.CORS())
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"service": cfg.Service.Name, "status": "ok"})
	})
	router.SetupRoutes(r, router.RouterConfig{
		NoticeHandler:      noticeHandler,
		VisitorHandler:     visitorHandler,
		ParkingHandler:     parkingHandler,
		PropertyFeeHandler: propertyFeeHandler,
		WorkorderHandler:   workorderHandler,
		StatsHandler:       statsHandler,
		ReportHandler:      reportHandler,
		MessageHandler:     messageHandler,
		PermProvider:       permProvider,
		JWTSecret:          jwtSecret,
		RedisClient:        rdb,
	})

	addr := fmt.Sprintf("%s:%d", cfg.Service.Host, cfg.Service.Port)
	logr.Info("starting service", "addr", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("run server failed: %v", err)
	}
}

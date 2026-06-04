package svc

import (
	"log"

	"smartcommunity-microservices/app/stats/rpc/internal/config"
	"smartcommunity-microservices/app/stats/rpc/internal/repository"
	"smartcommunity-microservices/app/stats/rpc/internal/service"
	"smartcommunity-microservices/common/db"
	"smartcommunity-microservices/common/logger"
	"smartcommunity-microservices/common/redis"

	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config     config.Config
	DB         *gorm.DB
	Redis      *goredis.Client
	StatsRepo  *repository.StatsRepo
	ReportRepo *repository.ReportRepo
	StatsSvc   *service.StatsService
	ReportSvc  *service.ReportService
}

func NewServiceContext(c config.Config) *ServiceContext {
	logr := logger.New(c.Name)

	// Initialize MySQL
	database, err := db.InitMySQL(c.MySQL)
	if err != nil {
		log.Fatalf("failed to init mysql in stats-rpc: %v", err)
	}

	// Initialize Redis
	rdb, err := redis.Init(c.BizRedis)
	if err != nil {
		log.Fatalf("failed to init redis in stats-rpc: %v", err)
	}

	// Initialize Repositories
	statsRepo := repository.NewStatsRepo(database)
	reportRepo := repository.NewReportRepo(database)

	// Initialize AI Service
	aiSvc := service.NewAIService(
		c.Agent.LlmApiKey,
		c.Agent.LlmBaseUrl,
		c.Agent.LlmModel,
	)

	// Initialize Services
	statsSvc := service.NewStatsService(statsRepo, rdb, logr)
	reportSvc := service.NewReportService(reportRepo, aiSvc)

	return &ServiceContext{
		Config:     c,
		DB:         database,
		Redis:      rdb,
		StatsRepo:  statsRepo,
		ReportRepo: reportRepo,
		StatsSvc:   statsSvc,
		ReportSvc:  reportSvc,
	}
}

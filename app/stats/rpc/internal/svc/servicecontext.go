package svc

import (
	"log"

	"smartcommunity-microservices/app/stats/rpc/internal/config"
	"smartcommunity-microservices/app/stats/rpc/internal/repository"
	"smartcommunity-microservices/app/stats/rpc/internal/service"
	"smartcommunity-microservices/app/stats/rpc/internal/model"
	"smartcommunity-microservices/common/db"
	"smartcommunity-microservices/common/logger"
	"smartcommunity-microservices/common/mq"
	"smartcommunity-microservices/common/redis"
	"smartcommunity-microservices/common/storage"

	"github.com/minio/minio-go/v7"
	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	DB          *gorm.DB
	Redis       *goredis.Client
	MQ          *mq.Client
	MinioClient *minio.Client
	StatsRepo   *repository.StatsRepo
	ReportRepo  *repository.ReportRepo
	StatsSvc    *service.StatsService
	ReportSvc   *service.ReportService
	Scheduler   *service.ReportScheduler
}

func NewServiceContext(c config.Config) *ServiceContext {
	logr := logger.New(c.Name)

	// Initialize Postgres
	database, err := db.InitPostgres(c.Postgres)
	if err != nil {
		log.Fatalf("failed to init postgres in stats-rpc: %v", err)
	}

	// Run AutoMigrate for PostgreSQL tables
	_ = database.AutoMigrate(
		&model.AIReport{},
	)

	// Initialize Redis
	rdb, err := redis.Init(c.BizRedis)
	if err != nil {
		log.Fatalf("failed to init redis in stats-rpc: %v", err)
	}

	// Initialize MQ
	var mqClient *mq.Client
	if c.RabbitMQ.URL() != "" {
		mqClient, err = mq.Connect(c.RabbitMQ)
		if err != nil {
			logr.Warn("rabbitmq connect failed in stats-rpc, async reports disabled", "error", err)
		} else {
			logr.Info("rabbitmq connected in stats-rpc")
		}
	}

	// Initialize MinIO
	var minioClient *minio.Client
	if c.MinIO.Endpoint != "" {
		minioClient, err = storage.Init(c.MinIO)
		if err != nil {
			logr.Warn("minio initialization failed in stats-rpc", "error", err)
		} else {
			logr.Info("minio initialized in stats-rpc")
		}
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

	// Initialize and start automatic nightly scheduler
	scheduler := service.NewReportScheduler(reportSvc, minioClient, c.MinIO.Bucket)
	scheduler.Start()

	return &ServiceContext{
		Config:      c,
		DB:          database,
		Redis:       rdb,
		MQ:          mqClient,
		MinioClient: minioClient,
		StatsRepo:   statsRepo,
		ReportRepo:  reportRepo,
		StatsSvc:    statsSvc,
		ReportSvc:   reportSvc,
		Scheduler:   scheduler,
	}
}



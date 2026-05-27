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

	"smartcommunity-microservices/services/mall-service/internal/handler"
	"smartcommunity-microservices/services/mall-service/internal/model"
	"smartcommunity-microservices/services/mall-service/internal/repository"
	"smartcommunity-microservices/services/mall-service/internal/router"
	"smartcommunity-microservices/services/mall-service/internal/service"

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
		cfg.Service.Name = "mall-service"
	}
	if cfg.Service.Port == 0 {
		cfg.Service.Port = 8002
	}

	logr := logger.New(cfg.Service.Name)

	// MySQL
	database, err := db.InitMySQL(cfg.MySQL)
	if err != nil {
		log.Fatalf("init mysql failed: %v", err)
	}

	// AutoMigrate
	if err := database.AutoMigrate(
		&model.Product{},
		&model.ProductCategory{},
		&model.Promotion{},
		&model.PromotionProduct{},
		&model.Cart{},
		&model.Order{},
		&model.OrderItem{},
		&model.Store{},
		&model.StoreProduct{},
		&model.Favorite{},
		&model.ProductComment{},
		&model.Wallet{},
		&model.WalletTransaction{},
		&model.ServiceArea{},
		&model.PaymentRecord{},
		&model.ProductViewLog{},
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
	_ = jwtTTL // used by JWTAuth middleware

	// Internal token for service-to-service auth
	internalToken := ""
	if v, ok := cfg.Raw["internal_token"].(string); ok {
		internalToken = v
	}

	// RabbitMQ (best-effort — service degrades gracefully without it)
	var mqClient *rabbitmq.Client
	if cfg.RabbitMQ.URL() != "" {
		mqClient, err = rabbitmq.Connect(cfg.RabbitMQ)
		if err != nil {
			logr.Warn("rabbitmq connect failed, events disabled", "error", err)
		} else {
			logr.Info("rabbitmq connected")
		}
	}

	// Nacos registration (best-effort)
	if err := nacos.RegisterService(context.Background(), cfg.Nacos, cfg.Service.Name, cfg.Service.Host, cfg.Service.Port, map[string]string{"domain": "mall"}, cfg.Service.RegisterIP); err != nil {
		logr.Warn("nacos registration skipped", "error", err)
	}

	// Dependency injection — repositories
	productRepo := repository.NewProductRepo(database)
	categoryRepo := repository.NewCategoryRepo(database)
	promotionRepo := repository.NewPromotionRepo(database)
	cartRepo := repository.NewCartRepo(database)
	orderRepo := repository.NewOrderRepo(database)
	storeRepo := repository.NewStoreRepo(database)
	storeProductRepo := repository.NewStoreProductRepo(database)
	favoriteRepo := repository.NewFavoriteRepo(database)
	commentRepo := repository.NewCommentRepo(database)
	walletRepo := repository.NewWalletRepo(database)
	serviceAreaRepo := repository.NewServiceAreaRepo(database)
	permissionRepo := repository.NewPermissionRepo(database)
	paymentRepo := repository.NewPaymentRepo(database)
	viewLogRepo := repository.NewViewLogRepo(database)
	userRepo := repository.NewUserRepo(database)

	// Event bus (wraps RabbitMQ, nil-safe)
	eventBus := service.NewEventBus(mqClient, logr)

	// Dependency injection — services
	productSvc := service.NewProductService(productRepo)
	categorySvc := service.NewCategoryService(categoryRepo)
	promoSvc := service.NewPromotionService(promotionRepo)
	cartSvc := service.NewCartService(cartRepo, productRepo)
	orderSvc := service.NewOrderService(database, orderRepo, cartRepo, productRepo, storeRepo, storeProductRepo, walletRepo, eventBus, rdb)
	storeSvc := service.NewStoreService(storeRepo, storeProductRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo)
	commentSvc := service.NewCommentService(commentRepo, productRepo)
	walletSvc := service.NewWalletService(database, walletRepo)
	serviceAreaSvc := service.NewServiceAreaService(serviceAreaRepo)
	paymentSvc := service.NewPaymentService(database, orderRepo, storeProductRepo, walletRepo, paymentRepo, userRepo, productRepo, eventBus, logr)

	// Order timeout service (polling fallback)
	timeoutSvc := service.NewOrderTimeoutService(database, orderRepo, storeProductRepo, productRepo, eventBus, logr)
	timeoutSvc.Start()
	defer timeoutSvc.Stop()

	// Permission provider (RBAC)
	permProvider := handler.NewPermissionProvider(permissionRepo)

	// Dependency injection — handlers
	productHandler := handler.NewProductHandler(productSvc, viewLogRepo)
	cartHandler := handler.NewCartHandler(cartSvc)
	orderHandler := handler.NewOrderHandler(orderSvc)
	paymentHandler := handler.NewPaymentHandler(paymentSvc)
	storeHandler := handler.NewStoreHandler(storeSvc)
	favoriteHandler := handler.NewFavoriteHandler(favoriteSvc)
	commentHandler := handler.NewCommentHandler(commentSvc)
	walletHandler := handler.NewWalletHandler(walletSvc)
	promotionHandler := handler.NewPromotionHandler(promoSvc)
	categoryHandler := handler.NewCategoryHandler(categorySvc)
	serviceAreaHandler := handler.NewServiceAreaHandler(serviceAreaSvc)
	adminOrderHandler := handler.NewAdminOrderHandler(orderSvc, productSvc, storeSvc, promoSvc)
	internalHandler := handler.NewInternalHandler(orderSvc, paymentSvc, walletSvc)

	// Router
	r := gin.New()
	r.Use(middleware.RequestID(), middleware.Logger(logr), middleware.Recovery(logr), middleware.CORS())
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"service": cfg.Service.Name, "status": "ok"})
	})

	router.SetupRoutes(r, router.RouterConfig{
		ProductHandler:     productHandler,
		CartHandler:        cartHandler,
		OrderHandler:       orderHandler,
		PaymentHandler:     paymentHandler,
		StoreHandler:       storeHandler,
		FavoriteHandler:    favoriteHandler,
		CommentHandler:     commentHandler,
		WalletHandler:      walletHandler,
		PromotionHandler:   promotionHandler,
		CategoryHandler:    categoryHandler,
		ServiceAreaHandler: serviceAreaHandler,
		AdminOrderHandler:  adminOrderHandler,
		InternalHandler:    internalHandler,
		PermProvider:       permProvider,
		JWTSecret:          jwtSecret,
		RedisClient:        rdb,
		InternalToken:      internalToken,
	})

	addr := fmt.Sprintf("%s:%d", cfg.Service.Host, cfg.Service.Port)
	logr.Info("starting service", "addr", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("run server failed: %v", err)
	}
}

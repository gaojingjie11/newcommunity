package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"smartcommunity-microservices/app/mall/rpc/internal/config"
	"smartcommunity-microservices/app/mall/rpc/internal/repository"
	"smartcommunity-microservices/app/mall/rpc/internal/service"
	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/common/db"
	"smartcommunity-microservices/common/logger"
	"smartcommunity-microservices/common/mq"
	"smartcommunity-microservices/common/redis"

	goredis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	DB     *gorm.DB
	Redis  *goredis.Client
	MQ     *mq.Client

	ProductRepo      *repository.ProductRepo
	CategoryRepo     *repository.CategoryRepo
	CartRepo         *repository.CartRepo
	OrderRepo        *repository.OrderRepo
	StoreRepo        *repository.StoreRepo
	StoreProductRepo *repository.StoreProductRepo
	FavoriteRepo     *repository.FavoriteRepo
	CommentRepo      *repository.CommentRepo
	WalletRepo       *repository.WalletRepo
	ServiceAreaRepo  *repository.ServiceAreaRepo
	PaymentRepo      *repository.PaymentRepo
	UserRepo         *repository.UserRepo

	ProductSvc     *service.ProductService
	CategorySvc    *service.CategoryService
	CartSvc        *service.CartService
	OrderSvc       *service.OrderService
	StoreSvc       *service.StoreService
	FavoriteSvc    *service.FavoriteService
	CommentSvc     *service.CommentService
	WalletSvc      *service.WalletService
	ServiceAreaSvc *service.ServiceAreaService
	PaymentSvc     *service.PaymentService
	AlipaySvc      *service.AlipayService
	EventBus       *service.EventBus
	TimeoutSvc     *service.OrderTimeoutService
}

func NewServiceContext(c config.Config) *ServiceContext {
	logr := logger.New(c.Name)

	// Init Postgres
	database, err := db.InitPostgres(c.Postgres)
	if err != nil {
		log.Fatalf("init postgres failed: %v", err)
	}

	// Run AutoMigrate for PostgreSQL tables
	_ = database.AutoMigrate(
		&model.Cart{},
		&model.ProductComment{},
		&model.UserProfile{},
		&model.Favorite{},
		&model.Order{},
		&model.OrderItem{},
		&model.PaymentRecord{},
		&model.Product{},
		&model.ProductCategory{},
		&model.ServiceArea{},
		&model.Store{},
		&model.StoreProduct{},
		&model.UserStore{},
		&model.Wallet{},
		&model.WalletTransaction{},
	)

	// Init Redis
	rdb, err := redis.Init(c.BizRedis)
	if err != nil {
		log.Fatalf("init redis failed: %v", err)
	}

	// Init MQ
	var mqClient *mq.Client
	if c.RabbitMQ.URL() != "" {
		mqClient, err = mq.Connect(c.RabbitMQ)
		if err != nil {
			logr.Warn("rabbitmq connect failed, events disabled", "error", err)
		} else {
			logr.Info("rabbitmq connected")
		}
	}

	// Event Bus
	eventBus := service.NewEventBus(mqClient, logr)

	// Repositories
	productRepo := repository.NewProductRepo(database)
	categoryRepo := repository.NewCategoryRepo(database)
	cartRepo := repository.NewCartRepo(database)
	orderRepo := repository.NewOrderRepo(database)
	storeRepo := repository.NewStoreRepo(database)
	storeProductRepo := repository.NewStoreProductRepo(database)
	favoriteRepo := repository.NewFavoriteRepo(database)
	commentRepo := repository.NewCommentRepo(database)
	walletRepo := repository.NewWalletRepo(database)
	serviceAreaRepo := repository.NewServiceAreaRepo(database)
	paymentRepo := repository.NewPaymentRepo(database)
	userRepo := repository.NewUserRepo(database)

	// Services
	productSvc := service.NewProductService(productRepo, rdb, eventBus)
	categorySvc := service.NewCategoryService(categoryRepo, rdb)
	cartSvc := service.NewCartService(cartRepo, productRepo)
	orderSvc := service.NewOrderService(database, orderRepo, cartRepo, productRepo, storeRepo, storeProductRepo, walletRepo, eventBus, rdb)
	storeSvc := service.NewStoreService(storeRepo, storeProductRepo)
	favoriteSvc := service.NewFavoriteService(favoriteRepo)
	commentSvc := service.NewCommentService(commentRepo, productRepo)
	walletSvc := service.NewWalletService(database, walletRepo)
	serviceAreaSvc := service.NewServiceAreaService(serviceAreaRepo)
	paymentSvc := service.NewPaymentService(database, orderRepo, storeProductRepo, walletRepo, paymentRepo, userRepo, productRepo, eventBus, logr)

	alipaySvc, err := service.NewAlipayService(c)
	if err != nil {
		logr.Warn("alipay init failed, alipay disabled", "error", err)
	}

	// Timeout service (polling fallback)
	timeoutSvc := service.NewOrderTimeoutService(database, orderRepo, storeProductRepo, productRepo, eventBus, logr)
	timeoutSvc.Start()

	// Asynchronously warm up the product details cache in Redis
	go warmUpProductCache(database, rdb)

	return &ServiceContext{
		Config:           c,
		DB:               database,
		Redis:            rdb,
		MQ:               mqClient,
		ProductRepo:      productRepo,
		CategoryRepo:     categoryRepo,
		CartRepo:         cartRepo,
		OrderRepo:        orderRepo,
		StoreRepo:        storeRepo,
		StoreProductRepo: storeProductRepo,
		FavoriteRepo:     favoriteRepo,
		CommentRepo:      commentRepo,
		WalletRepo:       walletRepo,
		ServiceAreaRepo:  serviceAreaRepo,
		PaymentRepo:      paymentRepo,
		UserRepo:         userRepo,
		ProductSvc:       productSvc,
		CategorySvc:      categorySvc,
		CartSvc:          cartSvc,
		OrderSvc:         orderSvc,
		StoreSvc:         storeSvc,
		FavoriteSvc:      favoriteSvc,
		CommentSvc:       commentSvc,
		WalletSvc:        walletSvc,
		ServiceAreaSvc:   serviceAreaSvc,
		PaymentSvc:       paymentSvc,
		AlipaySvc:        alipaySvc,
		EventBus:         eventBus,
		TimeoutSvc:       timeoutSvc,
	}
}

func warmUpProductCache(db *gorm.DB, rdb *goredis.Client) {
	var products []model.Product
	// Query only active products (status = 1)
	if err := db.Where("status = ?", 1).Find(&products).Error; err != nil {
		log.Printf("[Cache Warm-Up] Failed to query products from DB: %v", err)
		return
	}

	ctx := context.Background()
	count := 0
	for _, p := range products {
		cacheKey := fmt.Sprintf("mall:product:detail:%d", p.ID)
		productJSON, err := json.Marshal(p)
		if err == nil {
			// Set a randomized TTL (2 hours + up to 10 minutes random jitter)
			ttl := 7200 + rand.Intn(600)
			_ = rdb.Set(ctx, cacheKey, string(productJSON), time.Duration(ttl)*time.Second).Err()
			count++
		}
	}
	log.Printf("[Cache Warm-Up] Completed. Loaded %d/%d products into Redis successfully.", count, len(products))
}

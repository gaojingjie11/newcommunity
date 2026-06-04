package svc

import (
	"log"

	"smartcommunity-microservices/app/mall/rpc/internal/config"
	"smartcommunity-microservices/app/mall/rpc/internal/repository"
	"smartcommunity-microservices/app/mall/rpc/internal/service"
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
	ViewLogRepo      *repository.ViewLogRepo
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

	// Init MySQL
	database, err := db.InitMySQL(c.MySQL)
	if err != nil {
		log.Fatalf("init mysql failed: %v", err)
	}

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
	viewLogRepo := repository.NewViewLogRepo(database)
	userRepo := repository.NewUserRepo(database)

	// Services
	productSvc := service.NewProductService(productRepo)
	categorySvc := service.NewCategoryService(categoryRepo)
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

	return &ServiceContext{
		Config:           c,
		DB:               database,
		Redis:            rdb,
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
		ViewLogRepo:      viewLogRepo,
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

package router

import (
	"smartcommunity-microservices/pkg/middleware"
	"smartcommunity-microservices/services/mall-service/internal/handler"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
)

type RouterConfig struct {
	ProductHandler     *handler.ProductHandler
	CartHandler        *handler.CartHandler
	OrderHandler       *handler.OrderHandler
	PaymentHandler     *handler.PaymentHandler
	StoreHandler       *handler.StoreHandler
	FavoriteHandler    *handler.FavoriteHandler
	CommentHandler     *handler.CommentHandler
	WalletHandler      *handler.WalletHandler
	PromotionHandler   *handler.PromotionHandler
	CategoryHandler    *handler.CategoryHandler
	ServiceAreaHandler *handler.ServiceAreaHandler
	AdminOrderHandler  *handler.AdminOrderHandler
	InternalHandler    *handler.InternalHandler
	PermProvider       middleware.PermissionProvider
	JWTSecret          string
	RedisClient        *goredis.Client
	InternalToken      string // X-Internal-Token for service-to-service calls
}

func SetupRoutes(r *gin.Engine, cfg RouterConfig) {
	authMw := middleware.JWTAuth(cfg.JWTSecret, cfg.RedisClient)
	optionalAuthMw := middleware.OptionalJWTAuth(cfg.JWTSecret, cfg.RedisClient)
	perm := func(code string) gin.HandlerFunc {
		return middleware.RequirePermission(cfg.RedisClient, cfg.PermProvider, code)
	}

	// Public routes with optional JWT enrichment. Anonymous users are allowed,
	// but valid tokens set userID/role for product view logs and future personalization.
	publicMall := r.Group("/api/mall")
	publicMall.Use(optionalAuthMw)
	{
		publicMall.GET("/products", cfg.ProductHandler.List)                     // MALL-001/004
		publicMall.GET("/products/search", cfg.ProductHandler.Search)            // MALL-002
		publicMall.GET("/products/promotions", cfg.ProductHandler.GetPromotions) // MALL-003
		publicMall.GET("/stores", cfg.StoreHandler.List)                         // MALL-006
		publicMall.GET("/stores/:id", cfg.StoreHandler.GetDetail)
		publicMall.GET("/promotions", cfg.PromotionHandler.List)
		publicMall.GET("/promotions/:id", cfg.PromotionHandler.GetDetail)
		publicMall.GET("/categories", cfg.CategoryHandler.List)
		publicMall.GET("/categories/:id", cfg.CategoryHandler.GetDetail)
		publicMall.GET("/service-areas", cfg.ServiceAreaHandler.List)
	}

	// Authenticated user routes
	userMall := r.Group("/api/mall")
	userMall.Use(authMw)
	{
		// Cart (MALL-007/010/011)
		userMall.POST("/cart/items", cfg.CartHandler.Add)
		userMall.DELETE("/cart/items/:id", cfg.CartHandler.Remove)
		userMall.PUT("/cart/items/:id", cfg.CartHandler.UpdateQuantity)
		userMall.GET("/cart/items", cfg.CartHandler.List)

		// Favorites (MALL-008/009/018)
		userMall.POST("/favorites", cfg.FavoriteHandler.Add)
		userMall.DELETE("/favorites/:product_id", cfg.FavoriteHandler.Remove)
		userMall.GET("/favorites", cfg.FavoriteHandler.List)
		userMall.GET("/favorites/check/:product_id", cfg.FavoriteHandler.Check)

		// Product detail requires login by product policy (MALL-005)
		userMall.GET("/products/:id", cfg.ProductHandler.GetDetail)

		// Product comments are tied to the authenticated product detail flow.
		userMall.GET("/comments", cfg.CommentHandler.List)
		userMall.POST("/comments", cfg.CommentHandler.Create)

		// Orders (MALL-012~017)
		userMall.POST("/orders", cfg.OrderHandler.Create)
		userMall.GET("/orders/available-stores", cfg.OrderHandler.ListAvailableStores)
		userMall.POST("/orders/:id/pay", cfg.PaymentHandler.Pay)
		userMall.POST("/orders/:id/cancel", cfg.OrderHandler.Cancel)
		userMall.POST("/orders/:id/receive", cfg.OrderHandler.Receive)
		userMall.GET("/orders/:id", cfg.OrderHandler.GetDetail)
		userMall.GET("/orders/:id/payment-status", cfg.PaymentHandler.GetPaymentStatus)
		userMall.GET("/orders", cfg.OrderHandler.List)

		// Wallet (MALL-019~021)
		userMall.POST("/wallet/recharge", cfg.WalletHandler.Recharge)
		userMall.POST("/wallet/transfer", cfg.WalletHandler.Transfer)
		userMall.GET("/wallet/balance", cfg.WalletHandler.GetBalance)
		userMall.GET("/wallet/transactions", cfg.WalletHandler.ListTransactions)
	}

	// Admin routes (JWT + RequirePermission)
	admin := r.Group("/api/admin/mall")
	admin.Use(authMw)
	{
		// Categories (ADMIN-MALL-005)
		admin.POST("/categories", perm("mall:category:create"), cfg.CategoryHandler.Create)
		admin.PUT("/categories/:id", perm("mall:category:update"), cfg.CategoryHandler.Update)
		admin.DELETE("/categories/:id", perm("mall:category:delete"), cfg.CategoryHandler.Delete)

		// Products (ADMIN-MALL-006)
		admin.GET("/products", perm("mall:product:list"), cfg.AdminOrderHandler.ListProducts)
		admin.POST("/products", perm("mall:product:create"), cfg.AdminOrderHandler.CreateProduct)
		admin.PUT("/products/:id", perm("mall:product:update"), cfg.AdminOrderHandler.UpdateProduct)
		admin.DELETE("/products/:id", perm("mall:product:delete"), cfg.AdminOrderHandler.DeleteProduct)

		// Promotions (ADMIN-MALL-007)
		admin.POST("/promotions", perm("mall:promotion:create"), cfg.AdminOrderHandler.CreatePromotion)
		admin.PUT("/promotions/:id", perm("mall:promotion:update"), cfg.AdminOrderHandler.UpdatePromotion)
		admin.DELETE("/promotions/:id", perm("mall:promotion:delete"), cfg.AdminOrderHandler.DeletePromotion)
		admin.POST("/promotions/:id/products", perm("mall:promotion:bind_product"), cfg.AdminOrderHandler.BindPromotionProducts)

		// Service areas (ADMIN-MALL-008)
		admin.POST("/service-areas", perm("mall:service_area:create"), cfg.ServiceAreaHandler.Create)
		admin.PUT("/service-areas/:id", perm("mall:service_area:update"), cfg.ServiceAreaHandler.Update)
		admin.DELETE("/service-areas/:id", perm("mall:service_area:delete"), cfg.ServiceAreaHandler.Delete)

		// Stores (ADMIN-MALL-009)
		admin.POST("/stores", perm("mall:store:create"), cfg.AdminOrderHandler.CreateStore)
		admin.PUT("/stores/:id", perm("mall:store:update"), cfg.AdminOrderHandler.UpdateStore)
		admin.DELETE("/stores/:id", perm("mall:store:delete"), cfg.AdminOrderHandler.DeleteStore)

		// Store products (ADMIN-MALL-010)
		admin.POST("/store-products", perm("mall:store_product:bind"), cfg.AdminOrderHandler.BindStoreProduct)
		admin.DELETE("/store-products", perm("mall:store_product:unbind"), cfg.AdminOrderHandler.UnbindStoreProduct)
		admin.PUT("/store-products/status", perm("mall:store_product:status"), cfg.AdminOrderHandler.UpdateStoreProductStatus)
		admin.PUT("/store-products/stock", perm("mall:store_product:stock"), cfg.AdminOrderHandler.UpdateStoreProductStock)
		admin.GET("/store-products/:store_id", perm("mall:store_product:list"), cfg.AdminOrderHandler.ListStoreProducts)

		// Orders (ADMIN-MALL-011)
		admin.GET("/orders", perm("mall:order:list"), cfg.AdminOrderHandler.ListOrders)
		admin.POST("/orders/:id/ship", perm("mall:order:ship"), cfg.AdminOrderHandler.ShipOrder)
		admin.POST("/orders/:id/cancel", perm("mall:order:cancel"), cfg.AdminOrderHandler.CancelOrder)
	}

	// Internal routes (service-to-service, X-Internal-Token)
	if cfg.InternalToken != "" && cfg.InternalHandler != nil {
		internal := r.Group("/api/internal/mall")
		internal.Use(middleware.InternalToken(cfg.InternalToken))
		{
			internal.POST("/orders/:id/cancel", cfg.InternalHandler.CancelExpiredOrder)
			internal.GET("/orders/:id/payment-status", cfg.InternalHandler.GetPaymentStatus)
			internal.POST("/wallet/debit", cfg.InternalHandler.DebitWallet)
		}
	}
}

package router

import (
	"smartcommunity-microservices/pkg/middleware"
	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/community-service/internal/handler"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
)

type RouterConfig struct {
	NoticeHandler      *handler.NoticeHandler
	VisitorHandler     *handler.VisitorHandler
	ParkingHandler     *handler.ParkingHandler
	PropertyFeeHandler *handler.PropertyFeeHandler
	WorkorderHandler   *handler.WorkorderHandler
	StatsHandler       *handler.StatsHandler
	ReportHandler      *handler.ReportHandler
	MessageHandler     *handler.MessageHandler
	PermProvider       middleware.PermissionProvider
	JWTSecret          string
	RedisClient        *goredis.Client
}

func SetupRoutes(r *gin.Engine, cfg RouterConfig) {
	authMw := middleware.JWTAuth(cfg.JWTSecret, cfg.RedisClient)
	perm := func(code string) gin.HandlerFunc {
		return middleware.RequirePermission(cfg.RedisClient, cfg.PermProvider, code)
	}

	public := r.Group("/api/community")
	{
		public.GET("/ping", func(c *gin.Context) {
			response.Success(c, gin.H{"service": "community-service", "pong": true})
		})
		public.GET("/notices", cfg.NoticeHandler.List)       // COMM-001
		public.GET("/notices/:id", cfg.NoticeHandler.Detail) // COMM-001 detail + view_count
	}

	user := r.Group("/api/community")
	user.Use(authMw)
	{
		user.POST("/notices/:id/read", cfg.NoticeHandler.MarkRead) // ADMIN-COMM-002 support

		user.POST("/visitors", cfg.VisitorHandler.Create) // COMM-002
		user.GET("/visitors", cfg.VisitorHandler.MyList)

		user.GET("/parking-spaces/my", cfg.ParkingHandler.MyBindings)       // COMM-003
		user.PUT("/parking-spaces/:id/plate", cfg.ParkingHandler.BindPlate) // COMM-004

		user.GET("/property-fees", cfg.PropertyFeeHandler.MyFees)
		user.POST("/property-fees/:id/pay", cfg.PropertyFeeHandler.Pay)        // COMM-007
		user.GET("/property-fees/payments", cfg.PropertyFeeHandler.MyPayments) // ADMIN-COMM-007 support

		user.GET("/messages", cfg.MessageHandler.List) // community chat
		user.POST("/message", cfg.MessageHandler.Send)
	}

	admin := r.Group("/api/admin/community")
	admin.Use(authMw)
	{
		admin.GET("/notices", perm("community:notice:list"), cfg.NoticeHandler.AdminList)
		admin.POST("/notices", perm("community:notice:create"), cfg.NoticeHandler.Create) // ADMIN-COMM-001
		admin.DELETE("/notices/:id", perm("community:notice:delete"), cfg.NoticeHandler.Delete)
		admin.GET("/notices/:id/views", perm("community:notice:views"), cfg.NoticeHandler.Views) // ADMIN-COMM-002

		admin.GET("/visitors", perm("community:visitor:list"), cfg.VisitorHandler.AdminList)
		admin.POST("/visitors/:id/audit", perm("community:visitor:audit"), cfg.VisitorHandler.Audit) // ADMIN-COMM-003

		admin.GET("/parking-spaces", perm("community:parking:list"), cfg.ParkingHandler.AdminList)
		admin.POST("/parking-spaces", perm("community:parking:create"), cfg.ParkingHandler.Create)
		admin.POST("/parking-spaces/:id/assign", perm("community:parking:assign"), cfg.ParkingHandler.Assign)
		admin.GET("/parking-spaces/statistics", perm("community:parking:statistics"), cfg.ParkingHandler.Stats) // ADMIN-COMM-004

		admin.GET("/property-fees", perm("community:fee:list"), cfg.PropertyFeeHandler.AdminList)
		admin.POST("/property-fees", perm("community:fee:create"), cfg.PropertyFeeHandler.AdminCreate)
		admin.GET("/property-fees/payments", perm("community:fee:payment_list"), cfg.PropertyFeeHandler.AdminPayments) // ADMIN-COMM-007
	}

	// ── Workorder public ──
	wkPublic := r.Group("/api/workorders")
	{
		wkPublic.GET("/ping", func(c *gin.Context) {
			response.Success(c, gin.H{"service": "community-service", "pong": true})
		})
	}

	// ── Workorder user (auth) ──
	wkUser := r.Group("/api/workorders")
	wkUser.Use(authMw)
	{
		wkUser.POST("", cfg.WorkorderHandler.Create)
		wkUser.GET("", cfg.WorkorderHandler.MyList)
		wkUser.GET("/:id/logs", cfg.WorkorderHandler.Logs)
	}

	// ── Workorder admin (auth + perm) ──
	wkAdmin := r.Group("/api/admin/workorders")
	wkAdmin.Use(authMw)
	{
		wkAdmin.GET("", perm("workorder:repair:list"), cfg.WorkorderHandler.AdminList)
		wkAdmin.POST("/:id/process", perm("workorder:repair:process"), cfg.WorkorderHandler.Process)
	}

	// ── Statistics (auth + perm) ──
	stats := r.Group("/api/statistics")
	stats.Use(authMw)
	{
		stats.GET("/products/sales-rank", perm("statistics:product:sales_rank"), cfg.StatsHandler.ProductSalesRank)
		stats.GET("/products/view-rank", perm("statistics:product:view_rank"), cfg.StatsHandler.ProductViewRank)
		stats.GET("/community/overview", perm("statistics:community:overview"), cfg.StatsHandler.CommunityOverview)
		stats.GET("/orders", perm("statistics:order:summary"), cfg.StatsHandler.OrderStats)
		stats.GET("/workorders", perm("statistics:workorder:summary"), cfg.StatsHandler.WorkorderStats)

		stats.POST("/ai-report/generate", perm("statistics:ai_report:generate"), cfg.ReportHandler.GenerateReport)
		stats.GET("/ai-report/latest", perm("statistics:ai_report:read"), cfg.ReportHandler.GetLatestReport)
		stats.GET("/ai-report/list", perm("statistics:ai_report:read"), cfg.ReportHandler.ListReports)
		stats.GET("/ai-report/:id", perm("statistics:ai_report:read"), cfg.ReportHandler.GetReportDetail)
	}
}

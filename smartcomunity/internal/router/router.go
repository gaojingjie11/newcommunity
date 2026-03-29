package router

import (
	"smartcommunity/internal/controller"
	"smartcommunity/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitRouter(r *gin.Engine) {
	uploadHandler := controller.UploadHandler{}
	userHandler := controller.UserHandler{}
	productHandler := controller.ProductHandler{}
	cartHandler := controller.CartHandler{}
	orderHandler := controller.OrderHandler{}
	noticeHandler := controller.NoticeHandler{}
	repairHandler := controller.RepairHandler{}
	financeHandler := controller.FinanceHandler{}
	securityHandler := controller.SecurityHandler{}
	favoriteHandler := controller.FavoriteHandler{}
	storeHandler := controller.StoreHandler{}
	adminHandler := controller.AdminHandler{}
	marketingHandler := controller.MarketingHandler{}
	commentHandler := controller.CommentHandler{}
	aiHandler := controller.AIHandler{}
	greenPointHandler := controller.GreenPointHandler{}
	communityMessageHandler := controller.CommunityMessageHandler{}

	publicAPI := r.Group("/api/v1")
	{
		publicAPI.POST("/register", userHandler.Register)
		publicAPI.POST("/login", userHandler.Login)
		publicAPI.POST("/send_code", userHandler.SendCode)
		publicAPI.POST("/login_code", userHandler.LoginCode)
		publicAPI.POST("/forget_password", userHandler.ForgetPassword)

		publicAPI.GET("/products", productHandler.List)
		publicAPI.GET("/product/:id", productHandler.Detail)
		publicAPI.POST("/cart", productHandler.Create)
		publicAPI.GET("/stores", storeHandler.List)
		publicAPI.GET("/categories", productHandler.GetCategories)
		publicAPI.GET("/dashboard/stats", adminHandler.GetDashboardStats)
		publicAPI.GET("/comments", commentHandler.List)
		publicAPI.GET("/green-points/leaderboard", greenPointHandler.Leaderboard)
		publicAPI.GET("/notices", noticeHandler.List)
		publicAPI.GET("/notice/:id", noticeHandler.Detail)
	}

	private := r.Group("/api/v1")
	private.Use(middleware.JWTAuth())
	{
		private.POST("/cart/add", cartHandler.Add)
		private.GET("/cart/list", cartHandler.List)
		private.DELETE("/cart/:id", cartHandler.Delete)
		private.POST("/cart/:id", cartHandler.Update)

		private.POST("/order/create", orderHandler.Create)
		private.GET("/order/list", orderHandler.List)
		private.GET("/order/detail", orderHandler.Detail)
		private.POST("/order/pay", orderHandler.Pay)
		private.GET("/order/admin/list", middleware.RequireRole("admin", "store"), orderHandler.ListAll)
		private.POST("/order/ship", middleware.RequireRole("admin", "store"), orderHandler.Ship)
		private.POST("/order/receive", orderHandler.Receive)
		private.POST("/order/cancel", orderHandler.Cancel)

		private.POST("/repair/create", repairHandler.Create)
		private.GET("/repair/list", repairHandler.List)

		private.POST("/finance/pay", financeHandler.Pay)
		private.GET("/property/list", financeHandler.ListPropertyFee)
		private.POST("/finance/recharge", financeHandler.Recharge)
		private.POST("/finance/transfer", financeHandler.Transfer)
		private.GET("/finance/transactions", financeHandler.ListTransactions)

		private.POST("/green-points/upload-garbage", greenPointHandler.UploadGarbage)

		private.POST("/marketing/promotion/create", marketingHandler.Create)
		private.GET("/marketing/promotion/list", marketingHandler.List)
		private.DELETE("/marketing/promotion/:id", marketingHandler.Delete)

		private.POST("/store/create", middleware.RequireRole("admin", "store"), storeHandler.Create)
		private.POST("/store/update", middleware.RequireRole("admin", "store"), storeHandler.Update)
		private.DELETE("/store/:id", middleware.RequireRole("admin", "store"), storeHandler.Delete)
		private.POST("/store/bind_product", middleware.RequireRole("admin", "store"), storeHandler.BindProduct)

		private.POST("/product/create", middleware.RequireRole("admin", "store"), productHandler.Create)
		private.POST("/product/update", middleware.RequireRole("admin", "store"), productHandler.Update)
		private.DELETE("/product/:id", middleware.RequireRole("admin", "store"), productHandler.Delete)
		private.GET("/product/rank", productHandler.GetRank)

		private.POST("/visitor/create", securityHandler.CreateVisitor)
		private.GET("/visitor/list", securityHandler.ListVisitor)
		private.GET("/parking/my", securityHandler.MyParking)
		private.POST("/parking/bind", securityHandler.BindCar)
		private.GET("/visitor/admin/list", middleware.RequireRole("admin", "property"), securityHandler.ListAllVisitor)
		private.POST("/visitor/audit", middleware.RequireRole("admin", "property"), securityHandler.AuditVisitor)

		private.POST("/upload", uploadHandler.UploadFile)

		private.POST("/notice/create", middleware.RequireRole("admin", "property"), noticeHandler.Create)
		private.DELETE("/notice/:id", middleware.RequireRole("admin", "property"), noticeHandler.Delete)
		private.POST("/notice/read/:id", noticeHandler.Read)

		private.GET("/repair/admin/list", middleware.RequireRole("admin", "property"), repairHandler.ListAll)
		private.POST("/repair/process", middleware.RequireRole("admin", "property"), repairHandler.Process)

		private.POST("/favorite/add", favoriteHandler.Add)
		private.POST("/favorite/delete", favoriteHandler.Delete)
		private.GET("/favorites", favoriteHandler.List)
		private.GET("/favorite/check", favoriteHandler.Check)

		private.POST("/logout", userHandler.Logout)
		private.POST("/user/update", userHandler.Update)
		private.POST("/user/change_password", userHandler.ChangePassword)
		private.POST("/user/face/register", userHandler.RegisterFace)
		private.GET("/user/info", userHandler.Info)

		private.GET("/parking/admin/list", middleware.RequireRole("admin", "property"), securityHandler.GetAllParking)
		private.GET("/parking/admin/stats", middleware.RequireRole("admin", "property"), securityHandler.GetParkingStats)
		private.POST("/parking/admin/assign", middleware.RequireRole("admin", "property"), securityHandler.AssignParking)
		private.POST("/parking/admin/create", middleware.RequireRole("admin", "property"), securityHandler.CreateParking)

		private.POST("/property/admin/create", middleware.RequireRole("admin", "property"), financeHandler.CreatePropertyFee)
		private.GET("/property/admin/list", middleware.RequireRole("admin", "property"), financeHandler.ListAllPropertyFees)

		private.POST("/admin/role/create", middleware.RequireRole("admin"), adminHandler.CreateRole)
		private.GET("/admin/role/list", middleware.RequireRole("admin"), adminHandler.ListRoles)
		private.POST("/admin/menu/create", middleware.RequireRole("admin"), adminHandler.CreateMenu)
		private.GET("/admin/menu/list", middleware.RequireRole("admin"), adminHandler.ListMenus)
		private.POST("/admin/role/bind_menu", middleware.RequireRole("admin"), adminHandler.BindRoleMenu)
		private.GET("/admin/user/list", middleware.RequireRole("admin"), adminHandler.ListUsers)
		private.POST("/admin/user/freeze", middleware.RequireRole("admin"), adminHandler.FreezeUser)
		private.POST("/admin/user/assign_role", middleware.RequireRole("admin"), adminHandler.AssignRole)
		private.POST("/admin/user/update_balance", middleware.RequireRole("admin"), adminHandler.UpdateUserBalance)
		private.POST("/admin/ai-report/generate", middleware.RequireRole("admin"), adminHandler.GenerateAIReport)
		private.GET("/admin/ai-report/list", middleware.RequireRole("admin"), adminHandler.ListAIReports)
		private.GET("/admin/ai-report/:id", middleware.RequireRole("admin"), adminHandler.GetAIReportDetail)
		private.GET("/admin/ai-report", middleware.RequireRole("admin"), adminHandler.GetAIReport)

		private.POST("/comment/create", commentHandler.Create)
		private.POST("/chat/send", aiHandler.Send)
		private.GET("/chat/history", aiHandler.History)
		private.POST("/community/message", communityMessageHandler.Send)
		private.GET("/community/messages", communityMessageHandler.List)
	}
}

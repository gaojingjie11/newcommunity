package router

import (
	"strconv"

	"smartcommunity-microservices/pkg/middleware"
	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/user-service/internal/handler"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
)

type RouterConfig struct {
	AuthHandler   *handler.AuthHandler
	UserHandler   *handler.UserHandler
	AdminHandler  *handler.AdminHandler
	JWTSecret     string
	RedisClient   *goredis.Client
	InternalToken string // X-Internal-Token for service-to-service calls
}

func SetupRoutes(r *gin.Engine, cfg RouterConfig) {
	// Public routes (no auth)
	public := r.Group("/api/users")
	{
		public.POST("/register", cfg.AuthHandler.Register)                 // AUTH-001
		public.POST("/login", cfg.AuthHandler.Login)                       // AUTH-002
		public.POST("/sms-code", cfg.AuthHandler.SendLoginCode)            // AUTH-008a
		public.POST("/login-code", cfg.AuthHandler.LoginByCode)            // AUTH-008b
		public.POST("/password-reset/code", cfg.AuthHandler.SendResetCode) // AUTH-003a
		public.POST("/password-reset", cfg.AuthHandler.ResetPassword)      // AUTH-003b
	}

	// Authenticated user routes
	authMw := middleware.JWTAuth(cfg.JWTSecret, cfg.RedisClient)
	user := r.Group("/api/users")
	user.Use(authMw)
	{
		user.GET("/me", cfg.UserHandler.GetProfile)              // AUTH-005, AUTH-006
		user.PUT("/me", cfg.UserHandler.UpdateProfile)           // AUTH-005
		user.PUT("/me/password", cfg.AuthHandler.ChangePassword) // AUTH-004
		user.POST("/me/face", cfg.UserHandler.RegisterFace)      // Face registration
		user.POST("/logout", cfg.AuthHandler.Logout)             // AUTH-007
	}

	// Admin routes — permission-based access control
	admin := r.Group("/api/admin")
	admin.Use(authMw)

	perm := func(code string) gin.HandlerFunc {
		return middleware.RequirePermission(cfg.RedisClient, cfg.AdminHandler, code)
	}

	// RBAC: Role management (ADMIN-MALL-001)
	admin.POST("/roles", perm("rbac:role:create"), cfg.AdminHandler.CreateRole)
	admin.PUT("/roles", perm("rbac:role:update"), cfg.AdminHandler.UpdateRole)
	admin.DELETE("/roles", perm("rbac:role:delete"), cfg.AdminHandler.DeleteRole)
	admin.GET("/roles", perm("rbac:role:list"), cfg.AdminHandler.ListRoles)

	// RBAC: Role-menu binding (ADMIN-MALL-002)
	admin.POST("/roles/:id/menus", perm("rbac:role:bind_menu"), cfg.AdminHandler.BindRoleMenu)

	// RBAC: Role-permission binding
	admin.POST("/roles/:id/permissions", perm("rbac:role:bind_permission"), cfg.AdminHandler.BindRolePermissions)
	admin.GET("/roles/:id/permissions", perm("rbac:role:get_permissions"), cfg.AdminHandler.GetRolePermissions)

	// RBAC: User management (ADMIN-MALL-003)
	admin.GET("/users", perm("rbac:user:list"), cfg.AdminHandler.ListAdminUsers)
	admin.POST("/users/freeze", perm("rbac:user:freeze"), cfg.AdminHandler.FreezeUser)
	admin.POST("/users/assign-role", perm("rbac:user:assign_role"), cfg.AdminHandler.AssignRole)
	admin.POST("/users/update-balance", perm("rbac:user:update_balance"), cfg.AdminHandler.UpdateUserBalance)
	admin.POST("/users/:id/roles", perm("rbac:user:assign_roles"), cfg.AdminHandler.AssignUserRoles)
	admin.GET("/users/:id/roles", perm("rbac:user:get_roles"), cfg.AdminHandler.GetUserRoles)

	// RBAC: Member list (ADMIN-MALL-004)
	admin.GET("/members", perm("rbac:member:list"), cfg.AdminHandler.ListMembers)

	// RBAC: Permission and menu queries
	admin.GET("/permissions", perm("rbac:permission:list"), cfg.AdminHandler.ListPermissions)
	admin.GET("/menus", perm("rbac:menu:list"), cfg.AdminHandler.ListMenus)

	// Login logs (LOG-001, LOG-002)
	admin.GET("/user-login-logs", perm("log:user_login:list"), cfg.AdminHandler.QueryUserLoginLogs)
	admin.GET("/admin-login-logs", perm("log:admin_login:list"), cfg.AdminHandler.QueryAdminLoginLogs)

	// Internal routes (service-to-service, X-Internal-Token)
	if cfg.InternalToken != "" {
		internal := r.Group("/api/internal/users")
		internal.Use(middleware.InternalToken(cfg.InternalToken))
		{
			// Permission query for other services (e.g., gateway)
			internal.GET("/permissions/:user_id", func(c *gin.Context) {
				userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
				if err != nil {
					response.Error(c, 400, 400, "参数错误", nil)
					return
				}
				codes, err := cfg.AdminHandler.GetPermissionCodesByUserID(userID)
				if err != nil {
					response.Error(c, 500, 500, "查询失败", nil)
					return
				}
				response.Success(c, gin.H{"permissions": codes})
			})
		}
	}
}

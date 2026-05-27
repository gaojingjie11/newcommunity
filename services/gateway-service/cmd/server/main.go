package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http/httputil"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"smartcommunity-microservices/pkg/config"
	"smartcommunity-microservices/pkg/db"
	"smartcommunity-microservices/pkg/logger"
	"smartcommunity-microservices/pkg/middleware"
	storage "smartcommunity-microservices/pkg/minio"
	"smartcommunity-microservices/pkg/nacos"
	"smartcommunity-microservices/pkg/redis"
	"smartcommunity-microservices/pkg/response"

	"smartcommunity-microservices/services/gateway-service/internal/discovery"
	"smartcommunity-microservices/services/gateway-service/internal/perm"

	"github.com/gin-gonic/gin"
	miniogo "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
	goredis "github.com/redis/go-redis/v9"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "path to config file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}
	if cfg.Service.Name == "" {
		cfg.Service.Name = "gateway-service"
	}
	if cfg.Service.Port == 0 {
		cfg.Service.Port = 8000
	}

	logr := logger.New(cfg.Service.Name)

	// MySQL — required for permission queries
	database, err := db.InitMySQL(cfg.MySQL)
	if err != nil {
		log.Fatalf("init mysql failed: %v", err)
	}

	// Redis — required for JWT validation and permission cache
	rdb, err := redis.Init(cfg.Redis)
	if err != nil {
		log.Fatalf("init redis failed: %v", err)
	}

	minioClient, err := storage.Init(cfg.MinIO)
	if err != nil {
		logr.Warn("minio init failed, upload will be unavailable", "error", err)
	} else {
		ensureMinIOBucket(context.Background(), minioClient, cfg.MinIO, logr)
	}

	// JWT config
	jwtSecret := "dev-secret-change-me"
	if v, ok := cfg.Raw["jwt"].(map[string]interface{}); ok {
		if s, ok := v["secret"].(string); ok && s != "" {
			jwtSecret = s
		}
	}

	// Internal token
	internalToken := cfg.Gateway.InternalToken

	// Nacos registration (best-effort)
	if err := nacos.RegisterService(context.Background(), cfg.Nacos, cfg.Service.Name, cfg.Service.Host, cfg.Service.Port, map[string]string{"kind": "gateway"}, cfg.Service.RegisterIP); err != nil {
		logr.Warn("nacos registration skipped", "error", err)
	}

	// Service discovery: Nacos with local config fallback
	resolver := discovery.NewResolver(cfg.Nacos, cfg.Gateway.Services, logr)
	resolver.StartRefresh(context.Background())

	// Permission provider (queries sys_user_role + sys_role_permission)
	permProvider := perm.NewPermissionProvider(database)

	// Middleware
	authMw := middleware.JWTAuth(jwtSecret, rdb)

	// Router
	r := gin.New()
	r.Use(middleware.RequestID(), middleware.Logger(logr), middleware.Recovery(logr), middleware.CORS())

	// ── Public routes (no JWT) ──
	r.GET("/health", func(c *gin.Context) {
		response.Success(c, gin.H{"service": cfg.Service.Name, "status": "ok"})
	})
	r.GET("/api/gateway/services", func(c *gin.Context) {
		response.Success(c, gin.H{"services": resolver.Services(), "source": resolver.Source()})
	})
	// ── Public proxied routes ──
	publicPaths := []string{
		"/api/users/register",
		"/api/users/login",
		"/api/users/sms-code",
		"/api/users/login-code",
		"/api/users/password-reset/code",
		"/api/users/password-reset",
		"/api/mall/products",
		"/api/mall/products/search",
		"/api/mall/products/promotions",
		"/api/mall/stores",
		"/api/mall/stores/:id",
		"/api/mall/promotions",
		"/api/mall/promotions/:id",
		"/api/mall/categories",
		"/api/mall/categories/:id",
		"/api/mall/service-areas",
		"/api/community/ping",
		"/api/community/notices",
		"/api/community/notices/:id",
		"/api/workorders/ping",
		"/agent/health",
	}
	for _, p := range publicPaths {
		r.Any(p, proxyHandler(resolver, internalToken, logr))
	}

	// ── Authenticated proxied routes (JWT required, no specific permission) ──
	authGroup := r.Group("")
	authGroup.Use(authMw)
	authGroup.POST("/api/upload", uploadHandler(minioClient, cfg.MinIO))
	authenticatedPaths := []string{
		// user-service
		"/api/users/me",
		"/api/users/me/password",
		"/api/users/me/face",
		"/api/users/logout",
		// mall-service
		"/api/mall/cart/items",
		"/api/mall/cart/items/:id",
		"/api/mall/favorites",
		"/api/mall/favorites/:product_id",
		"/api/mall/favorites/check/:product_id",
		"/api/mall/products/:id",
		"/api/mall/comments",
		"/api/mall/orders",
		"/api/mall/orders/:id",
		"/api/mall/orders/:id/pay",
		"/api/mall/orders/:id/cancel",
		"/api/mall/orders/:id/receive",
		"/api/mall/orders/:id/payment-status",
		"/api/mall/wallet/recharge",
		"/api/mall/wallet/transfer",
		"/api/mall/wallet/balance",
		"/api/mall/wallet/transactions",
		// community-service
		"/api/community/notices/:id/read",
		"/api/community/visitors",
		"/api/community/parking-spaces/my",
		"/api/community/parking-spaces/:id/plate",
		"/api/community/property-fees",
		"/api/community/property-fees/:id/pay",
		"/api/community/property-fees/payments",
		// community-service (workorder)
		"/api/workorders",
		"/api/workorders/:id/logs",
		// community-service (chat)
		"/api/community/messages",
		"/api/community/message",
		// agent-service
		"/agent/chat",
		"/agent/repair-classify",
		"/agent/complaint-risk",
		"/agent/recommend",
	}
	for _, p := range authenticatedPaths {
		authGroup.Any(p, proxyHandler(resolver, internalToken, logr))
	}

	// ── Permission-required proxied routes (JWT + RBAC) ──
	permGroup := r.Group("")
	permGroup.Use(authMw, gatewayRequirePermission(rdb, permProvider))
	// All /api/admin/** and /api/statistics/** routes go through permission check
	permGroup.Any("/api/admin/*path", proxyHandler(resolver, internalToken, logr))
	permGroup.Any("/api/statistics/*path", proxyHandler(resolver, internalToken, logr))

	addr := fmt.Sprintf("%s:%d", cfg.Service.Host, cfg.Service.Port)
	logr.Info("starting service", "addr", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("run server failed: %v", err)
	}
}

func ensureMinIOBucket(ctx context.Context, client *miniogo.Client, cfg config.MinIOConfig, logr *slog.Logger) {
	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		logr.Warn("check minio bucket failed", "bucket", cfg.Bucket, "error", err)
		return
	}
	if !exists {
		if err := client.MakeBucket(ctx, cfg.Bucket, miniogo.MakeBucketOptions{}); err != nil {
			logr.Warn("create minio bucket failed", "bucket", cfg.Bucket, "error", err)
			return
		}
	}

	policy := fmt.Sprintf(`{
		"Version":"2012-10-17",
		"Statement":[{
			"Effect":"Allow",
			"Principal":{"AWS":["*"]},
			"Action":["s3:GetObject"],
			"Resource":["arn:aws:s3:::%s/*"]
		}]
	}`, cfg.Bucket)
	if err := client.SetBucketPolicy(ctx, cfg.Bucket, policy); err != nil {
		logr.Warn("set minio public read policy failed", "bucket", cfg.Bucket, "error", err)
	}

	// Set Lifecycle Policy for face/temp-pay/ to expire after 7 days
	lifecycleConfig := lifecycle.NewConfiguration()
	lifecycleConfig.Rules = []lifecycle.Rule{
		{
			ID:     "expire-temp-payment-faces",
			Status: "Enabled",
			Prefix: "face/temp-pay/",
			Expiration: lifecycle.Expiration{
				Days: 7,
			},
		},
	}
	if err := client.SetBucketLifecycle(ctx, cfg.Bucket, lifecycleConfig); err != nil {
		logr.Warn("set minio bucket lifecycle policy failed", "bucket", cfg.Bucket, "error", err)
	}
}

func uploadHandler(client *miniogo.Client, cfg config.MinIOConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		if client == nil {
			response.Error(c, 503, 503, "对象存储不可用", nil)
			return
		}

		file, err := c.FormFile("file")
		if err != nil {
			response.Error(c, 400, 400, "请选择要上传的文件", nil)
			return
		}
		src, err := file.Open()
		if err != nil {
			response.Error(c, 500, 500, "文件读取失败", nil)
			return
		}
		defer src.Close()

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext == "" {
			ext = ".bin"
		}
		dir := uploadDir(c.PostForm("dir"), file.Filename, file.Header.Get("Content-Type"))
		objectName := fmt.Sprintf("%s/%d%s", dir, time.Now().UnixNano(), ext)
		contentType := file.Header.Get("Content-Type")

		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()
		info, err := client.PutObject(ctx, cfg.Bucket, objectName, src, file.Size, miniogo.PutObjectOptions{ContentType: contentType})
		if err != nil {
			response.Error(c, 500, 500, "文件上传失败", nil)
			return
		}

		response.Success(c, gin.H{
			"url": minioObjectURL(cfg, info.Key),
			"key": info.Key,
		})
	}
}

func uploadDir(explicit, filename, contentType string) string {
	explicit = strings.Trim(strings.ToLower(explicit), "/ ")
	if explicit == "face/temp-pay" {
		return explicit
	}
	switch explicit {
	case "face", "image", "common":
		return explicit
	}
	name := strings.ToLower(filename)
	if strings.Contains(name, "face") {
		return "face"
	}
	if strings.HasPrefix(contentType, "image/") {
		return "image"
	}
	return "common"
}

func minioObjectURL(cfg config.MinIOConfig, key string) string {
	protocol := "http"
	if cfg.UseSSL {
		protocol = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", protocol, cfg.Endpoint, cfg.Bucket, key)
}

// proxyHandler returns a gin.HandlerFunc that proxies the request to the appropriate service.
func proxyHandler(resolver *discovery.Resolver, internalToken string, logr *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		serviceName := resolveServiceName(c.Request.URL.Path)
		if serviceName == "" {
			response.Error(c, 404, 404, "no upstream service for path", gin.H{"path": c.Request.URL.Path})
			return
		}

		base := resolver.Resolve(serviceName)
		if base == "" {
			response.Error(c, 502, 502, "service unavailable", gin.H{"service": serviceName})
			return
		}

		target, err := url.Parse(base)
		if err != nil {
			response.Error(c, 500, 500, "invalid upstream url", gin.H{"service": serviceName})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(target)
		c.Request.Host = target.Host

		// Inject headers
		c.Request.Header.Set("X-Gateway-Proxy", "gateway-service")
		c.Request.Header.Set("X-Gateway-Time", time.Now().Format(time.RFC3339))
		// Always forward the resolved request ID (from RequestID middleware context)
		if rid := c.GetString("request_id"); rid != "" {
			c.Request.Header.Set("X-Request-ID", rid)
		}
		if internalToken != "" {
			c.Request.Header.Set("X-Internal-Token", internalToken)
		}
		// Forward user identity from JWT context
		if uid, exists := c.Get("userID"); exists {
			c.Request.Header.Set("X-User-ID", fmt.Sprintf("%d", uid.(int64)))
		}
		if role, exists := c.Get("role"); exists {
			c.Request.Header.Set("X-User-Role", role.(string))
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

// resolveServiceName maps a URL path to a service name.
func resolveServiceName(path string) string {
	switch {
	case strings.HasPrefix(path, "/api/users"),
		strings.HasPrefix(path, "/api/admin/roles"),
		strings.HasPrefix(path, "/api/admin/users"),
		strings.HasPrefix(path, "/api/admin/members"),
		strings.HasPrefix(path, "/api/admin/permissions"),
		strings.HasPrefix(path, "/api/admin/menus"),
		strings.HasPrefix(path, "/api/admin/user-login-logs"),
		strings.HasPrefix(path, "/api/admin/admin-login-logs"):
		return "user-service"
	case strings.HasPrefix(path, "/api/mall"),
		strings.HasPrefix(path, "/api/admin/mall"):
		return "mall-service"
	case strings.HasPrefix(path, "/api/community"),
		strings.HasPrefix(path, "/api/admin/community"),
		strings.HasPrefix(path, "/api/workorders"),
		strings.HasPrefix(path, "/api/admin/workorders"),
		strings.HasPrefix(path, "/api/statistics"):
		return "community-service"
	case strings.HasPrefix(path, "/agent"):
		return "agent-service"
	default:
		return ""
	}
}

// gatewayRequirePermission returns middleware that checks RBAC permissions for admin/statistics routes.
func gatewayRequirePermission(rdb *goredis.Client, provider middleware.PermissionProvider) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			response.Error(c, 401, 401, "请先登录", nil)
			c.Abort()
			return
		}

		uid, ok := userID.(int64)
		if !ok {
			response.Error(c, 401, 401, "用户ID无效", nil)
			c.Abort()
			return
		}

		// Admin role shortcut
		if role, _ := c.Get("role"); role == "admin" {
			c.Next()
			return
		}

		// Lookup required permission for this method+path
		requiredPerm := perm.LookupPermission(c.Request.Method, c.Request.URL.Path)
		if requiredPerm == "" {
			// Fail closed: admin/statistics paths without a mapping are denied.
			// Add the missing route to internal/perm/mapping.go.
			response.Error(c, 403, 403, "该管理接口未配置权限映射", gin.H{
				"method": c.Request.Method,
				"path":   c.Request.URL.Path,
			})
			c.Abort()
			return
		}

		// Load permissions from cache or DB (reuse pkg/middleware logic)
		codes, err := loadPermissions(rdb, provider, uid)
		if err != nil {
			response.Error(c, 500, 500, "权限校验失败", nil)
			c.Abort()
			return
		}

		for _, code := range codes {
			if code == requiredPerm {
				c.Next()
				return
			}
		}

		response.Error(c, 403, 403, "无权限访问此资源", nil)
		c.Abort()
	}
}

func loadPermissions(rdb *goredis.Client, provider middleware.PermissionProvider, userID int64) ([]string, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("rbac:permissions:%d", userID)

	// Try cache first
	cached, err := rdb.SMembers(ctx, cacheKey).Result()
	if err == nil && len(cached) > 0 {
		return cached, nil
	}

	// Cache miss — load from DB
	codes, err := provider.GetPermissionCodesByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Write to cache with 10-min TTL
	if len(codes) > 0 {
		members := make([]interface{}, len(codes))
		for i, c := range codes {
			members[i] = c
		}
		_ = rdb.SAdd(ctx, cacheKey, members...).Err()
		_ = rdb.Expire(ctx, cacheKey, 600).Err()
	}

	return codes, nil
}

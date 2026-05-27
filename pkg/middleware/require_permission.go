package middleware

import (
	"context"
	"fmt"

	"smartcommunity-microservices/pkg/response"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
)

// PermissionProvider loads permission codes for a user.
// Implementations typically query sys_user_role + sys_role_permission from DB.
type PermissionProvider interface {
	GetPermissionCodesByUserID(userID int64) ([]string, error)
}

// RequirePermission returns middleware that checks if the current user
// has a specific permission code. It uses Redis cache with fallback to DB.
func RequirePermission(rdb *goredis.Client, provider PermissionProvider, permissionCode string) gin.HandlerFunc {
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

		// Admin role shortcut: if JWT role is "admin", allow all
		if role, _ := c.Get("role"); role == "admin" {
			c.Next()
			return
		}

		codes, err := loadPermissions(rdb, provider, uid)
		if err != nil {
			response.Error(c, 500, 500, "权限校验失败", nil)
			c.Abort()
			return
		}

		for _, code := range codes {
			if code == permissionCode {
				c.Next()
				return
			}
		}

		response.Error(c, 403, 403, "无权限访问此资源", nil)
		c.Abort()
	}
}

func loadPermissions(rdb *goredis.Client, provider PermissionProvider, userID int64) ([]string, error) {
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

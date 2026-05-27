package middleware

import (
	"context"
	"fmt"
	"strings"

	"smartcommunity-microservices/pkg/auth"
	"smartcommunity-microservices/pkg/response"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
)

func JWTAuth(secret string, rdb *goredis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, 401, 401, "请先登录", nil)
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, 401, 401, "Token格式错误", nil)
			c.Abort()
			return
		}

		claims, err := auth.ParseToken(secret, parts[1])
		if err != nil {
			response.Error(c, 401, 401, "Token无效或已过期", nil)
			c.Abort()
			return
		}

		redisKey := fmt.Sprintf("login:token:%d", claims.UserID)
		cachedToken, err := rdb.Get(context.Background(), redisKey).Result()
		if err != nil || cachedToken != parts[1] {
			response.Error(c, 401, 401, "登录已失效，请重新登录", nil)
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// OptionalJWTAuth enriches public routes with user identity when a valid token is
// present, while still allowing anonymous access when no token is provided.
func OptionalJWTAuth(secret string, rdb *goredis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || rdb == nil {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		claims, err := auth.ParseToken(secret, parts[1])
		if err != nil {
			c.Next()
			return
		}

		redisKey := fmt.Sprintf("login:token:%d", claims.UserID)
		cachedToken, err := rdb.Get(context.Background(), redisKey).Result()
		if err != nil || cachedToken != parts[1] {
			c.Next()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

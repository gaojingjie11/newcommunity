package middleware

import (
	"smartcommunity-microservices/pkg/response"

	"github.com/gin-gonic/gin"
)

// InternalToken creates a middleware that validates X-Internal-Token for service-to-service calls.
// If the header matches the configured token, the request proceeds.
// If empty or mismatched, returns 403 Forbidden.
//
// Usage:
//
//	internal := middleware.InternalToken("your-secret-token")
//	protectedGroup.POST("/internal/action", internal, handler.Action)
func InternalToken(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if token == "" {
			// No token configured — reject all internal calls
			response.Error(c, 403, 403, "服务间调用未配置", nil)
			c.Abort()
			return
		}

		header := c.GetHeader("X-Internal-Token")
		if header != token {
			response.Error(c, 403, 403, "服务间鉴权失败", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}

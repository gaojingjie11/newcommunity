package middleware

import (
	"smartcommunity-microservices/pkg/response"

	"github.com/gin-gonic/gin"
)

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			response.Error(c, 403, 403, "无权限: 未获取到角色信息", nil)
			c.Abort()
			return
		}
		roleStr, ok := userRole.(string)
		if !ok {
			response.Error(c, 403, 403, "无权限: 角色类型错误", nil)
			c.Abort()
			return
		}
		for _, role := range roles {
			if role == roleStr {
				c.Next()
				return
			}
		}
		response.Error(c, 403, 403, "无权限访问此资源", nil)
		c.Abort()
	}
}

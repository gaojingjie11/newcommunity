package handler

import (
	"net/http"
	"strconv"

	"smartcommunity-microservices/pkg/response"

	"github.com/gin-gonic/gin"
)

func currentUserID(c *gin.Context) (int64, bool) {
	v, ok := c.Get("userID")
	if !ok {
		response.Error(c, http.StatusUnauthorized, 401, "请先登录", nil)
		return 0, false
	}
	id, ok := v.(int64)
	if !ok || id <= 0 {
		response.Error(c, http.StatusUnauthorized, 401, "用户ID无效", nil)
		return 0, false
	}
	return id, true
}

func parseIDParam(c *gin.Context, name string) (int64, bool) {
	id, err := strconv.ParseInt(c.Param(name), 10, 64)
	if err != nil || id <= 0 {
		response.Error(c, http.StatusBadRequest, 400, "ID参数无效", nil)
		return 0, false
	}
	return id, true
}

func parseStatusQuery(c *gin.Context) (*int, bool) {
	raw := c.Query("status")
	if raw == "" {
		return nil, true
	}
	value, err := strconv.Atoi(raw)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "status参数无效", nil)
		return nil, false
	}
	return &value, true
}

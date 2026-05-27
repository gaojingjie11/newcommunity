package response

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ParsePage extracts page and size from query params with defaults (page=1, size=20, max=100).
func ParsePage(c *gin.Context) (page, size int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ = strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return
}

type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Body{Code: 0, Message: "success", Data: data})
}

func Error(c *gin.Context, status int, code int, message string, data interface{}) {
	c.JSON(status, Body{Code: code, Message: message, Data: data})
}

// PagedResult is the standard wrapper for paginated list responses.
// JSON shape: {"list": [...], "total": N, "page": P, "size": S}
type PagedResult struct {
	List  interface{} `json:"list"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// SuccessPaged sends a paginated success response.
func SuccessPaged(c *gin.Context, list interface{}, total int64, page, size int) {
	Success(c, PagedResult{List: list, Total: total, Page: page, Size: size})
}

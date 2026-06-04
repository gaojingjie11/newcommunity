package response

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
)

type Body struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func ParsePage(r *http.Request) (page, size int) {
	q := r.URL.Query()
	pageVal := q.Get("page")
	sizeVal := q.Get("size")

	page, _ = strconv.Atoi(pageVal)
	size, _ = strconv.Atoi(sizeVal)

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

func Response(w http.ResponseWriter, resp interface{}, err error) {
	if err != nil {
		code := 400
		msg := err.Error()

		// 1. 尝试从 gRPC 错误提取纯净业务描述
		if st, ok := status.FromError(err); ok {
			msg = st.Message()
		}

		// 2. 将难看或不友好的底层解析、语法和超时报错转换成友好提示
		if strings.Contains(msg, "strconv.Parse") || strings.Contains(msg, "invalid syntax") {
			msg = "请求参数格式不正确"
		} else if strings.Contains(msg, "sql:") || strings.Contains(msg, "database") || strings.Contains(msg, "GORM") {
			msg = "系统繁忙，数据查询异常"
		} else if strings.Contains(msg, "context deadline exceeded") || strings.Contains(msg, "timeout") {
			msg = "网络请求超时，请稍后重试"
		}

		type coder interface {
			ErrorCode() int
		}
		if c, ok := err.(coder); ok {
			code = c.ErrorCode()
		}

		httpx.WriteJson(w, http.StatusOK, Body{
			Code:    code,
			Message: msg,
			Data:    nil,
		})
		return
	}
	httpx.WriteJson(w, http.StatusOK, Body{
		Code:    0,
		Message: "success",
		Data:    resp,
	})
}

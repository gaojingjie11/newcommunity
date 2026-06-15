package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"smartcommunity-microservices/app/gateway/api/internal/logic"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/common/response"
)

func CreateConversationHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateConversationReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCreateConversationLogic(r.Context(), svcCtx)
		resp, err := l.CreateConversation(&req)
		response.Response(w, resp, err)
	}
}

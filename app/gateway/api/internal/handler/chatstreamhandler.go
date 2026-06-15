// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"smartcommunity-microservices/app/gateway/api/internal/logic"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
)

func ChatStreamHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ChatStreamReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewChatStreamLogic(r.Context(), svcCtx)
		err := l.ChatStream(w, &req)
		if err != nil {
			// If header is already sent, we can't write a JSON error anymore.
			// Just write an error block into the stream or log it.
			if w.Header().Get("Content-Type") != "text/event-stream" {
				httpx.ErrorCtx(r.Context(), w, err)
			}
		}
	}
}

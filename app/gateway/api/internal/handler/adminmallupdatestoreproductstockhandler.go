// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"smartcommunity-microservices/common/response"
	"smartcommunity-microservices/app/gateway/api/internal/logic"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
)

func AdminMallUpdateStoreProductStockHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminUpdateStoreProductStockReq
		if err := httpx.Parse(r, &req); err != nil {
			response.Response(w, nil, err)
			return
		}

		l := logic.NewAdminMallUpdateStoreProductStockLogic(r.Context(), svcCtx)
		resp, err := l.AdminMallUpdateStoreProductStock(&req)
		response.Response(w, resp, err)
	}
}

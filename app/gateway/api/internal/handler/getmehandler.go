package handler

import (
	"net/http"

	"smartcommunity-microservices/app/gateway/api/internal/logic"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/common/response"
)

func GetMeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetMeLogic(r.Context(), svcCtx)
		resp, err := l.GetMe()
		response.Response(w, resp, err)
	}
}

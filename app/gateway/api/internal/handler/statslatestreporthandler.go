// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"smartcommunity-microservices/common/response"
	"smartcommunity-microservices/app/gateway/api/internal/logic"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
)

func StatsLatestReportHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewStatsLatestReportLogic(r.Context(), svcCtx)
		resp, err := l.StatsLatestReport()
		response.Response(w, resp, err)
	}
}

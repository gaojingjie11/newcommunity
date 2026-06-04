package logic

import (
	"context"

	"smartcommunity-microservices/app/stats/rpc/internal/svc"
	"smartcommunity-microservices/app/stats/rpc/types/stats"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderStatsCombinedLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderStatsCombinedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderStatsCombinedLogic {
	return &GetOrderStatsCombinedLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOrderStatsCombinedLogic) GetOrderStatsCombined(in *stats.DaysReq) (*stats.OrderStatsResp, error) {
	summary, trend, err := l.svcCtx.StatsSvc.OrderStatsCombined(int(in.Days))
	if err != nil {
		return nil, err
	}

	var pbSummary []*stats.OrderSummaryInfo
	for _, s := range summary {
		pbSummary = append(pbSummary, &stats.OrderSummaryInfo{
			Status:      int32(s.Status),
			Count:       s.Count,
			TotalAmount: float64(s.TotalAmount) / 100.0,
		})
	}

	var pbTrend []*stats.OrderTrendInfo
	for _, t := range trend {
		pbTrend = append(pbTrend, &stats.OrderTrendInfo{
			Date:   t.Date,
			Count:  t.Count,
			Amount: float64(t.Amount) / 100.0,
		})
	}

	return &stats.OrderStatsResp{
		Summaries: pbSummary,
		Trends:    pbTrend,
	}, nil
}

package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/stats/rpc/statsrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatsOrdersCombinedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStatsOrdersCombinedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatsOrdersCombinedLogic {
	return &StatsOrdersCombinedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StatsOrdersCombinedLogic) StatsOrdersCombined(req *types.DaysReq) (resp *types.OrderStatsResp, err error) {
	rpcResp, err := l.svcCtx.StatsRpc.GetOrderStatsCombined(l.ctx, &statsrpc.DaysReq{
		Days: req.Days,
	})
	if err != nil {
		return nil, err
	}
	summaries := make([]types.OrderSummaryInfo, 0, len(rpcResp.Summaries))
	for _, item := range rpcResp.Summaries {
		summaries = append(summaries, toAPIOrderSummaryInfo(item))
	}
	trends := make([]types.OrderTrendInfo, 0, len(rpcResp.Trends))
	for _, item := range rpcResp.Trends {
		trends = append(trends, toAPIOrderTrendInfo(item))
	}
	return &types.OrderStatsResp{
		Summaries: summaries,
		Trends:    trends,
	}, nil
}

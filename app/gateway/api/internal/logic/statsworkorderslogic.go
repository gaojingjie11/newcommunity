package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/stats/rpc/statsrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatsWorkordersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStatsWorkordersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatsWorkordersLogic {
	return &StatsWorkordersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StatsWorkordersLogic) StatsWorkorders() (resp *types.WorkorderStatsResp, err error) {
	rpcResp, err := l.svcCtx.StatsRpc.GetWorkorderStats(l.ctx, &statsrpc.BaseResp{})
	if err != nil {
		return nil, err
	}
	summaries := make([]types.WorkorderSummaryInfo, 0, len(rpcResp.Summaries))
	for _, item := range rpcResp.Summaries {
		summaries = append(summaries, toAPIWorkorderSummaryInfo(item))
	}
	return &types.WorkorderStatsResp{
		Summaries: summaries,
	}, nil
}

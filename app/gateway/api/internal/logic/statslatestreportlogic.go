package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/stats/rpc/statsrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatsLatestReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStatsLatestReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatsLatestReportLogic {
	return &StatsLatestReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StatsLatestReportLogic) StatsLatestReport() (resp *types.ReportResp, err error) {
	rpcResp, err := l.svcCtx.StatsRpc.GetLatestAIReport(l.ctx, &statsrpc.BaseResp{})
	if err != nil {
		return nil, err
	}
	return &types.ReportResp{
		Report: toAPIAIReportInfo(rpcResp.Report),
	}, nil
}

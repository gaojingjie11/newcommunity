package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/stats/rpc/statsrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatsGenerateReportLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStatsGenerateReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatsGenerateReportLogic {
	return &StatsGenerateReportLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StatsGenerateReportLogic) StatsGenerateReport() (resp *types.ReportResp, err error) {
	rpcResp, err := l.svcCtx.StatsRpc.GenerateAIReport(l.ctx, &statsrpc.GenerateReportReq{
		UserId: getUserIDFromCtx(l.ctx),
	})
	if err != nil {
		return nil, err
	}
	return &types.ReportResp{
		Report: toAPIAIReportInfo(rpcResp.Report),
	}, nil
}

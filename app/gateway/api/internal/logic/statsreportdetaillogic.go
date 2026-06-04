package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/stats/rpc/statsrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatsReportDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStatsReportDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatsReportDetailLogic {
	return &StatsReportDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StatsReportDetailLogic) StatsReportDetail(req *types.ReportIDReq) (resp *types.ReportResp, err error) {
	rpcResp, err := l.svcCtx.StatsRpc.GetAIReportDetail(l.ctx, &statsrpc.ReportIDReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &types.ReportResp{
		Report: toAPIAIReportInfo(rpcResp.Report),
	}, nil
}

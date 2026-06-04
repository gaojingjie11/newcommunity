package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/stats/rpc/statsrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatsListReportsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStatsListReportsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatsListReportsLogic {
	return &StatsListReportsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StatsListReportsLogic) StatsListReports(req *types.ListReportsReq) (resp *types.ReportListResp, err error) {
	rpcResp, err := l.svcCtx.StatsRpc.ListAIReports(l.ctx, &statsrpc.ListReportsReq{
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.AIReportInfo, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		list = append(list, toAPIAIReportInfo(item))
	}
	return &types.ReportListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}

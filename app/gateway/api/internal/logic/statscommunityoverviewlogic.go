package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/stats/rpc/statsrpc"

	"github.com/zeromicro/go-zero/core/logx"
)

type StatsCommunityOverviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStatsCommunityOverviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StatsCommunityOverviewLogic {
	return &StatsCommunityOverviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *StatsCommunityOverviewLogic) StatsCommunityOverview() (resp *types.CommunityOverviewResp, err error) {
	rpcResp, err := l.svcCtx.StatsRpc.GetCommunityOverview(l.ctx, &statsrpc.BaseResp{})
	if err != nil {
		return nil, err
	}
	return &types.CommunityOverviewResp{
		UserCount:      rpcResp.UserCount,
		OrderCount:     rpcResp.OrderCount,
		PaidAmount:     rpcResp.PaidAmount,
		RepairCount:    rpcResp.RepairCount,
		ComplaintCount: rpcResp.ComplaintCount,
		FeeCount:       rpcResp.FeeCount,
		FeePaidCount:   rpcResp.FeePaidCount,
	}, nil
}

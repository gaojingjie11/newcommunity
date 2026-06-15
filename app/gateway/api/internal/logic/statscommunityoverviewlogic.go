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

	var repairStats []types.RepairStatInfo
	for _, item := range rpcResp.RepairStats {
		repairStats = append(repairStats, types.RepairStatInfo{
			Name:  item.Name,
			Value: item.Value,
		})
	}

	return &types.CommunityOverviewResp{
		UserCount:      rpcResp.UserCount,
		OrderCount:     rpcResp.OrderCount,
		PaidAmount:     rpcResp.PaidAmount,
		RepairCount:    rpcResp.RepairCount,
		ComplaintCount: rpcResp.ComplaintCount,
		FeeCount:       rpcResp.FeeCount,
		FeePaidCount:   rpcResp.FeePaidCount,

		TotalUsers:     rpcResp.TotalUsers,
		TodayOrders:    rpcResp.TodayOrders,
		ParkingRate:    rpcResp.ParkingRate,
		MonthIncome:    rpcResp.MonthIncome,
		RepairStats:    repairStats,
		IncomeDates:    rpcResp.IncomeDates,
		IncomeTrend:    rpcResp.IncomeTrend,
		CostStructure:  rpcResp.CostStructure,
	}, nil
}

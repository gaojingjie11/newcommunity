package logic

import (
	"context"

	"smartcommunity-microservices/app/stats/rpc/internal/svc"
	"smartcommunity-microservices/app/stats/rpc/types/stats"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCommunityOverviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCommunityOverviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCommunityOverviewLogic {
	return &GetCommunityOverviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCommunityOverviewLogic) GetCommunityOverview(in *stats.BaseResp) (*stats.CommunityOverviewResp, error) {
	ov, err := l.svcCtx.StatsSvc.CommunityOverview()
	if err != nil {
		return nil, err
	}

	var repairStats []*stats.RepairStatInfo
	for _, r := range ov.RepairStats {
		repairStats = append(repairStats, &stats.RepairStatInfo{
			Name:  r.Name,
			Value: r.Value,
		})
	}

	return &stats.CommunityOverviewResp{
		UserCount:      ov.UserCount,
		OrderCount:     ov.OrderCount,
		PaidAmount:     float64(ov.PaidAmount) / 100.0,
		RepairCount:    ov.RepairCount,
		ComplaintCount: ov.ComplaintCount,
		FeeCount:       ov.FeeCount,
		FeePaidCount:   ov.FeePaidCount,

		TotalUsers:     ov.TotalUsers,
		TodayOrders:    ov.TodayOrders,
		ParkingRate:    ov.ParkingRate,
		MonthIncome:    ov.MonthIncome,
		RepairStats:    repairStats,
		IncomeDates:    ov.IncomeDates,
		IncomeTrend:    ov.IncomeTrend,
		CostStructure:  ov.CostStructure,
	}, nil
}

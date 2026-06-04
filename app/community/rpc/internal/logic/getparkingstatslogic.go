package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetParkingStatsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetParkingStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetParkingStatsLogic {
	return &GetParkingStatsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetParkingStatsLogic) GetParkingStats(in *community.BaseResp) (*community.ParkingStatsResp, error) {
	stats, err := l.svcCtx.ParkingRepo.Stats()
	if err != nil {
		return nil, err
	}

	return &community.ParkingStatsResp{
		Total: stats["total"],
		Bound: stats["bound"],
		Free:  stats["free"],
	}, nil
}

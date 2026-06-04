package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCommunityParkingStatsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCommunityParkingStatsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCommunityParkingStatsLogic {
	return &AdminCommunityParkingStatsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminCommunityParkingStatsLogic) AdminCommunityParkingStats() (resp *types.ParkingStatsResp, err error) {
	rpcResp, err := l.svcCtx.CommunityRpc.GetParkingStats(l.ctx, &communityrpc.BaseResp{})
	if err != nil {
		return nil, err
	}
	return &types.ParkingStatsResp{
		Total: rpcResp.Total,
		Bound: rpcResp.Bound,
		Free:  rpcResp.Free,
	}, nil
}

package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommunityBindPlateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommunityBindPlateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommunityBindPlateLogic {
	return &CommunityBindPlateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommunityBindPlateLogic) CommunityBindPlate(req *types.BindPlateReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.CommunityRpc.BindPlate(l.ctx, &communityrpc.BindPlateReq{
		UserId:         getUserIDFromCtx(l.ctx),
		ParkingSpaceId: req.Id,
		CarPlate:       req.CarPlate,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}

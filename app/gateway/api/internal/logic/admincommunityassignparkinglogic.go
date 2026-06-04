package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCommunityAssignParkingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCommunityAssignParkingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCommunityAssignParkingLogic {
	return &AdminCommunityAssignParkingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminCommunityAssignParkingLogic) AdminCommunityAssignParking(req *types.AssignParkingReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.CommunityRpc.AssignParking(l.ctx, &communityrpc.AssignParkingReq{
		Id:       req.Id,
		Mobile:   req.Mobile,
		CarPlate: req.CarPlate,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}

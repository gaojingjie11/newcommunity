package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindPlateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindPlateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindPlateLogic {
	return &BindPlateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BindPlateLogic) BindPlate(in *community.BindPlateReq) (*community.BaseResp, error) {
	_, err := l.svcCtx.ParkingRepo.BindPlate(in.ParkingSpaceId, in.UserId, in.CarPlate)
	if err != nil {
		return nil, err
	}

	return &community.BaseResp{Code: 0, Message: "success"}, nil
}

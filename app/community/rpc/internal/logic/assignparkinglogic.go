package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	"github.com/zeromicro/go-zero/core/logx"
)

type AssignParkingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAssignParkingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AssignParkingLogic {
	return &AssignParkingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AssignParkingLogic) AssignParking(in *community.AssignParkingReq) (*community.BaseResp, error) {
	_, err := l.svcCtx.ParkingRepo.Assign(in.Id, in.Mobile, in.CarPlate)
	if err != nil {
		return nil, err
	}

	return &community.BaseResp{Code: 0, Message: "success"}, nil
}

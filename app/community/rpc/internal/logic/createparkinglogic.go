package logic

import (
	"context"
	"errors"
	"strings"

	"smartcommunity-microservices/app/community/rpc/internal/model"
	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateParkingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateParkingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateParkingLogic {
	return &CreateParkingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateParkingLogic) CreateParking(in *community.CreateParkingReq) (*community.BaseResp, error) {
	parkingNo := strings.TrimSpace(in.ParkingNo)
	if parkingNo == "" {
		return nil, errors.New("parking_no required")
	}

	space := &model.ParkingSpace{
		ParkingNo: parkingNo,
		Status:    0,
	}

	if err := l.svcCtx.ParkingRepo.Create(space); err != nil {
		return nil, err
	}

	return &community.BaseResp{Code: 0, Message: "success"}, nil
}

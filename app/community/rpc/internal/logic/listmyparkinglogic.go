package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMyParkingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMyParkingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyParkingLogic {
	return &ListMyParkingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListMyParkingLogic) ListMyParking(in *community.UserIDReq) (*community.ParkingListResp, error) {
	bindings, err := l.svcCtx.ParkingRepo.ListBindingsByUser(in.UserId)
	if err != nil {
		return nil, err
	}

	var list []*community.ParkingSpaceInfo
	for _, b := range bindings {
		list = append(list, &community.ParkingSpaceInfo{
			Id:        b.ParkingSpaceID,
			ParkingNo: b.ParkingSpace.ParkingNo,
			Status:    1,
			UserId:    b.UserID,
			CarPlate:  b.CarPlate,
			BindingId: b.ID,
		})
	}

	return &community.ParkingListResp{
		List:  list,
		Total: int64(len(list)),
	}, nil
}

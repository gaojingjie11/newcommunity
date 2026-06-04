package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListParkingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminListParkingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListParkingLogic {
	return &AdminListParkingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminListParkingLogic) AdminListParking(in *community.ListParkingReq) (*community.ParkingListResp, error) {
	spaces, total, err := l.svcCtx.ParkingRepo.List(int(in.Page), int(in.Size))
	if err != nil {
		return nil, err
	}

	var list []*community.ParkingSpaceInfo
	for _, s := range spaces {
		list = append(list, &community.ParkingSpaceInfo{
			Id:         s.ID,
			ParkingNo:  s.ParkingNo,
			Status:     int32(s.Status),
			UserId:     s.UserID,
			UserName:   s.UserName,
			UserMobile: s.UserMobile,
			CarPlate:   s.CarPlate,
			BindingId:  s.BindingID,
		})
	}

	return &community.ParkingListResp{
		List:  list,
		Total: total,
	}, nil
}

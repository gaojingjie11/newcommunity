package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommunityMyParkingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommunityMyParkingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommunityMyParkingLogic {
	return &CommunityMyParkingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommunityMyParkingLogic) CommunityMyParking() (resp *types.ParkingListResp, err error) {
	rpcResp, err := l.svcCtx.CommunityRpc.ListMyParking(l.ctx, &communityrpc.UserIDReq{
		UserId: getUserIDFromCtx(l.ctx),
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.ParkingSpaceInfo, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		list = append(list, toAPIParkingSpaceInfo(item))
	}
	return &types.ParkingListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}

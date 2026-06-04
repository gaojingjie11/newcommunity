package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCommunityParkingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCommunityParkingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCommunityParkingLogic {
	return &AdminCommunityParkingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminCommunityParkingLogic) AdminCommunityParking(req *types.ListParkingReq) (resp *types.ParkingListResp, err error) {
	rpcResp, err := l.svcCtx.CommunityRpc.AdminListParking(l.ctx, &communityrpc.ListParkingReq{
		Page: req.Page,
		Size: req.Size,
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

package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCommunityFeesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCommunityFeesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCommunityFeesLogic {
	return &AdminCommunityFeesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminCommunityFeesLogic) AdminCommunityFees(req *types.ListPropertyFeesReq) (resp *types.PropertyFeeListResp, err error) {
	rpcResp, err := l.svcCtx.CommunityRpc.AdminListFees(l.ctx, &communityrpc.ListPropertyFeesReq{
		UserId: 0,
		Page:   req.Page,
		Size:   req.Size,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.PropertyFeeInfo, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		list = append(list, toAPIPropertyFeeInfo(item))
	}
	return &types.PropertyFeeListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}

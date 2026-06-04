package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCommunityVisitorsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCommunityVisitorsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCommunityVisitorsLogic {
	return &AdminCommunityVisitorsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminCommunityVisitorsLogic) AdminCommunityVisitors(req *types.ListVisitorReq) (resp *types.VisitorListResp, err error) {
	rpcResp, err := l.svcCtx.CommunityRpc.AdminListVisitors(l.ctx, &communityrpc.ListVisitorReq{
		UserId: 0,
		Page:   req.Page,
		Size:   req.Size,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.VisitorInfo, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		list = append(list, toAPIVisitorInfo(item))
	}
	return &types.VisitorListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}

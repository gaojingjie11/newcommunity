package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommunityMyVisitorsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommunityMyVisitorsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommunityMyVisitorsLogic {
	return &CommunityMyVisitorsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommunityMyVisitorsLogic) CommunityMyVisitors(req *types.ListVisitorReq) (resp *types.VisitorListResp, err error) {
	rpcResp, err := l.svcCtx.CommunityRpc.ListMyVisitors(l.ctx, &communityrpc.ListVisitorReq{
		UserId: getUserIDFromCtx(l.ctx),
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

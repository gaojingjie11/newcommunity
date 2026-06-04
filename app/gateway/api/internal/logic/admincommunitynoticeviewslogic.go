package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCommunityNoticeViewsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCommunityNoticeViewsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCommunityNoticeViewsLogic {
	return &AdminCommunityNoticeViewsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminCommunityNoticeViewsLogic) AdminCommunityNoticeViews(req *types.NoticeIDReq) (resp *types.NoticeViewsResp, err error) {
	rpcResp, err := l.svcCtx.CommunityRpc.GetNoticeViews(l.ctx, &communityrpc.NoticeIDReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &types.NoticeViewsResp{
		Views: rpcResp.Views,
	}, nil
}

package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCommunityCreateNoticeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCommunityCreateNoticeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCommunityCreateNoticeLogic {
	return &AdminCommunityCreateNoticeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminCommunityCreateNoticeLogic) AdminCommunityCreateNotice(req *types.CreateNoticeReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.CommunityRpc.CreateNotice(l.ctx, &communityrpc.CreateNoticeReq{
		Title:     req.Title,
		Content:   req.Content,
		Publisher: req.Publisher,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}

package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCommunityDeleteNoticeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCommunityDeleteNoticeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCommunityDeleteNoticeLogic {
	return &AdminCommunityDeleteNoticeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminCommunityDeleteNoticeLogic) AdminCommunityDeleteNotice(req *types.NoticeIDReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.CommunityRpc.DeleteNotice(l.ctx, &communityrpc.NoticeIDReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}

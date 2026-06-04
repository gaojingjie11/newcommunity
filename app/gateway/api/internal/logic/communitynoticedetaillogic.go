package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommunityNoticeDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommunityNoticeDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommunityNoticeDetailLogic {
	return &CommunityNoticeDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommunityNoticeDetailLogic) CommunityNoticeDetail(req *types.NoticeIDReq) (resp *types.NoticeInfo, err error) {
	rpcResp, err := l.svcCtx.CommunityRpc.GetNoticeDetail(l.ctx, &communityrpc.NoticeIDReq{
		Id:     req.Id,
		UserId: getUserIDFromCtx(l.ctx),
	})
	if err != nil {
		return nil, err
	}
	info := toAPINoticeInfo(rpcResp.Notice)
	return &info, nil
}

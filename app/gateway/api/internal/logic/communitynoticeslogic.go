package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommunityNoticesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommunityNoticesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommunityNoticesLogic {
	return &CommunityNoticesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommunityNoticesLogic) CommunityNotices(req *types.ListNoticesReq) (resp *types.NoticeListResp, err error) {
	rpcResp, err := l.svcCtx.CommunityRpc.ListNotices(l.ctx, &communityrpc.ListNoticesReq{
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.NoticeInfo, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		list = append(list, toAPINoticeInfo(item))
	}
	return &types.NoticeListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}

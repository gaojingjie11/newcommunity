package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCommunityNoticesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminCommunityNoticesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCommunityNoticesLogic {
	return &AdminCommunityNoticesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminCommunityNoticesLogic) AdminCommunityNotices(req *types.ListNoticesReq) (resp *types.NoticeListResp, err error) {
	rpcResp, err := l.svcCtx.CommunityRpc.AdminListNotices(l.ctx, &communityrpc.ListNoticesReq{
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

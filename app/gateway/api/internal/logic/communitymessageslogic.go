package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/communityrpc"
	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CommunityMessagesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCommunityMessagesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CommunityMessagesLogic {
	return &CommunityMessagesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CommunityMessagesLogic) CommunityMessages(req *types.ListMessagesReq) (resp *types.MessageListResp, err error) {
	rpcResp, err := l.svcCtx.CommunityRpc.ListMessages(l.ctx, &communityrpc.ListMessagesReq{
		Page: req.Page,
		Size: req.Size,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.MessageInfo, 0, len(rpcResp.List))
	for _, item := range rpcResp.List {
		list = append(list, toAPIMessageInfo(item))
	}
	return &types.MessageListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}

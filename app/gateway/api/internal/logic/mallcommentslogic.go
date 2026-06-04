package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallCommentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallCommentsLogic {
	return &MallCommentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallCommentsLogic) MallComments(req *types.ListCommentsReq) (resp *types.CommentListResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.ListComments(l.ctx, &mall.ListCommentsReq{
		ProductId: req.ProductId,
		Page:      req.Page,
		Size:      req.Size,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.CommentInfo, 0, len(rpcResp.List))
	for _, c := range rpcResp.List {
		list = append(list, toAPICommentInfo(c))
	}
	return &types.CommentListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}

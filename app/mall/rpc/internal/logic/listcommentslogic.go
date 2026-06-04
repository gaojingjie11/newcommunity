package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCommentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCommentsLogic {
	return &ListCommentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListCommentsLogic) ListComments(in *mall.ListCommentsReq) (*mall.CommentListResp, error) {
	comments, total, err := l.svcCtx.CommentSvc.List(in.ProductId, int(in.Page), int(in.Size))
	if err != nil {
		return nil, err
	}
	var list []*mall.CommentInfo
	for _, c := range comments {
		list = append(list, toProtoComment(&c))
	}
	return &mall.CommentListResp{List: list, Total: total}, nil
}

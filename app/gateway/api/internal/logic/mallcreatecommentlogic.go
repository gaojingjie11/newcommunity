package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallCreateCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallCreateCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallCreateCommentLogic {
	return &MallCreateCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallCreateCommentLogic) MallCreateComment(req *types.CreateCommentReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.CreateComment(l.ctx, &mall.CreateCommentReq{
		UserId:    getUserIDFromCtx(l.ctx),
		ProductId: req.ProductId,
		Content:   req.Content,
		Rating:    req.Rating,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}

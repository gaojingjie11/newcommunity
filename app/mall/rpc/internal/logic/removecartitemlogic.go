package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveCartItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveCartItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveCartItemLogic {
	return &RemoveCartItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveCartItemLogic) RemoveCartItem(in *mall.RemoveCartItemReq) (*mall.BaseResp, error) {
	err := l.svcCtx.CartSvc.Remove(in.Id, in.UserId)
	if err != nil {
		return &mall.BaseResp{Code: 500, Message: err.Error()}, nil
	}
	return &mall.BaseResp{Code: 0, Message: "success"}, nil
}

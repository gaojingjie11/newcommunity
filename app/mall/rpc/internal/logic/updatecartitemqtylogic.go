package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCartItemQtyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCartItemQtyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCartItemQtyLogic {
	return &UpdateCartItemQtyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateCartItemQtyLogic) UpdateCartItemQty(in *mall.UpdateCartItemQtyReq) (*mall.BaseResp, error) {
	err := l.svcCtx.CartSvc.UpdateQuantity(in.Id, in.UserId, int64(in.Quantity))
	if err != nil {
		return &mall.BaseResp{Code: 500, Message: err.Error()}, nil
	}
	return &mall.BaseResp{Code: 0, Message: "success"}, nil
}

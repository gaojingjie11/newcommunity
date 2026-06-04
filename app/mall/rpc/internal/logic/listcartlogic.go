package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCartLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCartLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCartLogic {
	return &ListCartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListCartLogic) ListCart(in *mall.UserIDReq) (*mall.CartResp, error) {
	items, err := l.svcCtx.CartSvc.List(in.UserId)
	if err != nil {
		return nil, err
	}
	var list []*mall.CartItem
	for _, item := range items {
		list = append(list, toProtoCartItem(&item))
	}
	return &mall.CartResp{Items: list}, nil
}

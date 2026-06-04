package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallCartListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallCartListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallCartListLogic {
	return &MallCartListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallCartListLogic) MallCartList() (resp *types.CartResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.ListCart(l.ctx, &mall.UserIDReq{
		UserId: getUserIDFromCtx(l.ctx),
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.CartItem, 0, len(rpcResp.Items))
	for _, item := range rpcResp.Items {
		list = append(list, toAPICartItem(item))
	}
	return &types.CartResp{
		List: list,
	}, nil
}

package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallAddCartItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallAddCartItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallAddCartItemLogic {
	return &MallAddCartItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallAddCartItemLogic) MallAddCartItem(req *types.AddCartItemReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.AddCartItem(l.ctx, &mall.AddCartItemReq{
		UserId:    getUserIDFromCtx(l.ctx),
		ProductId: req.ProductId,
		Quantity:  req.Quantity,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}

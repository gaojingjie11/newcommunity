package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallUpdateCartItemQtyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallUpdateCartItemQtyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallUpdateCartItemQtyLogic {
	return &MallUpdateCartItemQtyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallUpdateCartItemQtyLogic) MallUpdateCartItemQty(req *types.UpdateCartQtyReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.UpdateCartItemQty(l.ctx, &mall.UpdateCartItemQtyReq{
		UserId:   getUserIDFromCtx(l.ctx),
		Id:       req.Id,
		Quantity: req.Quantity,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}

package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallPayOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallPayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallPayOrderLogic {
	return &MallPayOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallPayOrderLogic) MallPayOrder(req *types.PayOrderReq) (resp *types.PayOrderResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.PayOrder(l.ctx, &mall.PayOrderReq{
		Id:             req.Id,
		UserId:         getUserIDFromCtx(l.ctx),
		PayType:        req.PayType,
		Password:       req.Password,
		FaceImageUrl:   req.FaceImageUrl,
		IdempotencyKey: req.IdempotencyKey,
		ReturnUrl:      req.ReturnUrl,
	})
	if err != nil {
		return nil, err
	}
	return &types.PayOrderResp{
		Success: rpcResp.Success,
		OrderNo: rpcResp.OrderNo,
		PayUrl:  rpcResp.PayUrl,
	}, nil
}

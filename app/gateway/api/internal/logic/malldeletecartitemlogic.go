package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallDeleteCartItemLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallDeleteCartItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallDeleteCartItemLogic {
	return &MallDeleteCartItemLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallDeleteCartItemLogic) MallDeleteCartItem(req *types.CartIDReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.RemoveCartItem(l.ctx, &mall.RemoveCartItemReq{
		UserId: getUserIDFromCtx(l.ctx),
		Id:     req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}

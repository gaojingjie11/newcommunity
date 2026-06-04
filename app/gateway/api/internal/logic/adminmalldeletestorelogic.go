package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminMallDeleteStoreLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminMallDeleteStoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminMallDeleteStoreLogic {
	return &AdminMallDeleteStoreLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminMallDeleteStoreLogic) AdminMallDeleteStore(req *types.AdminDeleteStoreReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.DeleteStore(l.ctx, &mall.StoreIDReq{
		Id: req.Id,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}

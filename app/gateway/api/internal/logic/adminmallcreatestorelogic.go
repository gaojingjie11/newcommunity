package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminMallCreateStoreLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminMallCreateStoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminMallCreateStoreLogic {
	return &AdminMallCreateStoreLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminMallCreateStoreLogic) AdminMallCreateStore(req *types.AdminCreateStoreReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.CreateStore(l.ctx, &mall.CreateStoreReq{
		Name:      req.Name,
		Address:   req.Address,
		Phone:     req.Phone,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}

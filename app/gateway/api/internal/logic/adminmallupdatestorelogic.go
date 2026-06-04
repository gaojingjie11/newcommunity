package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminMallUpdateStoreLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminMallUpdateStoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminMallUpdateStoreLogic {
	return &AdminMallUpdateStoreLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminMallUpdateStoreLogic) AdminMallUpdateStore(req *types.AdminUpdateStoreReq) (resp *types.BaseResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.UpdateStore(l.ctx, &mall.UpdateStoreReq{
		Id:        req.Id,
		Name:      req.Name,
		Address:   req.Address,
		Phone:     req.Phone,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Status:    req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &types.BaseResp{
		Code:    int(rpcResp.Code),
		Message: rpcResp.Message,
	}, nil
}

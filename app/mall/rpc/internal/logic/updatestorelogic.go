package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateStoreLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateStoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStoreLogic {
	return &UpdateStoreLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateStoreLogic) UpdateStore(in *mall.UpdateStoreReq) (*mall.BaseResp, error) {
	err := l.svcCtx.StoreSvc.Update(&model.Store{
		ID:      in.Id,
		Name:    in.Name,
		Address: in.Address,
		Phone:   in.Phone,
	})
	if err != nil {
		return &mall.BaseResp{Code: 500, Message: err.Error()}, nil
	}
	return &mall.BaseResp{Code: 0, Message: "success"}, nil
}

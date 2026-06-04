package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateStoreLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateStoreLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateStoreLogic {
	return &CreateStoreLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateStoreLogic) CreateStore(in *mall.CreateStoreReq) (*mall.BaseResp, error) {
	err := l.svcCtx.StoreSvc.Create(&model.Store{
		Name:    in.Name,
		Address: in.Address,
		Phone:   in.Phone,
	})
	if err != nil {
		return &mall.BaseResp{Code: 500, Message: err.Error()}, nil
	}
	return &mall.BaseResp{Code: 0, Message: "success"}, nil
}

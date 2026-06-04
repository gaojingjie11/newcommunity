package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaymentStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPaymentStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaymentStatusLogic {
	return &GetPaymentStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPaymentStatusLogic) GetPaymentStatus(in *mall.OrderIDReq) (*mall.PaymentStatusResp, error) {
	status, err := l.svcCtx.PaymentSvc.GetPaymentStatus(in.Id)
	if err != nil {
		return nil, err
	}
	return &mall.PaymentStatusResp{Status: int32(status.Status)}, nil
}

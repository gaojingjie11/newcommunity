package logic

import (
	"context"
	"fmt"

	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/service"
	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayOrderLogic {
	return &PayOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PayOrderLogic) PayOrder(in *mall.PayOrderReq) (*mall.PayOrderResp, error) {
	order, err := l.svcCtx.OrderRepo.FindByID(in.Id)
	if err != nil {
		return &mall.PayOrderResp{Success: false}, err
	}

	if in.PayType == "alipay" {
		if l.svcCtx.AlipaySvc == nil {
			return &mall.PayOrderResp{Success: false}, fmt.Errorf("支付宝支付暂未启用配置")
		}

		idempotencyKey := fmt.Sprintf("pay:alipay:%d", in.Id)
		existing, err := l.svcCtx.PaymentRepo.FindByIdempotencyKey(idempotencyKey)
		if err == nil {
			if existing.Status == 1 { // consts.PaymentStatusSuccess
				return &mall.PayOrderResp{Success: true, OrderNo: order.OrderNo}, nil
			}
		} else {
			// Create record
			record := &model.PaymentRecord{
				OrderID:        order.ID,
				OrderNo:        order.OrderNo,
				UserID:         in.UserId,
				Amount:         order.TotalAmount,
				PaymentMethod:  "alipay",
				Status:         0, // consts.PaymentStatusInit
				IdempotencyKey: idempotencyKey,
			}
			if err := l.svcCtx.PaymentRepo.Create(record); err != nil {
				return &mall.PayOrderResp{Success: false}, err
			}
		}

		payURL, err := l.svcCtx.AlipaySvc.GetPaymentURL(order.OrderNo, order.TotalAmount, in.ReturnUrl)
		if err != nil {
			return &mall.PayOrderResp{Success: false}, err
		}
		return &mall.PayOrderResp{Success: true, OrderNo: order.OrderNo, PayUrl: payURL}, nil
	}

	// Default wallet pay path
	// Use frontend-provided idempotency key to avoid duplicate payment
	walletIdempotencyKey := in.IdempotencyKey
	if walletIdempotencyKey == "" {
		walletIdempotencyKey = fmt.Sprintf("pay:wallet:%d:%d", in.Id, in.UserId)
	}

	// Determine auth type - use what frontend provides
	payType := in.PayType
	if payType == "" || payType == "alipay" {
		payType = "password"
	}

	_, err = l.svcCtx.PaymentSvc.PayOrder(in.Id, in.UserId, service.PayOrderRequest{
		PayType:        payType,
		Password:       in.Password,
		FaceImageURL:   in.FaceImageUrl,
		IdempotencyKey: walletIdempotencyKey,
	})
	if err != nil {
		return &mall.PayOrderResp{Success: false}, err
	}
	return &mall.PayOrderResp{Success: true, OrderNo: order.OrderNo}, nil
}


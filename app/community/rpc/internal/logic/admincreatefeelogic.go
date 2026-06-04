package logic

import (
	"context"
	"errors"
	"time"

	"smartcommunity-microservices/app/community/rpc/internal/model"
	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminCreateFeeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminCreateFeeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminCreateFeeLogic {
	return &AdminCreateFeeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminCreateFeeLogic) AdminCreateFee(in *community.CreatePropertyFeeReq) (*community.BaseResp, error) {
	if in.Amount <= 0 {
		return nil, errors.New("amount must be greater than 0")
	}

	amountCents := int64(in.Amount * 100)
	fee := &model.PropertyFee{
		UserID: in.UserId,
		Month:  in.Month,
		Amount: amountCents,
		Status: 0,
	}

	if in.DueDate != "" {
		due, err := time.ParseInLocation("2006-01-02", in.DueDate, time.Local)
		if err != nil {
			due, err = time.ParseInLocation(time.RFC3339, in.DueDate, time.Local)
		}
		if err == nil {
			fee.DueDate = &due
		}
	}

	if err := l.svcCtx.PropertyFeeRepo.Create(fee); err != nil {
		return nil, err
	}

	return &community.BaseResp{Code: 0, Message: "success"}, nil
}

package logic

import (
	"context"
	"errors"
	"fmt"
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

	targetUserID := in.UserId
	// Check if in.UserId is an 11-digit mobile number (Chinese mobile range)
	if in.UserId >= 10000000000 && in.UserId <= 20000000000 {
		var resolvedID int64
		mobileStr := fmt.Sprintf("%d", in.UserId)
		err := l.svcCtx.DB.Table("sys_user").Where("mobile = ?", mobileStr).Pluck("id", &resolvedID).Error
		if err != nil {
			return nil, err
		}
		if resolvedID == 0 {
			return nil, errors.New("未找到手机号为 " + mobileStr + " 的业主用户，请确认手机号是否正确")
		}
		targetUserID = resolvedID
	} else {
		// Verify if the user ID exists
		var resolvedID int64
		err := l.svcCtx.DB.Table("sys_user").Where("id = ?", in.UserId).Pluck("id", &resolvedID).Error
		if err != nil {
			return nil, err
		}
		if resolvedID == 0 {
			return nil, errors.New("未找到ID为该数值的用户")
		}
	}

	amountCents := int64(in.Amount * 100)
	fee := &model.PropertyFee{
		UserID: targetUserID,
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

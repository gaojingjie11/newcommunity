package logic

import (
	"context"
	"fmt"
	"time"

	"smartcommunity-microservices/app/mall/rpc/internal/model"
	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type TransferWalletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTransferWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TransferWalletLogic {
	return &TransferWalletLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *TransferWalletLogic) TransferWallet(in *mall.TransferWalletReq) (*mall.BaseResp, error) {
	var targetUser model.SysUser
	err := l.svcCtx.DB.Where("mobile = ?", in.TargetMobile).First(&targetUser).Error
	if err != nil {
		return &mall.BaseResp{Code: 404, Message: "目标用户不存在"}, nil
	}

	idempotencyKey := fmt.Sprintf("transfer:%d:%d:%d", in.UserId, targetUser.ID, time.Now().UnixNano())
	err = l.svcCtx.WalletSvc.Transfer(in.UserId, targetUser.ID, in.Amount, idempotencyKey)
	if err != nil {
		return &mall.BaseResp{Code: 500, Message: err.Error()}, nil
	}
	return &mall.BaseResp{Code: 0, Message: "success"}, nil
}

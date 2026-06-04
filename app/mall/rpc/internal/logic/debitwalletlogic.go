package logic

import (
	"context"
	"fmt"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type DebitWalletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDebitWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DebitWalletLogic {
	return &DebitWalletLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DebitWalletLogic) DebitWallet(in *mall.DebitWalletReq) (*mall.DebitWalletResp, error) {
	txID, _, _, err := l.svcCtx.WalletSvc.DebitForExternal(
		in.UserId,
		in.Amount,
		in.BusinessType,
		in.OrderNo,
		in.IdempotencyKey,
		"外部服务扣款",
	)
	if err != nil {
		return &mall.DebitWalletResp{Success: false}, err
	}
	transactionNo := fmt.Sprintf("%d", txID)
	return &mall.DebitWalletResp{Success: true, TransactionNo: transactionNo}, nil
}

package logic

import (
	"context"

	"smartcommunity-microservices/app/mall/rpc/internal/svc"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListWalletTransactionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListWalletTransactionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListWalletTransactionsLogic {
	return &ListWalletTransactionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListWalletTransactionsLogic) ListWalletTransactions(in *mall.ListWalletTxReq) (*mall.WalletTxListResp, error) {
	if l.svcCtx.PaymentReconcileSvc != nil {
		l.svcCtx.PaymentReconcileSvc.ReconcilePendingRechargesForUser(l.ctx, in.UserId)
	}
	txs, total, err := l.svcCtx.WalletSvc.ListTransactions(in.UserId, int(in.Page), int(in.Size), nil)
	if err != nil {
		return nil, err
	}
	var list []*mall.WalletTransactionInfo
	for _, t := range txs {
		list = append(list, toProtoWalletTx(&t))
	}
	return &mall.WalletTxListResp{List: list, Total: total}, nil
}

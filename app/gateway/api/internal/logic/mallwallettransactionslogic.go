package logic

import (
	"context"

	"smartcommunity-microservices/app/gateway/api/internal/svc"
	"smartcommunity-microservices/app/gateway/api/internal/types"
	"smartcommunity-microservices/app/mall/rpc/types/mall"

	"github.com/zeromicro/go-zero/core/logx"
)

type MallWalletTransactionsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMallWalletTransactionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MallWalletTransactionsLogic {
	return &MallWalletTransactionsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MallWalletTransactionsLogic) MallWalletTransactions(req *types.ListTransactionsReq) (resp *types.WalletTxListResp, err error) {
	rpcResp, err := l.svcCtx.MallRpc.ListWalletTransactions(l.ctx, &mall.ListWalletTxReq{
		UserId: getUserIDFromCtx(l.ctx),
		Page:   req.Page,
		Size:   req.Size,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.WalletTransactionInfo, 0, len(rpcResp.List))
	for _, t := range rpcResp.List {
		list = append(list, toAPIWalletTransactionInfo(t))
	}
	return &types.WalletTxListResp{
		List:  list,
		Total: rpcResp.Total,
	}, nil
}

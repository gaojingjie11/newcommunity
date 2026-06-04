package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMyPaymentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMyPaymentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyPaymentsLogic {
	return &ListMyPaymentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListMyPaymentsLogic) ListMyPayments(in *community.ListPaymentsReq) (*community.PaymentListResp, error) {
	payments, total, err := l.svcCtx.PropertyFeeRepo.ListPaymentsByUser(in.UserId, int(in.Page), int(in.Size))
	if err != nil {
		return nil, err
	}

	var list []*community.PropertyFeePaymentInfo
	for _, p := range payments {
		paidAtStr := ""
		if p.PaidAt != nil {
			paidAtStr = p.PaidAt.Format("2006-01-02 15:04:05")
		}

		list = append(list, &community.PropertyFeePaymentInfo{
			Id:                  p.ID,
			PropertyFeeId:       p.PropertyFeeID,
			UserId:              p.UserID,
			Amount:              float64(p.Amount) / 100.0,
			WalletTransactionId: p.WalletTransactionID,
			IdempotencyKey:      p.IdempotencyKey,
			Status:              int32(p.Status),
			PaidAt:              paidAtStr,
		})
	}

	return &community.PaymentListResp{
		List:  list,
		Total: total,
	}, nil
}

package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListPaymentsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminListPaymentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListPaymentsLogic {
	return &AdminListPaymentsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminListPaymentsLogic) AdminListPayments(in *community.ListPaymentsReq) (*community.PaymentListResp, error) {
	payments, total, err := l.svcCtx.PropertyFeeRepo.ListPayments(int(in.Page), int(in.Size))
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

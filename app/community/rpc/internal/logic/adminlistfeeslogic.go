package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListFeesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminListFeesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListFeesLogic {
	return &AdminListFeesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminListFeesLogic) AdminListFees(in *community.ListPropertyFeesReq) (*community.PropertyFeeListResp, error) {
	fees, total, err := l.svcCtx.PropertyFeeRepo.ListAll(int(in.Page), int(in.Size))
	if err != nil {
		return nil, err
	}

	var list []*community.PropertyFeeInfo
	for _, f := range fees {
		dueDateStr := ""
		if f.DueDate != nil {
			dueDateStr = f.DueDate.Format("2006-01-02")
		}
		paidAtStr := ""
		if f.PaidAt != nil {
			paidAtStr = f.PaidAt.Format("2006-01-02 15:04:05")
		}

		list = append(list, &community.PropertyFeeInfo{
			Id:        f.ID,
			UserId:    f.UserID,
			Month:     f.Month,
			Amount:    float64(f.Amount) / 100.0,
			Status:    int32(f.Status),
			DueDate:   dueDateStr,
			PaidAt:    paidAtStr,
		})
	}

	return &community.PropertyFeeListResp{
		List:  list,
		Total: total,
	}, nil
}

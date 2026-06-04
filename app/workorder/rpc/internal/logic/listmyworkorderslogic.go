package logic

import (
	"context"

	"smartcommunity-microservices/app/workorder/rpc/internal/svc"
	"smartcommunity-microservices/app/workorder/rpc/types/workorder"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMyWorkordersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMyWorkordersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyWorkordersLogic {
	return &ListMyWorkordersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListMyWorkordersLogic) ListMyWorkorders(in *workorder.ListWorkorderReq) (*workorder.WorkorderListResp, error) {
	orders, total, err := l.svcCtx.WorkorderRepo.ListByUser(in.UserId, int(in.Page), int(in.Size))
	if err != nil {
		return nil, err
	}

	var list []*workorder.WorkorderInfo
	for _, o := range orders {
		procAtStr := ""
		if o.ProcessedAt != nil {
			procAtStr = o.ProcessedAt.Format("2006-01-02 15:04:05")
		}

		list = append(list, &workorder.WorkorderInfo{
			Id:          o.ID,
			Type:        o.Type,
			UserId:      o.UserID,
			Category:    o.Category,
			Description: o.Description,
			Status:      int32(o.Status),
			Result:      o.Result,
			ProcessorId: o.ProcessorID,
			ProcessedAt: procAtStr,
			CreatedAt:   o.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &workorder.WorkorderListResp{
		List:  list,
		Total: total,
	}, nil
}

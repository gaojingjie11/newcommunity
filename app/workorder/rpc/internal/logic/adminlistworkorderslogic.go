package logic

import (
	"context"

	"smartcommunity-microservices/app/workorder/rpc/internal/svc"
	"smartcommunity-microservices/app/workorder/rpc/types/workorder"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminListWorkordersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAdminListWorkordersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminListWorkordersLogic {
	return &AdminListWorkordersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AdminListWorkordersLogic) AdminListWorkorders(in *workorder.ListWorkorderReq) (*workorder.WorkorderListResp, error) {
	// Map status filter
	var statusFilter *int
	if in.Status >= 0 {
		statusVal := int(in.Status)
		statusFilter = &statusVal
	}

	// We don't have type in ListWorkorderReq directly, but let's check.
	// Since in workorder.proto ListWorkorderReq has status, page, size, and user_id, 
	// we pass empty type string and the status filter to ListAll.
	orders, total, err := l.svcCtx.WorkorderRepo.ListAll("", statusFilter, int(in.Page), int(in.Size))
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
			UserName:     o.UserName,
			UserMobile:   o.UserMobile,
		})
	}

	return &workorder.WorkorderListResp{
		List:  list,
		Total: total,
	}, nil
}

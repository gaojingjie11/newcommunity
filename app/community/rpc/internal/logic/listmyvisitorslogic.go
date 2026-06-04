package logic

import (
	"context"

	"smartcommunity-microservices/app/community/rpc/internal/svc"
	"smartcommunity-microservices/app/community/rpc/types/community"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListMyVisitorsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListMyVisitorsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListMyVisitorsLogic {
	return &ListMyVisitorsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListMyVisitorsLogic) ListMyVisitors(in *community.ListVisitorReq) (*community.VisitorListResp, error) {
	visitors, total, err := l.svcCtx.VisitorRepo.ListByUser(in.UserId, int(in.Page), int(in.Size))
	if err != nil {
		return nil, err
	}

	var list []*community.VisitorInfo
	for _, v := range visitors {
		auditAtStr := ""
		if v.AuditAt != nil {
			auditAtStr = v.AuditAt.Format("2006-01-02 15:04:05")
		}

		list = append(list, &community.VisitorInfo{
			Id:           v.ID,
			UserId:       v.UserID,
			VisitorName:  v.VisitorName,
			VisitorPhone: v.VisitorPhone,
			VisitPurpose: v.VisitPurpose,
			ReleaseTime:  v.ReleaseTime.Format("2006-01-02 15:04:05"),
			ValidDate:    v.ValidDate.Format("2006-01-02"),
			Status:       int32(v.Status),
			AuditRemark:  v.AuditRemark,
			AuditAt:      auditAtStr,
			CreatedAt:    v.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &community.VisitorListResp{
		List:  list,
		Total: total,
	}, nil
}

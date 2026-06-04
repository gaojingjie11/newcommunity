package logic

import (
	"context"

	"smartcommunity-microservices/app/stats/rpc/internal/svc"
	"smartcommunity-microservices/app/stats/rpc/types/stats"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAIReportsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAIReportsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAIReportsLogic {
	return &ListAIReportsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListAIReportsLogic) ListAIReports(in *stats.ListReportsReq) (*stats.ReportListResp, error) {
	reports, total, err := l.svcCtx.ReportSvc.ListReports(int(in.Page), int(in.Size))
	if err != nil {
		return nil, err
	}

	var list []*stats.AIReportInfo
	for _, r := range reports {
		list = append(list, &stats.AIReportInfo{
			Id:                 r.ID,
			RepairNewCount:     r.RepairNewCount,
			RepairPendingCount: r.RepairPendingCount,
			VisitorNewCount:    r.VisitorNewCount,
			PropertyPaidCount:  r.PropertyPaidCount,
			PropertyPaidAmount: r.PropertyPaidAmount,
			ReportSummary:      r.ReportSummary,
			ReportMarkdown:     r.Report,
			GeneratedBy:        r.GeneratedBy,
			CreatedAt:          r.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return &stats.ReportListResp{
		List:  list,
		Total: total,
	}, nil
}

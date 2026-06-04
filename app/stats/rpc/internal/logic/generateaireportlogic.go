package logic

import (
	"context"

	"smartcommunity-microservices/app/stats/rpc/internal/svc"
	"smartcommunity-microservices/app/stats/rpc/types/stats"

	"github.com/zeromicro/go-zero/core/logx"
)

type GenerateAIReportLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateAIReportLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateAIReportLogic {
	return &GenerateAIReportLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GenerateAIReportLogic) GenerateAIReport(in *stats.GenerateReportReq) (*stats.ReportResp, error) {
	report, err := l.svcCtx.ReportSvc.GenerateReport(in.UserId)
	if err != nil {
		return nil, err
	}

	return &stats.ReportResp{
		Report: &stats.AIReportInfo{
			Id:                 report.ID,
			RepairNewCount:     report.RepairNewCount,
			RepairPendingCount: report.RepairPendingCount,
			VisitorNewCount:    report.VisitorNewCount,
			PropertyPaidCount:  report.PropertyPaidCount,
			PropertyPaidAmount: report.PropertyPaidAmount,
			ReportSummary:      report.ReportSummary,
			ReportMarkdown:     report.Report,
			GeneratedBy:        report.GeneratedBy,
			CreatedAt:          report.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}

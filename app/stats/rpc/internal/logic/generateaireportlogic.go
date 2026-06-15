package logic

import (
	"context"
	"encoding/json"
	"time"

	"smartcommunity-microservices/app/stats/rpc/internal/model"
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
	// 1. Create a placeholder report record immediately
	report := &model.AIReport{
		ReportSummary: "报告正在生成中，请稍后刷新查看...",
		Report:        "## 💡 研报正在生成中\n\nAI 正在深度分析近7日的数据指标，预计需要 1-2 分钟，生成完成后该页内容将自动更新。请关闭此抽屉，稍后刷新报告列表查看完整内容。",
		GeneratedBy:   in.UserId,
	}

	if err := l.svcCtx.ReportRepo.Create(report); err != nil {
		l.Logger.Errorf("failed to create placeholder AI report: %v", err)
		return nil, err
	}

	// 2. Publish task to MQ asynchronously
	if l.svcCtx.MQ != nil {
		task := struct {
			ReportID int64 `json:"report_id"`
			UserID   int64 `json:"user_id"`
		}{
			ReportID: report.ID,
			UserID:   in.UserId,
		}

		body, marshalErr := json.Marshal(task)
		if marshalErr == nil {
			pubCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if pubErr := l.svcCtx.MQ.PublishEvent(pubCtx, "ai_report.generate", body); pubErr != nil {
				l.Logger.Errorf("failed to publish ai_report.generate task: %v", pubErr)
			} else {
				l.Logger.Infof("published ai_report.generate task successfully for report %d", report.ID)
			}
		} else {
			l.Logger.Errorf("failed to marshal ai_report.generate task: %v", marshalErr)
		}
	} else {
		l.Logger.Infof("RabbitMQ is nil, cannot publish async report task. Executing fallback background goroutine.")
		go func() {
			_ = l.svcCtx.ReportSvc.GenerateReportAsync(report.ID, in.UserId)
		}()
	}

	// 3. Return the placeholder record immediately
	return &stats.ReportResp{
		Report: &stats.AIReportInfo{
			Id:                 report.ID,
			RepairNewCount:     0,
			RepairPendingCount: 0,
			VisitorNewCount:    0,
			PropertyPaidCount:  0,
			PropertyPaidAmount: 0.0,
			ReportSummary:      report.ReportSummary,
			ReportMarkdown:     report.Report,
			GeneratedBy:        report.GeneratedBy,
			CreatedAt:          report.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	}, nil
}


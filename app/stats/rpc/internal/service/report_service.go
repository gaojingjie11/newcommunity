package service

import (
	"fmt"
	"log"
	"strings"

	"smartcommunity-microservices/app/stats/rpc/internal/model"
	"smartcommunity-microservices/app/stats/rpc/internal/repository"
)

type ReportService struct {
	repo  *repository.ReportRepo
	aiSvc *AIService
}

func NewReportService(repo *repository.ReportRepo, aiSvc *AIService) *ReportService {
	return &ReportService{repo: repo, aiSvc: aiSvc}
}

func (s *ReportService) GenerateReport(operatorID int64) (*model.AIReport, error) {
	repairNew, _ := s.repo.Count7DayRepairs()
	repairPending, _ := s.repo.CountPendingRepairs()
	visitorNew, _ := s.repo.Count7DayVisitors()
	feePaid, _ := s.repo.Count7DayPaidFees()
	feeAmount, _ := s.repo.Sum7DayPaidAmount()

	report := &model.AIReport{
		RepairNewCount:     repairNew,
		RepairPendingCount: repairPending,
		VisitorNewCount:    visitorNew,
		PropertyPaidCount:  feePaid,
		PropertyPaidAmount: feeAmount,
		GeneratedBy:        operatorID,
	}

	prompt := fmt.Sprintf(
		"你是一个高级社区物业经理。以下是本社区近7天的数据：报修新增%d条，未处理%d条；访客新增%d条；物业费缴费%d笔，收缴%.2f元。请用 Markdown 生成一份专业的数据分析报告，包含：1）核心数据概览；2）管理风险；3）可执行建议。语言简洁，条理清晰。",
		repairNew, repairPending, visitorNew, feePaid, feeAmount,
	)

	reportText, err := s.aiSvc.GenerateReport(prompt)
	if err != nil {
		log.Printf("generate AI report failed: %v, using fallback", err)
		reportText = buildFallbackReport(report)
	}
	report.Report = normalizeMarkdown(reportText)
	if report.Report == "" {
		report.Report = buildFallbackReport(report)
	}
	report.ReportSummary = buildSummary(report.Report)

	if err := s.repo.Create(report); err != nil {
		return nil, err
	}
	return report, nil
}

func (s *ReportService) GetLatestReport() (*model.AIReport, error) {
	return s.repo.FindLatest()
}

func (s *ReportService) GetReportDetail(id int64) (*model.AIReport, error) {
	return s.repo.FindByID(id)
}

func (s *ReportService) ListReports(page, size int) ([]model.AIReport, int64, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 || size > 50 {
		size = 10
	}
	return s.repo.List(page, size)
}

func buildFallbackReport(r *model.AIReport) string {
	return fmt.Sprintf(`# 社区运营周报

## 核心数据概览

| 指标 | 数值 |
|------|------|
| 报修新增 | %d 条 |
| 待处理报修 | %d 条 |
| 访客新增 | %d 条 |
| 物业费缴费 | %d 笔 |
| 收缴金额 | %.2f 元 |

## 管理风险

- 待处理报修 %d 条，需关注处理效率

## 建议

- 加强报修工单跟进，缩短处理周期
- 定期发布社区通知，提升业主满意度`, r.RepairNewCount, r.RepairPendingCount, r.VisitorNewCount, r.PropertyPaidCount, r.PropertyPaidAmount, r.RepairPendingCount)
}

func normalizeMarkdown(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return text
	}
	if strings.HasPrefix(text, "```markdown") {
		text = strings.TrimPrefix(text, "```markdown")
		text = strings.TrimSuffix(text, "```")
		text = strings.TrimSpace(text)
	} else if strings.HasPrefix(text, "```") {
		text = strings.TrimPrefix(text, "```")
		text = strings.TrimSuffix(text, "```")
		text = strings.TrimSpace(text)
	}
	return text
}

func buildSummary(report string) string {
	lines := strings.Split(report, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "|") && !strings.HasPrefix(line, "-") {
			if len(line) > 200 {
				return line[:200] + "..."
			}
			return line
		}
	}
	return "社区运营数据分析报告"
}

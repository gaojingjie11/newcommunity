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
	complaintsNew, _ := s.repo.Count7DayComplaints()
	mallOrders, mallAmount, _ := s.repo.Count7DayMallOrders()
	mallAmountYuan := float64(mallAmount) / 100.0

	report := &model.AIReport{
		RepairNewCount:     repairNew,
		RepairPendingCount: repairPending,
		VisitorNewCount:    visitorNew,
		PropertyPaidCount:  feePaid,
		PropertyPaidAmount: feeAmount,
		GeneratedBy:        operatorID,
	}

	prompt := fmt.Sprintf(
		"你是一个高级社区运营总监。以下是近7天的社区真实运营指标数据：\n"+
			"【物业与服务】\n"+
			"- 新增报修工单：%d 条，当前待处理报修：%d 条\n"+
			"- 新增投诉建议：%d 条\n"+
			"- 新增登记访客：%d 人次\n"+
			"【物业收费】\n"+
			"- 物业费成功收缴：%d 笔，收缴总额：%.2f 元\n"+
			"【社区商城】\n"+
			"- 商城订单成交：%d 笔，商城总交易额：%.2f 元\n\n"+
			"请基于上述真实数据生成一份详尽、专业的「智能社区深度分析运营研报」，并包含以下排版优美的章节：\n"+
			"1. 📈 社区运营核心指标概览 (请用 Markdown 表格汇总上述数据，对比分析核心变化)\n"+
			"2. ⚠️ 潜在管理风险剖析 (结合新增报修及未处理报修、投诉量分析物业服务和安全的潜在隐患)\n"+
			"3. 💰 物业费收缴与社区商城运营表现深度点评 (对收缴率及商业消费活力展开评估)\n"+
			"4. 💡 针对性运营改善建议与可执行措施 (从服务效率、租户关系、商业活动等维度，提供切实具体的策略)\n\n"+
			"要求：语言严谨、商业感强，使用标准的 Markdown 语法以便于精细排版，具有金融/管理研报级别的厚重感与专业度。内容要详细丰富，多角度深入分析，字数可以多一些。",
		repairNew, repairPending, complaintsNew, visitorNew,
		feePaid, feeAmount,
		mallOrders, mallAmountYuan,
	)

	reportText, err := s.aiSvc.GenerateReport(prompt)
	if err != nil {
		log.Printf("generate AI report failed: %v, using fallback", err)
		reportText = buildFallbackReport(report, complaintsNew, mallOrders, mallAmountYuan)
	}
	report.Report = normalizeMarkdown(reportText)
	if report.Report == "" {
		report.Report = buildFallbackReport(report, complaintsNew, mallOrders, mallAmountYuan)
	}
	report.ReportSummary = buildSummary(report.Report)

	if err := s.repo.Create(report); err != nil {
		return nil, err
	}
	return report, nil
}

func (s *ReportService) GenerateReportAsync(reportID int64, operatorID int64) error {
	report, err := s.repo.FindByID(reportID)
	if err != nil {
		return err
	}

	repairNew, _ := s.repo.Count7DayRepairs()
	repairPending, _ := s.repo.CountPendingRepairs()
	visitorNew, _ := s.repo.Count7DayVisitors()
	feePaid, _ := s.repo.Count7DayPaidFees()
	feeAmount, _ := s.repo.Sum7DayPaidAmount()
	complaintsNew, _ := s.repo.Count7DayComplaints()
	mallOrders, mallAmount, _ := s.repo.Count7DayMallOrders()
	mallAmountYuan := float64(mallAmount) / 100.0

	report.RepairNewCount = repairNew
	report.RepairPendingCount = repairPending
	report.VisitorNewCount = visitorNew
	report.PropertyPaidCount = feePaid
	report.PropertyPaidAmount = feeAmount
	report.GeneratedBy = operatorID

	prompt := fmt.Sprintf(
		"你是一个高级社区运营总监。以下是近7天的社区真实运营指标数据：\n"+
			"【物业与服务】\n"+
			"- 新增报修工单：%d 条，当前待处理报修：%d 条\n"+
			"- 新增投诉建议：%d 条\n"+
			"- 新增登记访客：%d 人次\n"+
			"【物业收费】\n"+
			"- 物业费成功收缴：%d 笔，收缴总额：%.2f 元\n"+
			"【社区商城】\n"+
			"- 商城订单成交：%d 笔，商城总交易额：%.2f 元\n\n"+
			"请基于上述真实数据生成一份详尽、专业的「智能社区深度分析运营研报」，并包含以下排版优美的章节：\n"+
			"1. 📈 社区运营核心指标概览 (请用 Markdown 表格汇总上述数据，对比分析核心变化)\n"+
			"2. ⚠️ 潜在管理风险剖析 (结合新增报修及未处理报修、投诉量分析物业服务和安全的潜在隐患)\n"+
			"3. 💰 物业费收缴与社区商城运营表现深度点评 (对收缴率及商业消费活力展开评估)\n"+
			"4. 💡 针对性运营改善建议与可执行措施 (从服务效率、租户关系、商业活动等维度，提供切实具体的策略)\n\n"+
			"要求：语言严谨、商业感强，使用标准的 Markdown 语法以便于精细排版，具有金融/管理研报级别的厚重感与专业度。内容要详细丰富，多角度深入分析，字数可以多一些。",
		repairNew, repairPending, complaintsNew, visitorNew,
		feePaid, feeAmount,
		mallOrders, mallAmountYuan,
	)

	reportText, err := s.aiSvc.GenerateReport(prompt)
	if err != nil {
		log.Printf("generate AI report failed: %v, using fallback", err)
		reportText = buildFallbackReport(report, complaintsNew, mallOrders, mallAmountYuan)
	}
	report.Report = normalizeMarkdown(reportText)
	if report.Report == "" {
		report.Report = buildFallbackReport(report, complaintsNew, mallOrders, mallAmountYuan)
	}
	report.ReportSummary = buildSummary(report.Report)

	return s.repo.Update(report)
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

func buildFallbackReport(r *model.AIReport, complaints int64, mallOrders int64, mallAmount float64) string {
	return fmt.Sprintf(`# 📊 社区运营运营周报 (系统生成)

## 📈 社区运营核心指标概览

| 运营项目 | 指标名称 | 数值 | 状态/建议 |
| :--- | :--- | :--- | :--- |
| **物业服务** | 新增报修工单 | %d 条 | 正常 |
| | 当前未处理工单 | %d 条 | 需加速处理 |
| | 新增投诉建议 | %d 条 | 需及时回访 |
| | 登记来访人员 | %d 人次 | 正常登记 |
| **物业收费** | 成功缴费笔数 | %d 笔 | - |
| | 物业费收缴总额 | %.2f 元 | 稳步收缴中 |
| **社区商城** | 商城交易笔数 | %d 笔 | 商业活跃 |
| | 商城交易总额 | %.2f 元 | 持续增长 |

## ⚠️ 潜在管理风险分析

1. **报修单处理时效问题**：目前存在 **%d 条** 待处理的报修工单。如果长期无法关闭，可能会导致业主满意度下降，引发潜在投诉。
2. **投诉与建议增多**：本周收到 **%d 条** 新增投诉，主要集中在物业日常维护或车位纠纷。客服组应在24小时内主动与投诉业主取得联系，跟进问题关闭。

## 💡 运营管理建议与可执行措施

* **缩短服务响应时效**：建议物业工程部建立工单催办和限时处理机制，对于紧急报修限时2小时内上门，普通报修在24小时内流转关闭。
* **丰富商城促销活动**：当前商城录得订单 **%d 笔**，建议联合社区优质商户推出“周末社区团购”或“绿色环保积分兑换商品”活动，进一步激活商城的商业变现能力与用户粘性。
* **促进物业收缴率**：可通过微信公众号、管家微信群以及楼道公告形式，温馨提醒未缴物业费的业主进行便捷线上支付。`,
		r.RepairNewCount, r.RepairPendingCount, complaints, r.VisitorNewCount,
		r.PropertyPaidCount, r.PropertyPaidAmount, mallOrders, mallAmount,
		r.RepairPendingCount, complaints, mallOrders)
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

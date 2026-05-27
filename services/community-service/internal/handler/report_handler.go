package handler

import (
	"net/http"
	"strconv"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/community-service/internal/service"

	"github.com/gin-gonic/gin"
)

type ReportHandler struct {
	reportSvc *service.ReportService
}

func NewReportHandler(reportSvc *service.ReportService) *ReportHandler {
	return &ReportHandler{reportSvc: reportSvc}
}

func (h *ReportHandler) GenerateReport(c *gin.Context) {
	operatorID := c.GetInt64("userID")
	report, err := h.reportSvc.GenerateReport(operatorID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "生成报告失败: "+err.Error(), nil)
		return
	}
	response.Success(c, report)
}

func (h *ReportHandler) GetLatestReport(c *gin.Context) {
	report, err := h.reportSvc.GetLatestReport()
	if err != nil {
		response.Error(c, http.StatusNotFound, 404, "暂无报告", nil)
		return
	}
	response.Success(c, report)
}

func (h *ReportHandler) ListReports(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	reports, total, err := h.reportSvc.ListReports(page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, "查询失败", nil)
		return
	}
	response.Success(c, gin.H{"list": reports, "total": total, "page": page, "size": size})
}

func (h *ReportHandler) GetReportDetail(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, 400, "ID格式错误", nil)
		return
	}
	report, err := h.reportSvc.GetReportDetail(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, 404, "报告不存在", nil)
		return
	}
	response.Success(c, report)
}

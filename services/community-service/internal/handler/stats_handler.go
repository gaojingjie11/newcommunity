package handler

import (
	"net/http"
	"strconv"

	"smartcommunity-microservices/pkg/response"
	"smartcommunity-microservices/services/community-service/internal/service"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	statsSvc *service.StatsService
}

func NewStatsHandler(statsSvc *service.StatsService) *StatsHandler {
	return &StatsHandler{statsSvc: statsSvc}
}

func (h *StatsHandler) ProductSalesRank(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	ranks, err := h.statsSvc.ProductSalesRank(limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, gin.H{"list": ranks, "total": len(ranks)})
}

func (h *StatsHandler) ProductViewRank(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	ranks, err := h.statsSvc.ProductViewRank(limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, gin.H{"list": ranks, "total": len(ranks)})
}

func (h *StatsHandler) CommunityOverview(c *gin.Context) {
	overview, err := h.statsSvc.CommunityOverview()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, overview)
}

func (h *StatsHandler) OrderStats(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "30"))
	summary, trend, err := h.statsSvc.OrderStatsCombined(days)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}

	response.Success(c, gin.H{"summary": summary, "trend": trend})
}

func (h *StatsHandler) WorkorderStats(c *gin.Context) {
	summary, err := h.statsSvc.WorkorderSummary()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, 500, err.Error(), nil)
		return
	}
	response.Success(c, gin.H{"list": summary, "total": len(summary)})
}

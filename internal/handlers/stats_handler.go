package handlers

import (
	"net/http"

	handlerInterfaces "pr_reviewer_service_go/internal/handlers/interfaces"

	serviceInterfaces "pr_reviewer_service_go/internal/services/interfaces"

	"github.com/gin-gonic/gin"
)

type statsHandler struct {
	svc serviceInterfaces.StatsService
}

func NewStatsHandler(s serviceInterfaces.StatsService) handlerInterfaces.StatsHandler {
	return &statsHandler{svc: s}
}

func (h *statsHandler) GetReviewerStats(c *gin.Context) {
	stats, err := h.svc.ReviewerStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"reviewer_stats": stats})
}

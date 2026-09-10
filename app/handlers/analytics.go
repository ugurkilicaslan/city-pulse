package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"city-pulse/app/models"
)

// TrackEvent — POST /api/1.0/analytics/track
func (h *Handler) TrackEvent(c *gin.Context) {
	var event models.AnalyticsEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if uid, err := userID(c); err == nil {
		event.UserID = uid
	}
	event.IP = c.ClientIP()
	event.UserAgent = c.GetHeader("User-Agent")
	event.CreatedAt = time.Now()

	h.analyticsSvc.TrackEvent(h.ctx, &event)
	c.JSON(http.StatusAccepted, gin.H{"status": "queued"})
}

// AnalyticsSummary — GET /api/1.0/analytics/summary
func (h *Handler) AnalyticsSummary(c *gin.Context) {
	summary, err := h.analyticsSvc.DailySummary(h.ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}
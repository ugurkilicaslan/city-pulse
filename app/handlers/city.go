package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CitySnapshot — GET /api/1.0/city/snapshot?city=istanbul&lang=tr
// Aggregator endpoint: döviz + haberler + oyun fırsatlarını tek sorguda döner
func (h *Handler) CitySnapshot(c *gin.Context) {
	city := c.DefaultQuery("city", "istanbul")
	lang := c.DefaultQuery("lang", "tr")

	snapshot, err := h.citySvc.GetSnapshot(c.Request.Context(), city, lang)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, snapshot)
}

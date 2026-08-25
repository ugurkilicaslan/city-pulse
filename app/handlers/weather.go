package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// WeatherCurrent — GET /api/1.0/weather?city=istanbul
func (h *Handler) WeatherCurrent(c *gin.Context) {
	city := c.DefaultQuery("city", "istanbul")
	data, err := h.weatherSvc.GetWeather(h.ctx, city)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "hava durumu alınamadı: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}

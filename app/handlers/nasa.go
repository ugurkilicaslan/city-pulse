package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// NasaApod — GET /api/1.0/nasa/apod
func (h *Handler) NasaApod(c *gin.Context) {
	data, err := h.nasaSvc.GetApod(h.ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "nasa verisi alinamadi: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
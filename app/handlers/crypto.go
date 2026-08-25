package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CryptoPrices — GET /api/1.0/crypto/prices
func (h *Handler) CryptoPrices(c *gin.Context) {
	coins, err := h.cryptoSvc.GetPrices(h.ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "kripto fiyatları alınamadı: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"coins": coins,
		"count": len(coins),
	})
}

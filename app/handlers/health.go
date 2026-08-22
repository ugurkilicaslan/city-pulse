package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Healthz — GET /api/1.0/health
// Docker ve load balancer'ın servisin ayakta olup olmadığını kontrol etmesi için
func (h *Handler) Healthz(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}

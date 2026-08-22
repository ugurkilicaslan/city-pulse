package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"city-pulse/internal/config"
)

// APIKeyMiddleware — Her isteğin X-API-Key header'ını kontrol eder.
// Yanlış veya eksikse 401 döner, handler'a hiç ulaşmaz.
func APIKeyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key == "" || key != config.Config.AuthAPIKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized: missing or invalid X-API-Key header",
			})
			return
		}
		c.Next()
	}
}

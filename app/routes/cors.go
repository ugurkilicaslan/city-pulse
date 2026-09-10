package routes

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var localOrigins = []string{
	"http://localhost:3000",
	"http://127.0.0.1:3000",
	"http://localhost:4200",
	"http://127.0.0.1:4200",
	"http://localhost:8080",
	"http://127.0.0.1:8080",
}

var productionOriginSuffixes = []string{
	".citypulse.com",
	".istanbulum.app",
	".antalyam.app",
}

const allowedRequestHeaders = "Content-Type, Authorization, Accept, Origin, Cache-Control, X-Requested-With, X-API-Key"

// corsMiddleware — Local'de belirli origin'lere, prod'da domain suffix'e göre izin verir
func corsMiddleware(isLocal bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if isOriginAllowed(origin, isLocal) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		c.Writer.Header().Set("Access-Control-Allow-Headers", allowedRequestHeaders)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func isOriginAllowed(origin string, isLocal bool) bool {
	if origin == "" {
		return false
	}

	if isLocal {
		for _, allowed := range localOrigins {
			if origin == allowed {
				return true
			}
		}
	}

	for _, suffix := range productionOriginSuffixes {
		if strings.HasSuffix(origin, suffix) {
			return true
		}
	}

	return false
}

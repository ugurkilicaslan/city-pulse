package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"city-pulse/internal/metrics"
)

// RequestLogger — Renkli HTTP istek logları + metrics kaydı
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		isError := status >= 400

		metrics.Global.Record(method+" "+path, latency.Milliseconds(), isError)

		statusColor := "\033[32m"
		switch {
		case status >= 500:
			statusColor = "\033[31m"
		case status >= 400:
			statusColor = "\033[33m"
		case status >= 300:
			statusColor = "\033[36m"
		}

		fmt.Printf("%s[%s] %-6s %-40s %d  %s\033[0m\n",
			statusColor,
			time.Now().Format("15:04:05"),
			method,
			path,
			status,
			latency.Round(time.Millisecond),
		)
	}
}

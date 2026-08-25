package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"city-pulse/internal/ratelimit"
)

// RateLimit — IP başına rate limiting middleware
// limit: dakikadaki max istek sayısı
func RateLimit(limit int) gin.HandlerFunc {
	limiter := ratelimit.New(limit, time.Minute)

	return func(c *gin.Context) {
		key := c.ClientIP()
		remaining := limiter.Remaining(key)

		c.Header("X-RateLimit-Limit", fmt.Sprint(limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprint(remaining))

		if !limiter.Allow(key) {
			c.Header("Retry-After", "60")
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error":               "çok fazla istek — dakikada en fazla " + fmt.Sprint(limit) + " istek yapabilirsiniz",
				"code":                "RATE_LIMITED",
				"retry_after_seconds": 60,
			})
			return
		}
		c.Next()
	}
}

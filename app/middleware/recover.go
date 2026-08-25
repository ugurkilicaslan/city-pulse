package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery — Panic recovery middleware, stack trace logar ve 500 döner
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				stack := debug.Stack()
				fmt.Printf("\033[31m[PANIC RECOVERED]\nPath: %s %s\nError: %v\nStack:\n%s\033[0m\n",
					c.Request.Method, c.Request.URL.Path, r, stack)

				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": "beklenmedik bir sunucu hatası oluştu",
					"code":  "PANIC_RECOVERED",
				})
			}
		}()
		c.Next()
	}
}

package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger is a simple structured request logger middleware.
// In production this would be replaced by a structured logger (e.g. zap, slog).
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		log.Printf(
			"method=%s path=%s status=%d duration=%s ip=%s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(start),
			c.ClientIP(),
		)
	}
}

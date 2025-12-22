package logger

import (
	"time"

	"github.com/gin-gonic/gin"
)

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()
		fields := map[string]interface{}{
			"status":      status,
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"latency_ms":  latency.Milliseconds(),
			"client_ip":   c.ClientIP(),
			"user_agent":  c.Request.UserAgent(),
			"error_count": len(c.Errors),
		}
		if status >= 500 {
			Error("http_request", fields)
			return
		}
		if status >= 400 {
			Warn("http_request", fields)
			return
		}
		Info("http_request", fields)
	}
}
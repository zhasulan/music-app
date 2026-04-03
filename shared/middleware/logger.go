package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestLogger logs basic request/response details with latency.
func RequestLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)

		rid, _ := c.Get(RequestIDKey)
		log.Info("request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.Int("status", c.Writer.Status()),
			zap.String("client_ip", c.ClientIP()),
			zap.String("request_id", ridString(rid)),
			zap.Float64("duration_ms", float64(duration.Microseconds())/1000),
		)
	}
}

func ridString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

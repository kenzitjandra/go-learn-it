package middlewares

import (
	"log/slog"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func SlogMiddleware(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		reqID, _ := c.Get("request_id")

		logger.Info(
			"HTTP request",
			slog.Int("status", status),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("latency", latency.String()),
			slog.Any("request_id", reqID),
		)
	}
}

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqID := "req-" + strconv.FormatInt(time.Now().UnixNano(), 10)
		c.Set("request_id", reqID)
		c.Writer.Header().Set("X-Request-ID", reqID)
		c.Next()
	}
}

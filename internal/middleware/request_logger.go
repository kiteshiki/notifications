package middleware

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"fandom/notifications/internal/models"
	"fandom/notifications/internal/service"
)

func RequestLogger(logService *service.LogService) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// Calculate response time
		responseTime := time.Since(start).Milliseconds()

		// Extract API key from query params (if present)
		apiKey := c.Query("api")
		// Mask API key for logging (show only first 8 chars)
		if len(apiKey) > 8 {
			apiKey = apiKey[:8] + "..."
		}

		// Create log entry
		logEntry := &models.RequestLog{
			Method:       c.Request.Method,
			Path:         c.Request.URL.Path,
			QueryParams:  c.Request.URL.RawQuery,
			StatusCode:   c.Writer.Status(),
			IPAddress:    c.ClientIP(),
			UserAgent:    c.Request.UserAgent(),
			APIKey:       apiKey,
			ResponseTime: responseTime,
			CreatedAt:    time.Now(),
		}

		// Log asynchronously to avoid blocking the response
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = logService.LogRequest(ctx, logEntry)
		}()
	}
}


package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"fandom/notifications/internal/service"
)

func APIKeyAuth(apiKeyService *service.APIKeyService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Try to get API key from cookie first, then fall back to query parameter
		apiKey, err := c.Cookie("api_key")
		if err != nil || apiKey == "" {
			apiKey = c.Query("api")
		}

		if apiKey == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "API key is required. Provide it as 'api' query parameter or authenticate via /auth page.",
			})
			c.Abort()
			return
		}

		valid, err := apiKeyService.ValidateKey(c.Request.Context(), apiKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
			c.Abort()
			return
		}

		if !valid {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Invalid or inactive API key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}


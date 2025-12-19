package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func MasterKeyAuth(masterKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if masterKey == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Master API key not configured. Set MASTER_API_KEY environment variable.",
			})
			c.Abort()
			return
		}

		// Try to get API key from cookie first, then fall back to query parameter
		apiKey, err := c.Cookie("api_key")
		if err != nil || apiKey == "" {
			apiKey = c.Query("api")
		}

		if apiKey == "" {
			// Redirect to auth page for HTML requests, return JSON error for API requests
			if c.GetHeader("Accept") == "text/html" || c.Request.URL.Path == "/dashboard" {
				c.Redirect(http.StatusFound, "/auth")
				c.Abort()
				return
			}
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Master API key is required. Provide it as 'api' query parameter or authenticate via /auth page.",
			})
			c.Abort()
			return
		}

		if apiKey != masterKey {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Invalid master API key",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

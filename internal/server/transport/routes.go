package transport

import (
	"github.com/gin-gonic/gin"

	"fandom/notifications/internal/database"
	"fandom/notifications/internal/service"
)

func RegisterAdminRoutes(rg *gin.RouterGroup, db *database.DB, apiKeyService *service.APIKeyService) {
	// Admin routes (require master API key via middleware)
	apiKeyHandler := NewAPIKeyHandler(apiKeyService)
	rg.POST("/api-keys", apiKeyHandler.CreateAPIKey)
}

func RegisterDashboardRoutes(rg *gin.RouterGroup, logService *service.LogService) {
	// Dashboard routes (require master API key via middleware)
	dashboardHandler := NewDashboardHandler(logService)
	rg.GET("", dashboardHandler.DashboardPage)
	rg.GET("/logs", dashboardHandler.GetLogs)
	rg.GET("/stats", dashboardHandler.GetStats)
}

func RegisterRoutes(rg *gin.RouterGroup, db *database.DB) {
	// Protected routes (require regular API key via middleware)
	rg.GET("/hello", hello)
	// TODO: Add bookmark routes here
}



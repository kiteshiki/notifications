package server

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"fandom/notifications/internal/config"
	"fandom/notifications/internal/database"
	"fandom/notifications/internal/middleware"
	"fandom/notifications/internal/repository"
	"fandom/notifications/internal/server/transport"
	"fandom/notifications/internal/service"
)

func NewRouter(cfg config.Config, db *database.DB) *gin.Engine {
	gin.SetMode(cfg.GinMode)

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// Swagger UI (public)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Auth routes (public)
	r.GET("/auth", transport.AuthPage)
	r.POST("/auth/set", transport.SetAuthCookie)
	r.POST("/auth/clear", transport.ClearAuthCookie)

	// Initialize services
	apiKeyRepo := repository.NewAPIKeyRepository(db)
	apiKeyService := service.NewAPIKeyService(apiKeyRepo)
	logRepo := repository.NewLogRepository(db)
	logService := service.NewLogService(logRepo)

	// Request logging middleware (applies to all routes except swagger)
	r.Use(middleware.RequestLogger(logService))

	// Admin routes (require master API key)
	admin := r.Group("/")
	admin.Use(middleware.MasterKeyAuth(cfg.MasterAPIKey))
	transport.RegisterAdminRoutes(admin, db, apiKeyService)

	// Dashboard routes (require master API key)
	dashboard := r.Group("/dashboard")
	dashboard.Use(middleware.MasterKeyAuth(cfg.MasterAPIKey))
	transport.RegisterDashboardRoutes(dashboard, logService)

	// Protected API routes (require regular API key)
	api := r.Group("/")
	api.Use(middleware.APIKeyAuth(apiKeyService))
	transport.RegisterRoutes(api, db)

	return r
}

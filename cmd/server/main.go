package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "fandom/notifications/docs"
	"fandom/notifications/internal/config"
	"fandom/notifications/internal/database"
	"fandom/notifications/internal/server"
)

// @title           Notifications API
// @version         1.0
// @description     Simple API built with Gin and Swagger.
// @BasePath        /
// @schemes         http
func main() {
	ctx := context.Background()
	cfg := config.Load()

	// Initialize database connection
	db, err := database.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate(ctx); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	router := server.NewRouter(cfg, db)

	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	log.Printf("Server starting on port %s", cfg.Port)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}

	log.Println("Server exited")
}

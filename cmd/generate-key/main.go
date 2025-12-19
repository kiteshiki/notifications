package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"fandom/notifications/internal/config"
	"fandom/notifications/internal/database"
	"fandom/notifications/internal/repository"
	"fandom/notifications/internal/service"
)

func main() {
	var name string
	flag.StringVar(&name, "name", "Initial API Key", "Name for the API key")
	flag.Parse()

	cfg := config.Load()

	ctx := context.Background()

	// Connect to database
	db, err := database.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize services
	apiKeyRepo := repository.NewAPIKeyRepository(db)
	apiKeyService := service.NewAPIKeyService(apiKeyRepo)

	// Generate API key
	response, err := apiKeyService.CreateAPIKey(ctx, name)
	if err != nil {
		log.Fatalf("Failed to create API key: %v", err)
	}

	fmt.Printf("API Key generated successfully!\n\n")
	fmt.Printf("Name: %s\n", response.Name)
	fmt.Printf("Key:  %s\n", response.Key)
	fmt.Printf("Created at: %s\n\n", response.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println("⚠️  IMPORTANT: Save this key securely. It cannot be retrieved later.")
	fmt.Printf("\nYou can use it like this:\n")
	fmt.Printf("  curl \"http://localhost:8080/hello?api=%s\"\n", response.Key)
}


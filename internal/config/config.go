package config

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	Port         string
	GinMode      string
	DatabaseURL  string
	DatabaseName string
	MasterAPIKey string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	ginMode := os.Getenv("GIN_MODE")
	if ginMode == "" {
		ginMode = "release"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		// Default connection string format: postgres://user:password@host:port/dbname
		user := os.Getenv("DB_USER")
		if user == "" {
			user = "postgres"
		}
		password := os.Getenv("DB_PASSWORD")
		// if password == "" {
		// 	password = "postgres"
		// }
		host := os.Getenv("DB_HOST")
		if host == "" {
			host = "localhost"
		}
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}
		dbName := os.Getenv("DB_NAME")
		if dbName == "" {
			dbName = "fandom_notifications"
		}
		databaseURL = "postgres://" + user + ":" + password + "@" + host + ":" + port + "/" + dbName + "?sslmode=disable"
	}

	databaseName := os.Getenv("DB_NAME")
	if databaseName == "" {
		databaseName = "fandom_notifications"
	}

	masterAPIKey := os.Getenv("MASTER_API_KEY")

	return Config{
		Port:         port,
		GinMode:      ginMode,
		DatabaseURL:  databaseURL,
		DatabaseName: databaseName,
		MasterAPIKey: masterAPIKey,
	}
}

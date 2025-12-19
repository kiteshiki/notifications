package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, databaseURL string) (*DB, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")

	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}

func (db *DB) Migrate(ctx context.Context) error {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS bookmarks (
		id SERIAL PRIMARY KEY,
		user_id VARCHAR(255) NOT NULL,
		publication_id VARCHAR(255) NOT NULL,
		chapter_id VARCHAR(255) NOT NULL,
		image VARCHAR(255),
		chapter VARCHAR(255),
		volume VARCHAR(255),
		name VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_bookmarks_user_id ON bookmarks(user_id);
	CREATE INDEX IF NOT EXISTS idx_bookmarks_publication_id ON bookmarks(publication_id);

	CREATE TABLE IF NOT EXISTS api_keys (
		id SERIAL PRIMARY KEY,
		key VARCHAR(255) NOT NULL UNIQUE,
		name VARCHAR(255),
		active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		last_used_at TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_api_keys_key ON api_keys(key);
	CREATE INDEX IF NOT EXISTS idx_api_keys_active ON api_keys(active);

	CREATE TABLE IF NOT EXISTS request_logs (
		id SERIAL PRIMARY KEY,
		method VARCHAR(10) NOT NULL,
		path VARCHAR(500) NOT NULL,
		query_params TEXT,
		status_code INTEGER NOT NULL,
		ip_address VARCHAR(45),
		user_agent TEXT,
		api_key VARCHAR(255),
		response_time_ms BIGINT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_request_logs_created_at ON request_logs(created_at DESC);
	CREATE INDEX IF NOT EXISTS idx_request_logs_method ON request_logs(method);
	CREATE INDEX IF NOT EXISTS idx_request_logs_status_code ON request_logs(status_code);
	CREATE INDEX IF NOT EXISTS idx_request_logs_path ON request_logs(path);
	CREATE INDEX IF NOT EXISTS idx_request_logs_api_key ON request_logs(api_key);
	`

	if _, err := db.Pool.Exec(ctx, createTableQuery); err != nil {
		return fmt.Errorf("failed to create tables: %w", err)
	}

	log.Println("Successfully created database tables")
	return nil
}

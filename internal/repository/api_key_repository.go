package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"fandom/notifications/internal/database"
	"fandom/notifications/internal/models"
)

type APIKeyRepository struct {
	db *database.DB
}

func NewAPIKeyRepository(db *database.DB) *APIKeyRepository {
	return &APIKeyRepository{db: db}
}

func (r *APIKeyRepository) Create(ctx context.Context, key, name string) (*models.APIKey, error) {
	query := `
		INSERT INTO api_keys (key, name, active)
		VALUES ($1, $2, TRUE)
		RETURNING id, key, name, active, created_at, last_used_at
	`

	var apiKey models.APIKey
	var lastUsedAt sql.NullTime

	err := r.db.Pool.QueryRow(ctx, query, key, name).Scan(
		&apiKey.ID,
		&apiKey.Key,
		&apiKey.Name,
		&apiKey.Active,
		&apiKey.CreatedAt,
		&lastUsedAt,
	)
	if err != nil {
		return nil, err
	}

	if lastUsedAt.Valid {
		apiKey.LastUsedAt = &lastUsedAt.Time
	}

	return &apiKey, nil
}

func (r *APIKeyRepository) FindByKey(ctx context.Context, key string) (*models.APIKey, error) {
	query := `
		SELECT id, key, name, active, created_at, last_used_at
		FROM api_keys
		WHERE key = $1 AND active = TRUE
	`

	var apiKey models.APIKey
	var lastUsedAt sql.NullTime

	err := r.db.Pool.QueryRow(ctx, query, key).Scan(
		&apiKey.ID,
		&apiKey.Key,
		&apiKey.Name,
		&apiKey.Active,
		&apiKey.CreatedAt,
		&lastUsedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if lastUsedAt.Valid {
		apiKey.LastUsedAt = &lastUsedAt.Time
	}

	return &apiKey, nil
}

func (r *APIKeyRepository) UpdateLastUsed(ctx context.Context, key string) error {
	query := `
		UPDATE api_keys
		SET last_used_at = $1
		WHERE key = $2
	`

	_, err := r.db.Pool.Exec(ctx, query, time.Now(), key)
	return err
}

func (r *APIKeyRepository) List(ctx context.Context) ([]models.APIKey, error) {
	query := `
		SELECT id, key, name, active, created_at, last_used_at
		FROM api_keys
		ORDER BY created_at DESC
	`

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var apiKeys []models.APIKey
	for rows.Next() {
		var apiKey models.APIKey
		var lastUsedAt sql.NullTime

		err := rows.Scan(
			&apiKey.ID,
			&apiKey.Key,
			&apiKey.Name,
			&apiKey.Active,
			&apiKey.CreatedAt,
			&lastUsedAt,
		)
		if err != nil {
			return nil, err
		}

		if lastUsedAt.Valid {
			apiKey.LastUsedAt = &lastUsedAt.Time
		}

		apiKeys = append(apiKeys, apiKey)
	}

	return apiKeys, rows.Err()
}


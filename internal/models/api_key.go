package models

import "time"

type APIKey struct {
	ID         int       `json:"id" db:"id"`
	Key        string    `json:"key" db:"key"`
	Name       string    `json:"name,omitempty" db:"name"`
	Active     bool      `json:"active" db:"active"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty" db:"last_used_at"`
}

type CreateAPIKeyRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateAPIKeyResponse struct {
	Key       string    `json:"key"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}


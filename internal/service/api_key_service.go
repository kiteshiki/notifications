package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"fandom/notifications/internal/models"
	"fandom/notifications/internal/repository"
)

type APIKeyService struct {
	repo *repository.APIKeyRepository
}

func NewAPIKeyService(repo *repository.APIKeyRepository) *APIKeyService {
	return &APIKeyService{repo: repo}
}

func (s *APIKeyService) GenerateKey() (string, error) {
	bytes := make([]byte, 32) // 256 bits
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func (s *APIKeyService) CreateAPIKey(ctx context.Context, name string) (*models.CreateAPIKeyResponse, error) {
	key, err := s.GenerateKey()
	if err != nil {
		return nil, err
	}

	apiKey, err := s.repo.Create(ctx, key, name)
	if err != nil {
		return nil, err
	}

	return &models.CreateAPIKeyResponse{
		Key:       apiKey.Key,
		Name:      apiKey.Name,
		CreatedAt: apiKey.CreatedAt,
	}, nil
}

func (s *APIKeyService) ValidateKey(ctx context.Context, key string) (bool, error) {
	apiKey, err := s.repo.FindByKey(ctx, key)
	if err != nil {
		return false, err
	}

	if apiKey == nil {
		return false, nil
	}

	// Update last used timestamp
	_ = s.repo.UpdateLastUsed(ctx, key)

	return true, nil
}


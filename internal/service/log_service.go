package service

import (
	"context"
	"time"

	"fandom/notifications/internal/models"
	"fandom/notifications/internal/repository"
)

type LogService struct {
	repo *repository.LogRepository
}

func NewLogService(repo *repository.LogRepository) *LogService {
	return &LogService{repo: repo}
}

func (s *LogService) LogRequest(ctx context.Context, log *models.RequestLog) error {
	if log.CreatedAt.IsZero() {
		log.CreatedAt = time.Now()
	}
	return s.repo.Create(ctx, log)
}

func (s *LogService) GetLogs(ctx context.Context, params models.LogQueryParams) ([]models.RequestLog, int64, error) {
	return s.repo.List(ctx, params)
}

func (s *LogService) GetStats(ctx context.Context, startDate, endDate time.Time) (*models.LogStats, error) {
	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, 0, -7) // Default to last 7 days
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}
	return s.repo.GetStats(ctx, startDate, endDate)
}


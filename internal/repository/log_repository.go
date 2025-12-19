package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"fandom/notifications/internal/database"
	"fandom/notifications/internal/models"
)

type LogRepository struct {
	db *database.DB
}

func NewLogRepository(db *database.DB) *LogRepository {
	return &LogRepository{db: db}
}

func (r *LogRepository) Create(ctx context.Context, log *models.RequestLog) error {
	query := `
		INSERT INTO request_logs (method, path, query_params, status_code, ip_address, user_agent, api_key, response_time_ms, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id
	`

	err := r.db.Pool.QueryRow(
		ctx,
		query,
		log.Method,
		log.Path,
		log.QueryParams,
		log.StatusCode,
		log.IPAddress,
		log.UserAgent,
		log.APIKey,
		log.ResponseTime,
		log.CreatedAt,
	).Scan(&log.ID)

	return err
}

func (r *LogRepository) List(ctx context.Context, params models.LogQueryParams) ([]models.RequestLog, int64, error) {
	where := []string{"1=1"}
	args := []interface{}{}
	argPos := 1

	if params.Method != "" {
		where = append(where, fmt.Sprintf("method = $%d", argPos))
		args = append(args, params.Method)
		argPos++
	}

	if params.StatusCode > 0 {
		where = append(where, fmt.Sprintf("status_code = $%d", argPos))
		args = append(args, params.StatusCode)
		argPos++
	}

	if params.Path != "" {
		where = append(where, fmt.Sprintf("path LIKE $%d", argPos))
		args = append(args, "%"+params.Path+"%")
		argPos++
	}

	if !params.StartDate.IsZero() {
		where = append(where, fmt.Sprintf("created_at >= $%d", argPos))
		args = append(args, params.StartDate)
		argPos++
	}

	if !params.EndDate.IsZero() {
		where = append(where, fmt.Sprintf("created_at <= $%d", argPos))
		args = append(args, params.EndDate)
		argPos++
	}

	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + where[0]
		for i := 1; i < len(where); i++ {
			whereClause += " AND " + where[i]
		}
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM request_logs %s", whereClause)
	var total int64
	err := r.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	if params.Limit <= 0 {
		params.Limit = 100
	}
	if params.Limit > 1000 {
		params.Limit = 1000
	}

	query := fmt.Sprintf(`
		SELECT id, method, path, query_params, status_code, ip_address, user_agent, api_key, response_time_ms, created_at
		FROM request_logs
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argPos, argPos+1)

	args = append(args, params.Limit, params.Offset)

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []models.RequestLog
	for rows.Next() {
		var log models.RequestLog
		var queryParams sql.NullString
		var userAgent sql.NullString
		var apiKey sql.NullString

		err := rows.Scan(
			&log.ID,
			&log.Method,
			&log.Path,
			&queryParams,
			&log.StatusCode,
			&log.IPAddress,
			&userAgent,
			&apiKey,
			&log.ResponseTime,
			&log.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		if queryParams.Valid {
			log.QueryParams = queryParams.String
		}
		if userAgent.Valid {
			log.UserAgent = userAgent.String
		}
		if apiKey.Valid {
			log.APIKey = apiKey.String
		}

		logs = append(logs, log)
	}

	return logs, total, rows.Err()
}

func (r *LogRepository) GetStats(ctx context.Context, startDate, endDate time.Time) (*models.LogStats, error) {
	stats := &models.LogStats{
		StatusCodes: make(map[int]int64),
	}

	// Total requests
	var total sql.NullInt64
	var avgResponseTime sql.NullFloat64

	query := `
		SELECT 
			COUNT(*) as total,
			COALESCE(AVG(response_time_ms), 0) as avg_response_time
		FROM request_logs
		WHERE created_at >= $1 AND created_at <= $2
	`

	err := r.db.Pool.QueryRow(ctx, query, startDate, endDate).Scan(&total, &avgResponseTime)
	if err != nil && err != pgx.ErrNoRows {
		return nil, err
	}

	if total.Valid {
		stats.TotalRequests = total.Int64
	}
	if avgResponseTime.Valid {
		stats.AverageResponseTime = avgResponseTime.Float64
	}

	// Status code distribution
	statusQuery := `
		SELECT status_code, COUNT(*) as count
		FROM request_logs
		WHERE created_at >= $1 AND created_at <= $2
		GROUP BY status_code
		ORDER BY count DESC
	`

	rows, err := r.db.Pool.Query(ctx, statusQuery, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var code int
		var count int64
		if err := rows.Scan(&code, &count); err != nil {
			return nil, err
		}
		stats.StatusCodes[code] = count
	}

	// Top paths
	pathQuery := `
		SELECT path, COUNT(*) as count
		FROM request_logs
		WHERE created_at >= $1 AND created_at <= $2
		GROUP BY path
		ORDER BY count DESC
		LIMIT 10
	`

	rows, err = r.db.Pool.Query(ctx, pathQuery, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pc models.PathCount
		if err := rows.Scan(&pc.Path, &pc.Count); err != nil {
			return nil, err
		}
		stats.TopPaths = append(stats.TopPaths, pc)
	}

	// Top methods
	methodQuery := `
		SELECT method, COUNT(*) as count
		FROM request_logs
		WHERE created_at >= $1 AND created_at <= $2
		GROUP BY method
		ORDER BY count DESC
	`

	rows, err = r.db.Pool.Query(ctx, methodQuery, startDate, endDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var mc models.MethodCount
		if err := rows.Scan(&mc.Method, &mc.Count); err != nil {
			return nil, err
		}
		stats.TopMethods = append(stats.TopMethods, mc)
	}

	return stats, nil
}


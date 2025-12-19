package models

import "time"

type RequestLog struct {
	ID           int64     `json:"id" db:"id"`
	Method       string    `json:"method" db:"method"`
	Path         string    `json:"path" db:"path"`
	QueryParams  string    `json:"query_params,omitempty" db:"query_params"`
	StatusCode   int       `json:"status_code" db:"status_code"`
	IPAddress    string    `json:"ip_address" db:"ip_address"`
	UserAgent    string    `json:"user_agent,omitempty" db:"user_agent"`
	APIKey       string    `json:"api_key,omitempty" db:"api_key"`
	ResponseTime int64     `json:"response_time_ms" db:"response_time_ms"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type LogQueryParams struct {
	Limit      int       `form:"limit"`
	Page       int       `form:"page"`
	Offset     int       `form:"offset"`
	Method     string    `form:"method"`
	StatusCode int       `form:"status_code"`
	Path       string    `form:"path"`
	StartDate  time.Time `form:"start_date" time_format:"2006-01-02T15:04:05Z07:00"`
	EndDate    time.Time `form:"end_date" time_format:"2006-01-02T15:04:05Z07:00"`
}

type LogStats struct {
	TotalRequests       int64         `json:"total_requests"`
	AverageResponseTime float64       `json:"average_response_time_ms"`
	StatusCodes         map[int]int64 `json:"status_codes"`
	TopPaths            []PathCount   `json:"top_paths"`
	TopMethods          []MethodCount `json:"top_methods"`
}

type PathCount struct {
	Path  string `json:"path"`
	Count int64  `json:"count"`
}

type MethodCount struct {
	Method string `json:"method"`
	Count  int64  `json:"count"`
}

package models

import "time"

type Bookmark struct {
	ID            int       `json:"id" db:"id"`
	UserID        string    `json:"user_id" db:"user_id"`
	PublicationID string    `json:"publication_id" db:"publication_id"`
	ChapterID     string    `json:"chapter_id" db:"chapter_id"`
	Image         string    `json:"image,omitempty" db:"image"`
	Chapter       string    `json:"chapter,omitempty" db:"chapter"`
	Volume        string    `json:"volume,omitempty" db:"volume"`
	Name          string    `json:"name,omitempty" db:"name"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}


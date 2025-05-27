package models

import (
	"database/sql"
	"time"
)

type Progress struct {
	ID                   int          `json:"id"`
	UserID               int          `json:"user_id"`
	ProjectID            int          `json:"project_id"`
	Status               string       `json:"status"`
	CompletionPercentage int          `json:"completion_percentage"`
	StartedAt            sql.NullTime `json:"started_at"`
	CompletedAt          sql.NullTime `json:"completed_at"`
	CreatedAt            time.Time    `json:"created_at"`
	UpdatedAt            time.Time    `json:"updated_at"`
}

type UpdateProgressRequest struct {
	Status               string `json:"status"`
	CompletionPercentage int    `json:"completion_percentage"`
}

const (
	StatusNotStarted = "not_started"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusOnHold     = "on_hold"
	StatusAbandoned  = "abandoned"
)

package models

import (
	"time"
)

type Curriculum struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Projects    []Project `json:"projects,omitempty"`
}

type CreateCurriculumRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateCurriculumRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CurriculumWithStats struct {
	Curriculum
	TotalProjects     int `json:"total_projects"`
	CompletedProjects int `json:"completed_projects"`
	TotalTimeSpent    int `json:"total_time_spent"`
}

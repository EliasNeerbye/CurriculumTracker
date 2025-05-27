package main

import (
	"time"

	"github.com/google/uuid"
)

// Client-side models that mirror the server models but without database dependencies

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Curriculum struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProjectType string

const (
	ProjectTypeRoot            ProjectType = "root"
	ProjectTypeBase            ProjectType = "base"
	ProjectTypeLowerBranch     ProjectType = "lower_branch"
	ProjectTypeMiddleBranch    ProjectType = "middle_branch"
	ProjectTypeUpperBranch     ProjectType = "upper_branch"
	ProjectTypeFlowerMilestone ProjectType = "flower_milestone"
	ProjectTypeRootTest        ProjectType = "root_test"
	ProjectTypeBaseTest        ProjectType = "base_test"
)

type Project struct {
	ID                 uuid.UUID   `json:"id"`
	CurriculumID       uuid.UUID   `json:"curriculum_id"`
	Identifier         string      `json:"identifier"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	LearningObjectives []string    `json:"learning_objectives"`
	EstimatedTime      string      `json:"estimated_time"`
	Prerequisites      []string    `json:"prerequisites"`
	ProjectType        ProjectType `json:"project_type"`
	OrderIndex         int         `json:"order_index"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
}

type ProjectStatus string

const (
	ProjectStatusNotStarted ProjectStatus = "not_started"
	ProjectStatusInProgress ProjectStatus = "in_progress"
	ProjectStatusCompleted  ProjectStatus = "completed"
	ProjectStatusPaused     ProjectStatus = "paused"
)

type ProjectProgress struct {
	ID               uuid.UUID     `json:"id"`
	ProjectID        uuid.UUID     `json:"project_id"`
	UserID           uuid.UUID     `json:"user_id"`
	Status           ProjectStatus `json:"status"`
	TimeSpentMinutes int           `json:"time_spent_minutes"`
	CompletedAt      *time.Time    `json:"completed_at"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
}

type Note struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
	UserID    uuid.UUID `json:"user_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Reflection struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
	UserID    uuid.UUID `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TimeEntry struct {
	ID          uuid.UUID `json:"id"`
	ProjectID   uuid.UUID `json:"project_id"`
	UserID      uuid.UUID `json:"user_id"`
	Minutes     int       `json:"minutes"`
	Description string    `json:"description"`
	LoggedAt    time.Time `json:"logged_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// Request/Response types
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type CurriculumRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type ProjectRequest struct {
	Identifier         string      `json:"identifier"`
	Name               string      `json:"name"`
	Description        string      `json:"description"`
	LearningObjectives []string    `json:"learning_objectives"`
	EstimatedTime      string      `json:"estimated_time"`
	Prerequisites      []string    `json:"prerequisites"`
	ProjectType        ProjectType `json:"project_type"`
	OrderIndex         int         `json:"order_index"`
}

type ProgressUpdateRequest struct {
	Status           ProjectStatus `json:"status"`
	TimeSpentMinutes int           `json:"time_spent_minutes"`
}

type NoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ReflectionRequest struct {
	Content string `json:"content"`
}

type TimeEntryRequest struct {
	Minutes     int    `json:"minutes"`
	Description string `json:"description"`
}

type AnalyticsResponse struct {
	TotalProjects      int                     `json:"total_projects"`
	CompletedProjects  int                     `json:"completed_projects"`
	InProgressProjects int                     `json:"in_progress_projects"`
	TotalTimeSpent     int                     `json:"total_time_spent"`
	WeeklyTimeSpent    []int                   `json:"weekly_time_spent"`
	ProjectsByType     map[ProjectType]int     `json:"projects_by_type"`
	CompletionByType   map[ProjectType]float64 `json:"completion_by_type"`
	RecentActivity     []TimeEntry             `json:"recent_activity"`
}

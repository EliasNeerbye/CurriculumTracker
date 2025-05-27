package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Name         string    `json:"name" db:"name"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type Curriculum struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
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
	ID                 uuid.UUID      `json:"id" db:"id"`
	CurriculumID       uuid.UUID      `json:"curriculum_id" db:"curriculum_id"`
	Identifier         string         `json:"identifier" db:"identifier"`
	Name               string         `json:"name" db:"name"`
	Description        string         `json:"description" db:"description"`
	LearningObjectives pq.StringArray `json:"learning_objectives" db:"learning_objectives"`
	EstimatedTime      string         `json:"estimated_time" db:"estimated_time"`
	Prerequisites      pq.StringArray `json:"prerequisites" db:"prerequisites"`
	ProjectType        ProjectType    `json:"project_type" db:"project_type"`
	OrderIndex         int            `json:"order_index" db:"order_index"`
	CreatedAt          time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at" db:"updated_at"`
}

type ProjectStatus string

const (
	StatusNotStarted ProjectStatus = "not_started"
	StatusInProgress ProjectStatus = "in_progress"
	StatusCompleted  ProjectStatus = "completed"
	StatusPaused     ProjectStatus = "paused"
)

type ProjectProgress struct {
	ID               uuid.UUID     `json:"id" db:"id"`
	UserID           uuid.UUID     `json:"user_id" db:"user_id"`
	ProjectID        uuid.UUID     `json:"project_id" db:"project_id"`
	Status           ProjectStatus `json:"status" db:"status"`
	TimeSpentMinutes int           `json:"time_spent_minutes" db:"time_spent_minutes"`
	StartedAt        *time.Time    `json:"started_at" db:"started_at"`
	CompletedAt      *time.Time    `json:"completed_at" db:"completed_at"`
	CreatedAt        time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at" db:"updated_at"`
}

type Note struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	ProjectID uuid.UUID `json:"project_id" db:"project_id"`
	Title     string    `json:"title" db:"title"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type Reflection struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	ProjectID uuid.UUID `json:"project_id" db:"project_id"`
	Content   string    `json:"content" db:"content"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type TimeEntry struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	ProjectID   uuid.UUID `json:"project_id" db:"project_id"`
	Minutes     int       `json:"minutes" db:"minutes"`
	Description string    `json:"description" db:"description"`
	LoggedAt    time.Time `json:"logged_at" db:"logged_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type LoginRequest struct {
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
	ProjectsByType     map[ProjectType]int     `json:"projects_by_type"`
	CompletionByType   map[ProjectType]float64 `json:"completion_by_type"`
	RecentActivity     []TimeEntry             `json:"recent_activity"`
	WeeklyTimeSpent    []int                   `json:"weekly_time_spent"`
}

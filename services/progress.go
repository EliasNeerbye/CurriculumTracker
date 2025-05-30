package services

import (
	"curriculum-tracker/models"
	"database/sql"
	"fmt"
	"time"
)

type ProgressService struct {
	db *sql.DB
}

func NewProgressService(db *sql.DB) *ProgressService {
	return &ProgressService{db: db}
}

func (s *ProgressService) UpdateProgress(userID, projectID int, req models.UpdateProgressRequest) (*models.Progress, error) {
	// Get current progress to determine state transitions
	var currentStatus string
	var currentStartedAt sql.NullTime
	getCurrentQuery := `
		SELECT status, started_at 
		FROM progress 
		WHERE user_id = $1 AND project_id = $2
	`
	err := s.db.QueryRow(getCurrentQuery, userID, projectID).Scan(&currentStatus, &currentStartedAt)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get current progress: %w", err)
	}

	// Determine started_at and completed_at based on state transitions
	var startedAt, completedAt sql.NullTime

	// If transitioning from not_started to any other status, set started_at
	if (currentStatus == "" || currentStatus == models.StatusNotStarted) &&
		req.Status != models.StatusNotStarted {
		startedAt = sql.NullTime{Time: time.Now(), Valid: true}
	} else if currentStartedAt.Valid {
		// Preserve existing started_at
		startedAt = currentStartedAt
	}

	// Set completed_at only when transitioning to completed
	if req.Status == models.StatusCompleted {
		completedAt = sql.NullTime{Time: time.Now(), Valid: true}
		req.CompletionPercentage = 100
	}

	// Validate completion percentage
	if req.Status == models.StatusNotStarted {
		req.CompletionPercentage = 0
	} else if req.Status == models.StatusAbandoned && req.CompletionPercentage == 100 {
		req.CompletionPercentage = 99 // Can't be 100% if abandoned
	}

	query := `
		INSERT INTO progress (user_id, project_id, status, completion_percentage, started_at, completed_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, project_id)
		DO UPDATE SET 
			status = EXCLUDED.status,
			completion_percentage = EXCLUDED.completion_percentage,
			started_at = CASE 
				WHEN progress.started_at IS NULL AND EXCLUDED.started_at IS NOT NULL 
				THEN EXCLUDED.started_at 
				ELSE progress.started_at 
			END,
			completed_at = CASE
				WHEN EXCLUDED.status = 'completed' THEN EXCLUDED.completed_at
				ELSE NULL
			END,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, user_id, project_id, status, completion_percentage, started_at, completed_at, created_at, updated_at
	`

	var progress models.Progress
	err = s.db.QueryRow(query, userID, projectID, req.Status, req.CompletionPercentage, startedAt, completedAt).Scan(
		&progress.ID, &progress.UserID, &progress.ProjectID, &progress.Status,
		&progress.CompletionPercentage, &progress.StartedAt, &progress.CompletedAt,
		&progress.CreatedAt, &progress.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update progress: %w", err)
	}

	return &progress, nil
}

func (s *ProgressService) GetProgressByProjectID(userID, projectID int) (*models.Progress, error) {
	query := `
		SELECT id, user_id, project_id, status, completion_percentage, started_at, completed_at, created_at, updated_at
		FROM progress
		WHERE user_id = $1 AND project_id = $2
	`

	var progress models.Progress
	err := s.db.QueryRow(query, userID, projectID).Scan(
		&progress.ID, &progress.UserID, &progress.ProjectID, &progress.Status,
		&progress.CompletionPercentage, &progress.StartedAt, &progress.CompletedAt,
		&progress.CreatedAt, &progress.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			// Return a default progress object instead of error
			return &models.Progress{
				UserID:               userID,
				ProjectID:            projectID,
				Status:               models.StatusNotStarted,
				CompletionPercentage: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to query progress: %w", err)
	}

	return &progress, nil
}

func (s *ProgressService) GetProgressByCurriculumID(userID, curriculumID int) ([]models.Progress, error) {
	query := `
		SELECT pr.id, pr.user_id, pr.project_id, pr.status, pr.completion_percentage, 
		       pr.started_at, pr.completed_at, pr.created_at, pr.updated_at
		FROM progress pr
		JOIN projects p ON pr.project_id = p.id
		JOIN curricula c ON p.curriculum_id = c.id
		WHERE pr.user_id = $1 AND c.id = $2 AND c.user_id = $1
		ORDER BY p.position_order, p.created_at
	`

	rows, err := s.db.Query(query, userID, curriculumID)
	if err != nil {
		return nil, fmt.Errorf("failed to query progress: %w", err)
	}
	defer rows.Close()

	progressList := make([]models.Progress, 0)
	for rows.Next() {
		var progress models.Progress
		err := rows.Scan(
			&progress.ID, &progress.UserID, &progress.ProjectID, &progress.Status,
			&progress.CompletionPercentage, &progress.StartedAt, &progress.CompletedAt,
			&progress.CreatedAt, &progress.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan progress: %w", err)
		}
		progressList = append(progressList, progress)
	}

	return progressList, nil
}

func (s *ProgressService) CanStartProject(userID, projectID int) (bool, error) {
	// Check if all prerequisites are completed
	query := `
		SELECT COUNT(*)
		FROM projects p
		JOIN curricula c ON p.curriculum_id = c.id
		JOIN unnest(p.prerequisites) AS prereq_id ON true
		JOIN projects prereq ON prereq.identifier = prereq_id AND prereq.curriculum_id = p.curriculum_id
		LEFT JOIN progress pr ON pr.project_id = prereq.id AND pr.user_id = $1
		WHERE p.id = $2 AND c.user_id = $1
		AND (pr.status IS NULL OR pr.status != 'completed')
	`

	var incompletePrereqs int
	err := s.db.QueryRow(query, userID, projectID).Scan(&incompletePrereqs)
	if err != nil {
		return false, fmt.Errorf("failed to check prerequisites: %w", err)
	}

	return incompletePrereqs == 0, nil
}

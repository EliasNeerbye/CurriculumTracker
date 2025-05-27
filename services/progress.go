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
	var startedAt, completedAt sql.NullTime

	if req.Status == models.StatusInProgress {
		startedAt = sql.NullTime{Time: time.Now(), Valid: true}
	}
	if req.Status == models.StatusCompleted {
		completedAt = sql.NullTime{Time: time.Now(), Valid: true}
		req.CompletionPercentage = 100
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
			completed_at = EXCLUDED.completed_at,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, user_id, project_id, status, completion_percentage, started_at, completed_at, created_at, updated_at
	`

	var progress models.Progress
	err := s.db.QueryRow(query, userID, projectID, req.Status, req.CompletionPercentage, startedAt, completedAt).Scan(
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
			return nil, fmt.Errorf("progress not found")
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
		ORDER BY p.project_type, p.position_order
	`

	rows, err := s.db.Query(query, userID, curriculumID)
	if err != nil {
		return nil, fmt.Errorf("failed to query progress: %w", err)
	}
	defer rows.Close()

	var progressList []models.Progress
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

package services

import (
	"curriculum-tracker/models"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type ProjectService struct {
	db *sql.DB
}

func NewProjectService(db *sql.DB) *ProjectService {
	return &ProjectService{db: db}
}

func (s *ProjectService) CreateProject(curriculumID int, req models.CreateProjectRequest) (*models.Project, error) {
	query := `
		INSERT INTO projects (curriculum_id, identifier, name, description, learning_objectives, estimated_time, prerequisites, project_type, position_order)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, curriculum_id, identifier, name, description, learning_objectives, estimated_time, prerequisites, project_type, position_order, created_at, updated_at
	`

	var project models.Project
	err := s.db.QueryRow(query, curriculumID, req.Identifier, req.Name, req.Description,
		pq.Array(req.LearningObjectives), req.EstimatedTime, pq.Array(req.Prerequisites),
		req.ProjectType, req.PositionOrder).Scan(
		&project.ID, &project.CurriculumID, &project.Identifier, &project.Name,
		&project.Description, pq.Array(&project.LearningObjectives), &project.EstimatedTime,
		pq.Array(&project.Prerequisites), &project.ProjectType, &project.PositionOrder,
		&project.CreatedAt, &project.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create project: %w", err)
	}

	return &project, nil
}

func (s *ProjectService) GetProjectsByCurriculumID(userID, curriculumID int) ([]models.Project, error) {
	query := `
		SELECT 
			p.id, p.curriculum_id, p.identifier, p.name, p.description, 
			p.learning_objectives, p.estimated_time, p.prerequisites, 
			p.project_type, p.position_order, p.created_at, p.updated_at,
			pr.id, pr.user_id, pr.project_id, pr.status, pr.completion_percentage,
			pr.started_at, pr.completed_at, pr.created_at, pr.updated_at
		FROM projects p
		INNER JOIN curricula c ON p.curriculum_id = c.id
		LEFT JOIN progress pr ON p.id = pr.project_id AND pr.user_id = $1
		WHERE p.curriculum_id = $2 AND c.user_id = $1
		ORDER BY p.project_type, p.position_order
	`

	rows, err := s.db.Query(query, userID, curriculumID)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		var progress models.Progress
		var progressID sql.NullInt64

		err := rows.Scan(
			&p.ID, &p.CurriculumID, &p.Identifier, &p.Name, &p.Description,
			pq.Array(&p.LearningObjectives), &p.EstimatedTime, pq.Array(&p.Prerequisites),
			&p.ProjectType, &p.PositionOrder, &p.CreatedAt, &p.UpdatedAt,
			&progressID, &progress.UserID, &progress.ProjectID, &progress.Status,
			&progress.CompletionPercentage, &progress.StartedAt, &progress.CompletedAt,
			&progress.CreatedAt, &progress.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}

		if progressID.Valid {
			progress.ID = int(progressID.Int64)
			p.Progress = &progress
		}

		projects = append(projects, p)
	}

	return projects, nil
}

func (s *ProjectService) GetProjectByID(userID, projectID int) (*models.Project, error) {
	query := `
		SELECT 
			p.id, p.curriculum_id, p.identifier, p.name, p.description, 
			p.learning_objectives, p.estimated_time, p.prerequisites, 
			p.project_type, p.position_order, p.created_at, p.updated_at
		FROM projects p
		INNER JOIN curricula c ON p.curriculum_id = c.id
		WHERE p.id = $1 AND c.user_id = $2
	`

	var project models.Project
	err := s.db.QueryRow(query, projectID, userID).Scan(
		&project.ID, &project.CurriculumID, &project.Identifier, &project.Name,
		&project.Description, pq.Array(&project.LearningObjectives), &project.EstimatedTime,
		pq.Array(&project.Prerequisites), &project.ProjectType, &project.PositionOrder,
		&project.CreatedAt, &project.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found")
		}
		return nil, fmt.Errorf("failed to query project: %w", err)
	}

	return &project, nil
}

func (s *ProjectService) UpdateProject(userID, projectID int, req models.UpdateProjectRequest) (*models.Project, error) {
	query := `
		UPDATE projects
		SET identifier = $1, name = $2, description = $3, learning_objectives = $4,
		    estimated_time = $5, prerequisites = $6, project_type = $7, position_order = $8,
		    updated_at = CURRENT_TIMESTAMP
		FROM curricula c
		WHERE projects.id = $9 AND projects.curriculum_id = c.id AND c.user_id = $10
		RETURNING projects.id, projects.curriculum_id, projects.identifier, projects.name, 
		         projects.description, projects.learning_objectives, projects.estimated_time, 
		         projects.prerequisites, projects.project_type, projects.position_order, 
		         projects.created_at, projects.updated_at
	`

	var project models.Project
	err := s.db.QueryRow(query, req.Identifier, req.Name, req.Description,
		pq.Array(req.LearningObjectives), req.EstimatedTime, pq.Array(req.Prerequisites),
		req.ProjectType, req.PositionOrder, projectID, userID).Scan(
		&project.ID, &project.CurriculumID, &project.Identifier, &project.Name,
		&project.Description, pq.Array(&project.LearningObjectives), &project.EstimatedTime,
		pq.Array(&project.Prerequisites), &project.ProjectType, &project.PositionOrder,
		&project.CreatedAt, &project.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project not found")
		}
		return nil, fmt.Errorf("failed to update project: %w", err)
	}

	return &project, nil
}

func (s *ProjectService) DeleteProject(userID, projectID int) error {
	query := `
		DELETE FROM projects
		USING curricula c
		WHERE projects.id = $1 AND projects.curriculum_id = c.id AND c.user_id = $2
	`

	result, err := s.db.Exec(query, projectID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("project not found")
	}

	return nil
}

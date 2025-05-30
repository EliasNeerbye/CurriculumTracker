package services

import (
	"curriculum-tracker/models"
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type ProjectService struct {
	db *sql.DB
}

func NewProjectService(db *sql.DB) *ProjectService {
	return &ProjectService{db: db}
}

func (s *ProjectService) generateIdentifier(curriculumID int, projectType string) (string, error) {
	switch projectType {
	case models.ProjectTypeRoot:
		return s.generateSequentialIdentifier(curriculumID, projectType, "R")
	case models.ProjectTypeRootTest:
		return s.generateTestIdentifier(curriculumID, projectType, "RT")
	case models.ProjectTypeBase:
		return s.generateSequentialIdentifier(curriculumID, projectType, "B")
	case models.ProjectTypeBaseTest:
		return s.generateTestIdentifier(curriculumID, projectType, "BT")
	case models.ProjectTypeLowerBranch:
		return s.generateBranchIdentifier(curriculumID, projectType, "LB")
	case models.ProjectTypeMiddleBranch:
		return s.generateBranchIdentifier(curriculumID, projectType, "MB")
	case models.ProjectTypeUpperBranch:
		return s.generateBranchIdentifier(curriculumID, projectType, "UB")
	case models.ProjectTypeFlowerMilestone:
		return s.generateSequentialIdentifier(curriculumID, projectType, "F")
	default:
		return "", fmt.Errorf("unknown project type: %s", projectType)
	}
}

func (s *ProjectService) generateSequentialIdentifier(curriculumID int, projectType, prefix string) (string, error) {
	query := `
		SELECT COUNT(*) + 1 
		FROM projects 
		WHERE curriculum_id = $1 AND project_type = $2
	`

	var nextNum int
	err := s.db.QueryRow(query, curriculumID, projectType).Scan(&nextNum)
	if err != nil {
		return "", fmt.Errorf("failed to generate identifier: %w", err)
	}

	return fmt.Sprintf("%s%d", prefix, nextNum), nil
}

func (s *ProjectService) generateTestIdentifier(curriculumID int, projectType, prefix string) (string, error) {
	// Check if test project already exists
	query := `
		SELECT COUNT(*) 
		FROM projects 
		WHERE curriculum_id = $1 AND project_type = $2
	`

	var count int
	err := s.db.QueryRow(query, curriculumID, projectType).Scan(&count)
	if err != nil {
		return "", fmt.Errorf("failed to check existing test projects: %w", err)
	}

	if count > 0 {
		return "", fmt.Errorf("test project of type %s already exists in this curriculum", projectType)
	}

	return prefix, nil
}

func (s *ProjectService) generateBranchIdentifier(curriculumID int, projectType, prefix string) (string, error) {
	// For now, using simple sequential numbering like LB1, LB2, LB3
	// Can be enhanced later to support grouping like LB1_1, LB1_2, LB2_1, LB2_2
	query := `
		SELECT COUNT(*) + 1 
		FROM projects 
		WHERE curriculum_id = $1 AND project_type = $2
	`

	var nextNum int
	err := s.db.QueryRow(query, curriculumID, projectType).Scan(&nextNum)
	if err != nil {
		return "", fmt.Errorf("failed to generate identifier: %w", err)
	}

	return fmt.Sprintf("%s%d", prefix, nextNum), nil
}

func (s *ProjectService) validatePrerequisites(curriculumID int, prerequisites []string, currentIdentifier string) error {
	if len(prerequisites) == 0 {
		return nil
	}

	// Get all projects in the curriculum with their identifiers and order
	query := `
		SELECT identifier, position_order 
		FROM projects 
		WHERE curriculum_id = $1
		ORDER BY position_order
	`

	rows, err := s.db.Query(query, curriculumID)
	if err != nil {
		return fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	projectOrders := make(map[string]int)
	for rows.Next() {
		var identifier string
		var order int
		if err := rows.Scan(&identifier, &order); err != nil {
			return fmt.Errorf("failed to scan project: %w", err)
		}
		projectOrders[identifier] = order
	}

	// Get current project's order if it exists (for updates)
	currentOrder := -1
	if currentIdentifier != "" {
		if order, exists := projectOrders[currentIdentifier]; exists {
			currentOrder = order
		}
	}

	// Validate each prerequisite
	for _, prereq := range prerequisites {
		prereq = strings.TrimSpace(prereq)
		if prereq == "" {
			continue
		}

		order, exists := projectOrders[prereq]
		if !exists {
			return fmt.Errorf("prerequisite '%s' does not exist in this curriculum", prereq)
		}

		// For new projects, all prerequisites must have lower order
		// For existing projects, prerequisites must have order less than current
		if currentOrder != -1 && order >= currentOrder {
			return fmt.Errorf("prerequisite '%s' must come before this project", prereq)
		}
	}

	return nil
}

func (s *ProjectService) CreateProject(curriculumID int, req models.CreateProjectRequest) (*models.Project, error) {
	// Validate curriculum exists
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM curricula WHERE id = $1)", curriculumID).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check curriculum: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("curriculum not found")
	}

	// Generate identifier based on project type
	identifier, err := s.generateIdentifier(curriculumID, req.ProjectType)
	if err != nil {
		return nil, err
	}

	// Validate prerequisites
	if err := s.validatePrerequisites(curriculumID, req.Prerequisites, ""); err != nil {
		return nil, err
	}

	// Validate project type
	if !isValidProjectType(req.ProjectType) {
		return nil, fmt.Errorf("invalid project type: %s", req.ProjectType)
	}

	query := `
		INSERT INTO projects (curriculum_id, identifier, name, description, learning_objectives, estimated_time, prerequisites, project_type, position_order)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, curriculum_id, identifier, name, description, learning_objectives, estimated_time, prerequisites, project_type, position_order, created_at, updated_at
	`

	var project models.Project
	err = s.db.QueryRow(query, curriculumID, identifier, req.Name, req.Description,
		pq.Array(req.LearningObjectives), req.EstimatedTime, pq.Array(req.Prerequisites),
		req.ProjectType, req.PositionOrder).Scan(
		&project.ID, &project.CurriculumID, &project.Identifier, &project.Name,
		&project.Description, &project.LearningObjectives, &project.EstimatedTime,
		&project.Prerequisites, &project.ProjectType, &project.PositionOrder,
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
		ORDER BY p.position_order, p.created_at
	`

	rows, err := s.db.Query(query, userID, curriculumID)
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %w", err)
	}
	defer rows.Close()

	projects := make([]models.Project, 0)
	for rows.Next() {
		var p models.Project
		var progressID, progressUserID, progressProjectID, progressCompletionPercentage sql.NullInt64
		var progressStatus sql.NullString
		var progressStartedAt, progressCompletedAt, progressCreatedAt, progressUpdatedAt sql.NullTime

		err := rows.Scan(
			&p.ID, &p.CurriculumID, &p.Identifier, &p.Name, &p.Description,
			&p.LearningObjectives, &p.EstimatedTime, &p.Prerequisites,
			&p.ProjectType, &p.PositionOrder, &p.CreatedAt, &p.UpdatedAt,
			&progressID, &progressUserID, &progressProjectID, &progressStatus,
			&progressCompletionPercentage, &progressStartedAt, &progressCompletedAt,
			&progressCreatedAt, &progressUpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan project: %w", err)
		}

		if progressID.Valid {
			progress := models.Progress{
				ID:                   int(progressID.Int64),
				UserID:               int(progressUserID.Int64),
				ProjectID:            int(progressProjectID.Int64),
				Status:               progressStatus.String,
				CompletionPercentage: int(progressCompletionPercentage.Int64),
				StartedAt:            progressStartedAt,
				CompletedAt:          progressCompletedAt,
				CreatedAt:            progressCreatedAt.Time,
				UpdatedAt:            progressUpdatedAt.Time,
			}
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
		&project.Description, &project.LearningObjectives, &project.EstimatedTime,
		&project.Prerequisites, &project.ProjectType, &project.PositionOrder,
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
	// Get current project to validate prerequisites
	currentProject, err := s.GetProjectByID(userID, projectID)
	if err != nil {
		return nil, err
	}

	// Validate prerequisites
	if err := s.validatePrerequisites(currentProject.CurriculumID, req.Prerequisites, currentProject.Identifier); err != nil {
		return nil, err
	}

	// Validate project type
	if !isValidProjectType(req.ProjectType) {
		return nil, fmt.Errorf("invalid project type: %s", req.ProjectType)
	}

	query := `
		UPDATE projects
		SET name = $1, description = $2, learning_objectives = $3,
		    estimated_time = $4, prerequisites = $5, project_type = $6, position_order = $7,
		    updated_at = CURRENT_TIMESTAMP
		FROM curricula c
		WHERE projects.id = $8 AND projects.curriculum_id = c.id AND c.user_id = $9
		RETURNING projects.id, projects.curriculum_id, projects.identifier, projects.name, 
		         projects.description, projects.learning_objectives, projects.estimated_time, 
		         projects.prerequisites, projects.project_type, projects.position_order, 
		         projects.created_at, projects.updated_at
	`

	var project models.Project
	err = s.db.QueryRow(query, req.Name, req.Description,
		pq.Array(req.LearningObjectives), req.EstimatedTime, pq.Array(req.Prerequisites),
		req.ProjectType, req.PositionOrder, projectID, userID).Scan(
		&project.ID, &project.CurriculumID, &project.Identifier, &project.Name,
		&project.Description, &project.LearningObjectives, &project.EstimatedTime,
		&project.Prerequisites, &project.ProjectType, &project.PositionOrder,
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
	// Check if any other projects depend on this one
	var dependentCount int
	checkQuery := `
		SELECT COUNT(*)
		FROM projects p
		INNER JOIN curricula c ON p.curriculum_id = c.id
		WHERE c.user_id = $1 AND $2 = ANY(p.prerequisites)
	`

	err := s.db.QueryRow(checkQuery, userID, projectID).Scan(&dependentCount)
	if err != nil {
		return fmt.Errorf("failed to check dependencies: %w", err)
	}

	if dependentCount > 0 {
		return fmt.Errorf("cannot delete project: %d other projects depend on it", dependentCount)
	}

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

func isValidProjectType(projectType string) bool {
	validTypes := map[string]bool{
		models.ProjectTypeRoot:            true,
		models.ProjectTypeRootTest:        true,
		models.ProjectTypeBase:            true,
		models.ProjectTypeBaseTest:        true,
		models.ProjectTypeLowerBranch:     true,
		models.ProjectTypeMiddleBranch:    true,
		models.ProjectTypeUpperBranch:     true,
		models.ProjectTypeFlowerMilestone: true,
	}
	return validTypes[projectType]
}

package services

import (
	"curriculum-tracker/models"
	"database/sql"
	"fmt"
)

type CurriculumService struct {
	db *sql.DB
}

func NewCurriculumService(db *sql.DB) *CurriculumService {
	return &CurriculumService{db: db}
}

func (s *CurriculumService) CreateCurriculum(userID int, req models.CreateCurriculumRequest) (*models.Curriculum, error) {
	query := `
		INSERT INTO curricula (user_id, name, description)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, name, description, created_at, updated_at
	`

	var curriculum models.Curriculum
	err := s.db.QueryRow(query, userID, req.Name, req.Description).Scan(
		&curriculum.ID, &curriculum.UserID, &curriculum.Name, &curriculum.Description,
		&curriculum.CreatedAt, &curriculum.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create curriculum: %w", err)
	}

	return &curriculum, nil
}

func (s *CurriculumService) GetCurriculumsByUserID(userID int) ([]models.CurriculumWithStats, error) {
	query := `
		SELECT 
			c.id, c.user_id, c.name, c.description, c.created_at, c.updated_at,
			COUNT(p.id) as total_projects,
			COUNT(CASE WHEN pr.status = 'completed' THEN 1 END) as completed_projects,
			COALESCE(SUM(te.minutes), 0) as total_time_spent
		FROM curricula c
		LEFT JOIN projects p ON c.id = p.curriculum_id
		LEFT JOIN progress pr ON p.id = pr.project_id AND pr.user_id = $1
		LEFT JOIN time_entries te ON p.id = te.project_id AND te.user_id = $1
		WHERE c.user_id = $1
		GROUP BY c.id, c.user_id, c.name, c.description, c.created_at, c.updated_at
		ORDER BY c.created_at DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query curricula: %w", err)
	}
	defer rows.Close()

	// Initialize as empty slice, not nil
	curricula := make([]models.CurriculumWithStats, 0)
	for rows.Next() {
		var c models.CurriculumWithStats
		err := rows.Scan(
			&c.ID, &c.UserID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt,
			&c.TotalProjects, &c.CompletedProjects, &c.TotalTimeSpent,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan curriculum: %w", err)
		}
		curricula = append(curricula, c)
	}

	return curricula, nil
}

func (s *CurriculumService) GetCurriculumByID(userID, curriculumID int) (*models.Curriculum, error) {
	query := `
		SELECT id, user_id, name, description, created_at, updated_at
		FROM curricula
		WHERE id = $1 AND user_id = $2
	`

	var curriculum models.Curriculum
	err := s.db.QueryRow(query, curriculumID, userID).Scan(
		&curriculum.ID, &curriculum.UserID, &curriculum.Name, &curriculum.Description,
		&curriculum.CreatedAt, &curriculum.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("curriculum not found")
		}
		return nil, fmt.Errorf("failed to query curriculum: %w", err)
	}

	return &curriculum, nil
}

func (s *CurriculumService) UpdateCurriculum(userID, curriculumID int, req models.UpdateCurriculumRequest) (*models.Curriculum, error) {
	query := `
		UPDATE curricula
		SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND user_id = $4
		RETURNING id, user_id, name, description, created_at, updated_at
	`

	var curriculum models.Curriculum
	err := s.db.QueryRow(query, req.Name, req.Description, curriculumID, userID).Scan(
		&curriculum.ID, &curriculum.UserID, &curriculum.Name, &curriculum.Description,
		&curriculum.CreatedAt, &curriculum.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("curriculum not found")
		}
		return nil, fmt.Errorf("failed to update curriculum: %w", err)
	}

	return &curriculum, nil
}

func (s *CurriculumService) DeleteCurriculum(userID, curriculumID int) error {
	query := `DELETE FROM curricula WHERE id = $1 AND user_id = $2`

	result, err := s.db.Exec(query, curriculumID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete curriculum: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("curriculum not found")
	}

	return nil
}

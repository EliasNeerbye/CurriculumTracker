package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"curriculum-tracker/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type DB struct {
	Pool *pgxpool.Pool
}

func New(databaseURL string) (*DB, error) {
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &DB{Pool: pool}, nil
}

func (db *DB) Close() {
	db.Pool.Close()
}

func (db *DB) CreateUser(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users (id, email, password_hash, name, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)`

	now := time.Now()
	user.ID = uuid.New()
	user.CreatedAt = now
	user.UpdatedAt = now

	_, err := db.Pool.Exec(ctx, query, user.ID, user.Email, user.PasswordHash, user.Name, user.CreatedAt, user.UpdatedAt)
	return err
}

func (db *DB) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, name, created_at, updated_at FROM users WHERE email = $1`

	var user models.User
	err := db.Pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `SELECT id, email, password_hash, name, created_at, updated_at FROM users WHERE id = $1`

	var user models.User
	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.Name, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) CreateCurriculum(ctx context.Context, curriculum *models.Curriculum) error {
	query := `
        INSERT INTO curricula (id, user_id, name, description, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)`

	now := time.Now()
	curriculum.ID = uuid.New()
	curriculum.CreatedAt = now
	curriculum.UpdatedAt = now

	_, err := db.Pool.Exec(ctx, query, curriculum.ID, curriculum.UserID, curriculum.Name, curriculum.Description, curriculum.CreatedAt, curriculum.UpdatedAt)
	return err
}

func (db *DB) GetCurriculaByUserID(ctx context.Context, userID uuid.UUID) ([]models.Curriculum, error) {
	query := `SELECT id, user_id, name, description, created_at, updated_at FROM curricula WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var curricula []models.Curriculum
	for rows.Next() {
		var c models.Curriculum
		err := rows.Scan(&c.ID, &c.UserID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		curricula = append(curricula, c)
	}
	return curricula, rows.Err()
}

func (db *DB) GetCurriculumByID(ctx context.Context, id uuid.UUID) (*models.Curriculum, error) {
	query := `SELECT id, user_id, name, description, created_at, updated_at FROM curricula WHERE id = $1`

	var c models.Curriculum
	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&c.ID, &c.UserID, &c.Name, &c.Description, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (db *DB) UpdateCurriculum(ctx context.Context, curriculum *models.Curriculum) error {
	query := `UPDATE curricula SET name = $1, description = $2, updated_at = $3 WHERE id = $4`
	curriculum.UpdatedAt = time.Now()
	_, err := db.Pool.Exec(ctx, query, curriculum.Name, curriculum.Description, curriculum.UpdatedAt, curriculum.ID)
	return err
}

func (db *DB) DeleteCurriculum(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM curricula WHERE id = $1`
	_, err := db.Pool.Exec(ctx, query, id)
	return err
}

func (db *DB) CreateProject(ctx context.Context, project *models.Project) error {
	query := `
        INSERT INTO projects (id, curriculum_id, identifier, name, description, learning_objectives, estimated_time, prerequisites, project_type, order_index, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

	now := time.Now()
	project.ID = uuid.New()
	project.CreatedAt = now
	project.UpdatedAt = now

	_, err := db.Pool.Exec(ctx, query,
		project.ID, project.CurriculumID, project.Identifier, project.Name, project.Description,
		pq.Array(project.LearningObjectives), project.EstimatedTime, pq.Array(project.Prerequisites),
		project.ProjectType, project.OrderIndex, project.CreatedAt, project.UpdatedAt)
	return err
}

func (db *DB) GetProjectsByCurriculumID(ctx context.Context, curriculumID uuid.UUID) ([]models.Project, error) {
	query := `
        SELECT id, curriculum_id, identifier, name, description, learning_objectives, estimated_time, prerequisites, project_type, order_index, created_at, updated_at 
        FROM projects 
        WHERE curriculum_id = $1 
        ORDER BY project_type, order_index`

	rows, err := db.Pool.Query(ctx, query, curriculumID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var p models.Project
		err := rows.Scan(
			&p.ID, &p.CurriculumID, &p.Identifier, &p.Name, &p.Description,
			pq.Array(&p.LearningObjectives), &p.EstimatedTime, pq.Array(&p.Prerequisites),
			&p.ProjectType, &p.OrderIndex, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}

func (db *DB) GetProjectByID(ctx context.Context, id uuid.UUID) (*models.Project, error) {
	query := `
        SELECT id, curriculum_id, identifier, name, description, learning_objectives, estimated_time, prerequisites, project_type, order_index, created_at, updated_at 
        FROM projects 
        WHERE id = $1`

	var p models.Project
	err := db.Pool.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.CurriculumID, &p.Identifier, &p.Name, &p.Description,
		pq.Array(&p.LearningObjectives), &p.EstimatedTime, pq.Array(&p.Prerequisites),
		&p.ProjectType, &p.OrderIndex, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (db *DB) UpdateProject(ctx context.Context, project *models.Project) error {
	query := `
        UPDATE projects 
        SET identifier = $1, name = $2, description = $3, learning_objectives = $4, estimated_time = $5, prerequisites = $6, project_type = $7, order_index = $8, updated_at = $9 
        WHERE id = $10`

	project.UpdatedAt = time.Now()
	_, err := db.Pool.Exec(ctx, query,
		project.Identifier, project.Name, project.Description, pq.Array(project.LearningObjectives),
		project.EstimatedTime, pq.Array(project.Prerequisites), project.ProjectType, project.OrderIndex,
		project.UpdatedAt, project.ID)
	return err
}

func (db *DB) DeleteProject(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := db.Pool.Exec(ctx, query, id)
	return err
}

func (db *DB) GetOrCreateProjectProgress(ctx context.Context, userID, projectID uuid.UUID) (*models.ProjectProgress, error) {
	query := `
        SELECT id, user_id, project_id, status, time_spent_minutes, started_at, completed_at, created_at, updated_at 
        FROM project_progress 
        WHERE user_id = $1 AND project_id = $2`

	var progress models.ProjectProgress
	err := db.Pool.QueryRow(ctx, query, userID, projectID).Scan(
		&progress.ID, &progress.UserID, &progress.ProjectID, &progress.Status,
		&progress.TimeSpentMinutes, &progress.StartedAt, &progress.CompletedAt,
		&progress.CreatedAt, &progress.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		progress = models.ProjectProgress{
			ID:               uuid.New(),
			UserID:           userID,
			ProjectID:        projectID,
			Status:           models.StatusNotStarted,
			TimeSpentMinutes: 0,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}

		insertQuery := `
            INSERT INTO project_progress (id, user_id, project_id, status, time_spent_minutes, created_at, updated_at)
            VALUES ($1, $2, $3, $4, $5, $6, $7)`

		_, err = db.Pool.Exec(ctx, insertQuery, progress.ID, progress.UserID, progress.ProjectID,
			progress.Status, progress.TimeSpentMinutes, progress.CreatedAt, progress.UpdatedAt)
		if err != nil {
			return nil, err
		}
		return &progress, nil
	} else if err != nil {
		return nil, err
	}

	return &progress, nil
}

func (db *DB) UpdateProjectProgress(ctx context.Context, progress *models.ProjectProgress) error {
	query := `
        UPDATE project_progress 
        SET status = $1, time_spent_minutes = $2, started_at = $3, completed_at = $4, updated_at = $5 
        WHERE id = $6`

	progress.UpdatedAt = time.Now()

	if progress.Status == models.StatusInProgress && progress.StartedAt == nil {
		now := time.Now()
		progress.StartedAt = &now
	}

	if progress.Status == models.StatusCompleted && progress.CompletedAt == nil {
		now := time.Now()
		progress.CompletedAt = &now
	}

	_, err := db.Pool.Exec(ctx, query, progress.Status, progress.TimeSpentMinutes,
		progress.StartedAt, progress.CompletedAt, progress.UpdatedAt, progress.ID)
	return err
}

func (db *DB) CreateNote(ctx context.Context, note *models.Note) error {
	query := `
        INSERT INTO notes (id, user_id, project_id, title, content, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`

	now := time.Now()
	note.ID = uuid.New()
	note.CreatedAt = now
	note.UpdatedAt = now

	_, err := db.Pool.Exec(ctx, query, note.ID, note.UserID, note.ProjectID, note.Title, note.Content, note.CreatedAt, note.UpdatedAt)
	return err
}

func (db *DB) GetNotesByProjectID(ctx context.Context, userID, projectID uuid.UUID) ([]models.Note, error) {
	query := `
        SELECT id, user_id, project_id, title, content, created_at, updated_at 
        FROM notes 
        WHERE user_id = $1 AND project_id = $2 
        ORDER BY created_at DESC`

	rows, err := db.Pool.Query(ctx, query, userID, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note
	for rows.Next() {
		var n models.Note
		err := rows.Scan(&n.ID, &n.UserID, &n.ProjectID, &n.Title, &n.Content, &n.CreatedAt, &n.UpdatedAt)
		if err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	return notes, rows.Err()
}

func (db *DB) UpdateNote(ctx context.Context, note *models.Note) error {
	query := `UPDATE notes SET title = $1, content = $2, updated_at = $3 WHERE id = $4 AND user_id = $5`
	note.UpdatedAt = time.Now()
	_, err := db.Pool.Exec(ctx, query, note.Title, note.Content, note.UpdatedAt, note.ID, note.UserID)
	return err
}

func (db *DB) DeleteNote(ctx context.Context, id, userID uuid.UUID) error {
	query := `DELETE FROM notes WHERE id = $1 AND user_id = $2`
	_, err := db.Pool.Exec(ctx, query, id, userID)
	return err
}

func (db *DB) CreateReflection(ctx context.Context, reflection *models.Reflection) error {
	query := `
        INSERT INTO reflections (id, user_id, project_id, content, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)`

	now := time.Now()
	reflection.ID = uuid.New()
	reflection.CreatedAt = now
	reflection.UpdatedAt = now

	_, err := db.Pool.Exec(ctx, query, reflection.ID, reflection.UserID, reflection.ProjectID, reflection.Content, reflection.CreatedAt, reflection.UpdatedAt)
	return err
}

func (db *DB) GetReflectionsByProjectID(ctx context.Context, userID, projectID uuid.UUID) ([]models.Reflection, error) {
	query := `
        SELECT id, user_id, project_id, content, created_at, updated_at 
        FROM reflections 
        WHERE user_id = $1 AND project_id = $2 
        ORDER BY created_at DESC`

	rows, err := db.Pool.Query(ctx, query, userID, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reflections []models.Reflection
	for rows.Next() {
		var r models.Reflection
		err := rows.Scan(&r.ID, &r.UserID, &r.ProjectID, &r.Content, &r.CreatedAt, &r.UpdatedAt)
		if err != nil {
			return nil, err
		}
		reflections = append(reflections, r)
	}
	return reflections, rows.Err()
}

func (db *DB) UpdateReflection(ctx context.Context, reflection *models.Reflection) error {
	query := `UPDATE reflections SET content = $1, updated_at = $2 WHERE id = $3 AND user_id = $4`
	reflection.UpdatedAt = time.Now()
	_, err := db.Pool.Exec(ctx, query, reflection.Content, reflection.UpdatedAt, reflection.ID, reflection.UserID)
	return err
}

func (db *DB) DeleteReflection(ctx context.Context, id, userID uuid.UUID) error {
	query := `DELETE FROM reflections WHERE id = $1 AND user_id = $2`
	_, err := db.Pool.Exec(ctx, query, id, userID)
	return err
}

func (db *DB) CreateTimeEntry(ctx context.Context, entry *models.TimeEntry) error {
	query := `
        INSERT INTO time_entries (id, user_id, project_id, minutes, description, logged_at, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`

	now := time.Now()
	entry.ID = uuid.New()
	entry.CreatedAt = now
	if entry.LoggedAt.IsZero() {
		entry.LoggedAt = now
	}

	_, err := db.Pool.Exec(ctx, query, entry.ID, entry.UserID, entry.ProjectID, entry.Minutes, entry.Description, entry.LoggedAt, entry.CreatedAt)

	updateQuery := `
        UPDATE project_progress 
        SET time_spent_minutes = time_spent_minutes + $1, updated_at = $2 
        WHERE user_id = $3 AND project_id = $4`
	_, err2 := db.Pool.Exec(ctx, updateQuery, entry.Minutes, now, entry.UserID, entry.ProjectID)
	if err2 != nil {
		return err2
	}

	return err
}

func (db *DB) GetTimeEntriesByProjectID(ctx context.Context, userID, projectID uuid.UUID) ([]models.TimeEntry, error) {
	query := `
        SELECT id, user_id, project_id, minutes, description, logged_at, created_at 
        FROM time_entries 
        WHERE user_id = $1 AND project_id = $2 
        ORDER BY logged_at DESC`

	rows, err := db.Pool.Query(ctx, query, userID, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []models.TimeEntry
	for rows.Next() {
		var e models.TimeEntry
		err := rows.Scan(&e.ID, &e.UserID, &e.ProjectID, &e.Minutes, &e.Description, &e.LoggedAt, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}
	return entries, rows.Err()
}

func (db *DB) GetAnalytics(ctx context.Context, userID uuid.UUID) (*models.AnalyticsResponse, error) {
	totalProjectsQuery := `
        SELECT COUNT(*) 
        FROM projects p 
        JOIN curricula c ON p.curriculum_id = c.id 
        WHERE c.user_id = $1`

	var totalProjects int
	err := db.Pool.QueryRow(ctx, totalProjectsQuery, userID).Scan(&totalProjects)
	if err != nil {
		return nil, err
	}

	progressQuery := `
        SELECT 
            COUNT(CASE WHEN status = 'completed' THEN 1 END) as completed,
            COUNT(CASE WHEN status = 'in_progress' THEN 1 END) as in_progress,
            COALESCE(SUM(time_spent_minutes), 0) as total_time
        FROM project_progress pp
        WHERE pp.user_id = $1`

	var completed, inProgress, totalTime int
	err = db.Pool.QueryRow(ctx, progressQuery, userID).Scan(&completed, &inProgress, &totalTime)
	if err != nil {
		return nil, err
	}

	projectsByTypeQuery := `
        SELECT p.project_type, COUNT(*) 
        FROM projects p 
        JOIN curricula c ON p.curriculum_id = c.id 
        WHERE c.user_id = $1 
        GROUP BY p.project_type`

	rows, err := db.Pool.Query(ctx, projectsByTypeQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projectsByType := make(map[models.ProjectType]int)
	for rows.Next() {
		var projectType models.ProjectType
		var count int
		err := rows.Scan(&projectType, &count)
		if err != nil {
			return nil, err
		}
		projectsByType[projectType] = count
	}

	completionByTypeQuery := `
        SELECT p.project_type, 
               COUNT(*) as total,
               COUNT(CASE WHEN pp.status = 'completed' THEN 1 END) as completed
        FROM projects p 
        JOIN curricula c ON p.curriculum_id = c.id 
        LEFT JOIN project_progress pp ON p.id = pp.project_id AND pp.user_id = $1
        WHERE c.user_id = $1 
        GROUP BY p.project_type`

	rows2, err := db.Pool.Query(ctx, completionByTypeQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows2.Close()

	completionByType := make(map[models.ProjectType]float64)
	for rows2.Next() {
		var projectType models.ProjectType
		var total, completedCount int
		err := rows2.Scan(&projectType, &total, &completedCount)
		if err != nil {
			return nil, err
		}
		if total > 0 {
			completionByType[projectType] = float64(completedCount) / float64(total) * 100
		}
	}

	recentActivityQuery := `
        SELECT id, user_id, project_id, minutes, description, logged_at, created_at 
        FROM time_entries 
        WHERE user_id = $1 
        ORDER BY logged_at DESC 
        LIMIT 10`

	rows3, err := db.Pool.Query(ctx, recentActivityQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows3.Close()

	var recentActivity []models.TimeEntry
	for rows3.Next() {
		var e models.TimeEntry
		err := rows3.Scan(&e.ID, &e.UserID, &e.ProjectID, &e.Minutes, &e.Description, &e.LoggedAt, &e.CreatedAt)
		if err != nil {
			return nil, err
		}
		recentActivity = append(recentActivity, e)
	}

	weeklyTimeQuery := `
        SELECT 
            DATE_TRUNC('week', logged_at) as week,
            SUM(minutes) as total_minutes
        FROM time_entries 
        WHERE user_id = $1 AND logged_at >= NOW() - INTERVAL '8 weeks'
        GROUP BY week 
        ORDER BY week`

	rows4, err := db.Pool.Query(ctx, weeklyTimeQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows4.Close()

	weeklyTimeSpent := make([]int, 8)
	for rows4.Next() {
		var week time.Time
		var minutes int
		err := rows4.Scan(&week, &minutes)
		if err != nil {
			return nil, err
		}

		weeksAgo := int(time.Since(week).Hours() / 24 / 7)
		if weeksAgo >= 0 && weeksAgo < 8 {
			weeklyTimeSpent[7-weeksAgo] = minutes
		}
	}

	return &models.AnalyticsResponse{
		TotalProjects:      totalProjects,
		CompletedProjects:  completed,
		InProgressProjects: inProgress,
		TotalTimeSpent:     totalTime,
		ProjectsByType:     projectsByType,
		CompletionByType:   completionByType,
		RecentActivity:     recentActivity,
		WeeklyTimeSpent:    weeklyTimeSpent,
	}, nil
}

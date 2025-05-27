package services

import (
	"curriculum-tracker/models"
	"database/sql"
	"fmt"
	"time"
)

type AnalyticsService struct {
	db *sql.DB
}

func NewAnalyticsService(db *sql.DB) *AnalyticsService {
	return &AnalyticsService{db: db}
}

func (s *AnalyticsService) CreateTimeEntry(userID int, req models.CreateTimeEntryRequest) (*models.TimeEntry, error) {
	parsedDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, fmt.Errorf("invalid date format: %w", err)
	}

	query := `
		INSERT INTO time_entries (user_id, project_id, minutes, description, date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, project_id, minutes, description, date, created_at
	`

	var timeEntry models.TimeEntry
	err = s.db.QueryRow(query, userID, req.ProjectID, req.Minutes, req.Description, parsedDate).Scan(
		&timeEntry.ID, &timeEntry.UserID, &timeEntry.ProjectID, &timeEntry.Minutes,
		&timeEntry.Description, &timeEntry.Date, &timeEntry.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create time entry: %w", err)
	}

	return &timeEntry, nil
}

func (s *AnalyticsService) GetTimeEntriesByProjectID(userID, projectID int) ([]models.TimeEntry, error) {
	query := `
		SELECT te.id, te.user_id, te.project_id, te.minutes, te.description, te.date, te.created_at
		FROM time_entries te
		JOIN projects p ON te.project_id = p.id
		JOIN curricula c ON p.curriculum_id = c.id
		WHERE te.user_id = $1 AND te.project_id = $2 AND c.user_id = $1
		ORDER BY te.date DESC, te.created_at DESC
	`

	rows, err := s.db.Query(query, userID, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query time entries: %w", err)
	}
	defer rows.Close()

	var timeEntries []models.TimeEntry
	for rows.Next() {
		var te models.TimeEntry
		err := rows.Scan(
			&te.ID, &te.UserID, &te.ProjectID, &te.Minutes,
			&te.Description, &te.Date, &te.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan time entry: %w", err)
		}
		timeEntries = append(timeEntries, te)
	}

	return timeEntries, nil
}

func (s *AnalyticsService) GetTimeStatsByCurriculumID(userID, curriculumID int) (*models.TimeStats, error) {
	query := `
		SELECT 
			COALESCE(SUM(te.minutes), 0) as total_minutes,
			te.date,
			COALESCE(SUM(te.minutes), 0) as daily_minutes,
			p.name as project_name,
			COALESCE(SUM(te.minutes), 0) as project_minutes
		FROM time_entries te
		JOIN projects p ON te.project_id = p.id
		JOIN curricula c ON p.curriculum_id = c.id
		WHERE te.user_id = $1 AND c.id = $2 AND c.user_id = $1
		GROUP BY te.date, p.name
		ORDER BY te.date DESC
	`

	rows, err := s.db.Query(query, userID, curriculumID)
	if err != nil {
		return nil, fmt.Errorf("failed to query time stats: %w", err)
	}
	defer rows.Close()

	stats := &models.TimeStats{
		DailyBreakdown:   make(map[string]int),
		ProjectBreakdown: make(map[string]int),
	}

	for rows.Next() {
		var totalMinutes, dailyMinutes, projectMinutes int
		var date time.Time
		var projectName string

		err := rows.Scan(&totalMinutes, &date, &dailyMinutes, &projectName, &projectMinutes)
		if err != nil {
			return nil, fmt.Errorf("failed to scan time stats: %w", err)
		}

		stats.TotalMinutes += totalMinutes
		dateStr := date.Format("2006-01-02")
		stats.DailyBreakdown[dateStr] += dailyMinutes
		stats.ProjectBreakdown[projectName] += projectMinutes
	}

	if len(stats.DailyBreakdown) > 0 {
		weekCount := 0
		weeklyTotal := 0
		now := time.Now()
		for i := 0; i < 7; i++ {
			date := now.AddDate(0, 0, -i)
			dateStr := date.Format("2006-01-02")
			if minutes, exists := stats.DailyBreakdown[dateStr]; exists {
				weeklyTotal += minutes
				weekCount++
			}
		}
		if weekCount > 0 {
			stats.WeeklyAverage = float64(weeklyTotal) / float64(weekCount)
		}
	}

	return stats, nil
}

func (s *AnalyticsService) GetUserOverallStats(userID int) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(DISTINCT c.id) as total_curricula,
			COUNT(DISTINCT p.id) as total_projects,
			COUNT(DISTINCT CASE WHEN pr.status = 'completed' THEN p.id END) as completed_projects,
			COUNT(DISTINCT CASE WHEN pr.status = 'in_progress' THEN p.id END) as in_progress_projects,
			COALESCE(SUM(te.minutes), 0) as total_time_minutes,
			COUNT(DISTINCT n.id) as total_notes
		FROM curricula c
		LEFT JOIN projects p ON c.id = p.curriculum_id
		LEFT JOIN progress pr ON p.id = pr.project_id AND pr.user_id = $1
		LEFT JOIN time_entries te ON p.id = te.project_id AND te.user_id = $1
		LEFT JOIN notes n ON p.id = n.project_id AND n.user_id = $1
		WHERE c.user_id = $1
	`

	var totalCurricula, totalProjects, completedProjects, inProgressProjects, totalTimeMinutes, totalNotes int

	err := s.db.QueryRow(query, userID).Scan(
		&totalCurricula, &totalProjects, &completedProjects,
		&inProgressProjects, &totalTimeMinutes, &totalNotes,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query user stats: %w", err)
	}

	completionRate := 0.0
	if totalProjects > 0 {
		completionRate = float64(completedProjects) / float64(totalProjects) * 100
	}

	return map[string]interface{}{
		"total_curricula":      totalCurricula,
		"total_projects":       totalProjects,
		"completed_projects":   completedProjects,
		"in_progress_projects": inProgressProjects,
		"total_time_minutes":   totalTimeMinutes,
		"total_time_hours":     float64(totalTimeMinutes) / 60,
		"total_notes":          totalNotes,
		"completion_rate":      completionRate,
	}, nil
}

package models

import (
	"time"
)

type TimeEntry struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	ProjectID   int       `json:"project_id"`
	Minutes     int       `json:"minutes"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateTimeEntryRequest struct {
	ProjectID   int    `json:"project_id"`
	Minutes     int    `json:"minutes"`
	Description string `json:"description"`
	Date        string `json:"date"`
}

type TimeStats struct {
	TotalMinutes     int            `json:"total_minutes"`
	DailyBreakdown   map[string]int `json:"daily_breakdown"`
	ProjectBreakdown map[string]int `json:"project_breakdown"`
	WeeklyAverage    float64        `json:"weekly_average"`
}

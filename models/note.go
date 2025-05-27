package models

import (
	"time"
)

type Note struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ProjectID int       `json:"project_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	NoteType  string    `json:"note_type"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateNoteRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	NoteType string `json:"note_type"`
}

type UpdateNoteRequest struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	NoteType string `json:"note_type"`
}

const (
	NoteTypeNote       = "note"
	NoteTypeReflection = "reflection"
	NoteTypeLearning   = "learning"
	NoteTypeQuestion   = "question"
)

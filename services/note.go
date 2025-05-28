package services

import (
	"curriculum-tracker/models"
	"database/sql"
	"fmt"
)

type NoteService struct {
	db *sql.DB
}

func NewNoteService(db *sql.DB) *NoteService {
	return &NoteService{db: db}
}

func (s *NoteService) CreateNote(userID, projectID int, req models.CreateNoteRequest) (*models.Note, error) {
	query := `
		INSERT INTO notes (user_id, project_id, title, content, note_type)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, project_id, title, content, note_type, created_at, updated_at
	`

	var note models.Note
	err := s.db.QueryRow(query, userID, projectID, req.Title, req.Content, req.NoteType).Scan(
		&note.ID, &note.UserID, &note.ProjectID, &note.Title, &note.Content,
		&note.NoteType, &note.CreatedAt, &note.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create note: %w", err)
	}

	return &note, nil
}

func (s *NoteService) GetNotesByProjectID(userID, projectID int) ([]models.Note, error) {
	query := `
		SELECT n.id, n.user_id, n.project_id, n.title, n.content, n.note_type, n.created_at, n.updated_at
		FROM notes n
		JOIN projects p ON n.project_id = p.id
		JOIN curricula c ON p.curriculum_id = c.id
		WHERE n.user_id = $1 AND n.project_id = $2 AND c.user_id = $1
		ORDER BY n.created_at DESC
	`

	rows, err := s.db.Query(query, userID, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to query notes: %w", err)
	}
	defer rows.Close()

	// Initialize as empty slice, not nil
	notes := make([]models.Note, 0)
	for rows.Next() {
		var note models.Note
		err := rows.Scan(
			&note.ID, &note.UserID, &note.ProjectID, &note.Title, &note.Content,
			&note.NoteType, &note.CreatedAt, &note.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan note: %w", err)
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (s *NoteService) GetNoteByID(userID, noteID int) (*models.Note, error) {
	query := `
		SELECT n.id, n.user_id, n.project_id, n.title, n.content, n.note_type, n.created_at, n.updated_at
		FROM notes n
		JOIN projects p ON n.project_id = p.id
		JOIN curricula c ON p.curriculum_id = c.id
		WHERE n.id = $1 AND n.user_id = $2 AND c.user_id = $2
	`

	var note models.Note
	err := s.db.QueryRow(query, noteID, userID).Scan(
		&note.ID, &note.UserID, &note.ProjectID, &note.Title, &note.Content,
		&note.NoteType, &note.CreatedAt, &note.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("note not found")
		}
		return nil, fmt.Errorf("failed to query note: %w", err)
	}

	return &note, nil
}

func (s *NoteService) UpdateNote(userID, noteID int, req models.UpdateNoteRequest) (*models.Note, error) {
	query := `
		UPDATE notes
		SET title = $1, content = $2, note_type = $3, updated_at = CURRENT_TIMESTAMP
		FROM projects p, curricula c
		WHERE notes.id = $4 AND notes.user_id = $5 AND notes.project_id = p.id 
		      AND p.curriculum_id = c.id AND c.user_id = $5
		RETURNING notes.id, notes.user_id, notes.project_id, notes.title, notes.content, 
		         notes.note_type, notes.created_at, notes.updated_at
	`

	var note models.Note
	err := s.db.QueryRow(query, req.Title, req.Content, req.NoteType, noteID, userID).Scan(
		&note.ID, &note.UserID, &note.ProjectID, &note.Title, &note.Content,
		&note.NoteType, &note.CreatedAt, &note.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("note not found")
		}
		return nil, fmt.Errorf("failed to update note: %w", err)
	}

	return &note, nil
}

func (s *NoteService) DeleteNote(userID, noteID int) error {
	query := `
		DELETE FROM notes
		USING projects p, curricula c
		WHERE notes.id = $1 AND notes.user_id = $2 AND notes.project_id = p.id 
		      AND p.curriculum_id = c.id AND c.user_id = $2
	`

	result, err := s.db.Exec(query, noteID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("note not found")
	}

	return nil
}

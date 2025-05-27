package database

import (
	"database/sql"
	"fmt"
)

func RunMigrations(db *sql.DB) error {
	migrations := []string{
		createUsersTable,
		createCurriculaTable,
		createProjectsTable,
		createProgressTable,
		createNotesTable,
		createTimeEntriesTable,
		createIndexes,
	}

	for i, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration %d failed: %w", i+1, err)
		}
	}

	return nil
}

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	email VARCHAR(255) UNIQUE NOT NULL,
	password_hash VARCHAR(255) NOT NULL,
	name VARCHAR(255) NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

const createCurriculaTable = `
CREATE TABLE IF NOT EXISTS curricula (
	id SERIAL PRIMARY KEY,
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	name VARCHAR(255) NOT NULL,
	description TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

const createProjectsTable = `
CREATE TABLE IF NOT EXISTS projects (
	id SERIAL PRIMARY KEY,
	curriculum_id INTEGER NOT NULL REFERENCES curricula(id) ON DELETE CASCADE,
	identifier VARCHAR(50),
	name VARCHAR(255) NOT NULL,
	description TEXT,
	learning_objectives TEXT[],
	estimated_time VARCHAR(100),
	prerequisites TEXT[],
	project_type VARCHAR(50) NOT NULL,
	position_order INTEGER NOT NULL DEFAULT 0,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

const createProgressTable = `
CREATE TABLE IF NOT EXISTS progress (
	id SERIAL PRIMARY KEY,
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
	status VARCHAR(50) NOT NULL DEFAULT 'not_started',
	completion_percentage INTEGER DEFAULT 0,
	started_at TIMESTAMP,
	completed_at TIMESTAMP,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(user_id, project_id)
);
`

const createNotesTable = `
CREATE TABLE IF NOT EXISTS notes (
	id SERIAL PRIMARY KEY,
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
	title VARCHAR(255),
	content TEXT NOT NULL,
	note_type VARCHAR(50) DEFAULT 'note',
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

const createTimeEntriesTable = `
CREATE TABLE IF NOT EXISTS time_entries (
	id SERIAL PRIMARY KEY,
	user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	project_id INTEGER NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
	minutes INTEGER NOT NULL,
	description TEXT,
	date DATE NOT NULL DEFAULT CURRENT_DATE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

const createIndexes = `
CREATE INDEX IF NOT EXISTS idx_curricula_user_id ON curricula(user_id);
CREATE INDEX IF NOT EXISTS idx_projects_curriculum_id ON projects(curriculum_id);
CREATE INDEX IF NOT EXISTS idx_progress_user_id ON progress(user_id);
CREATE INDEX IF NOT EXISTS idx_progress_project_id ON progress(project_id);
CREATE INDEX IF NOT EXISTS idx_notes_user_id ON notes(user_id);
CREATE INDEX IF NOT EXISTS idx_notes_project_id ON notes(project_id);
CREATE INDEX IF NOT EXISTS idx_time_entries_user_id ON time_entries(user_id);
CREATE INDEX IF NOT EXISTS idx_time_entries_project_id ON time_entries(project_id);
CREATE INDEX IF NOT EXISTS idx_time_entries_date ON time_entries(date);
`

-- Performance indexes for CurriculumTracker
-- Run this script to add indexes that will improve query performance and reduce timeouts

-- Index for curricula queries by user_id
CREATE INDEX IF NOT EXISTS idx_curricula_user_id ON curricula(user_id);

-- Index for projects queries by curriculum_id
CREATE INDEX IF NOT EXISTS idx_projects_curriculum_id ON projects(curriculum_id);

-- Index for project_progress queries by user_id and project_id
CREATE INDEX IF NOT EXISTS idx_project_progress_user_project ON project_progress(user_id, project_id);

-- Index for project_progress status queries
CREATE INDEX IF NOT EXISTS idx_project_progress_status ON project_progress(status);

-- Index for time_entries queries by user_id
CREATE INDEX IF NOT EXISTS idx_time_entries_user_id ON time_entries(user_id);

-- Index for time_entries queries by logged_at (for analytics)
CREATE INDEX IF NOT EXISTS idx_time_entries_logged_at ON time_entries(logged_at);

-- Index for time_entries queries by project_id
CREATE INDEX IF NOT EXISTS idx_time_entries_project_id ON time_entries(project_id);

-- Index for notes queries by user_id and project_id
CREATE INDEX IF NOT EXISTS idx_notes_user_project ON notes(user_id, project_id);

-- Index for reflections queries by user_id and project_id
CREATE INDEX IF NOT EXISTS idx_reflections_user_project ON reflections(user_id, project_id);

-- Composite index for the analytics query (most important for performance)
CREATE INDEX IF NOT EXISTS idx_curricula_projects_analytics ON projects(curriculum_id, project_type);

-- Index for projects order
CREATE INDEX IF NOT EXISTS idx_projects_order ON projects(curriculum_id, order_index);

-- Analyze tables to update statistics
ANALYZE users;
ANALYZE curricula;
ANALYZE projects;
ANALYZE project_progress;
ANALYZE time_entries;
ANALYZE notes;
ANALYZE reflections;

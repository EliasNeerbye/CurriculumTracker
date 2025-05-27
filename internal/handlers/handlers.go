package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"curriculum-tracker/internal/auth"
	"curriculum-tracker/internal/database"
	"curriculum-tracker/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type Handler struct {
	db   *database.DB
	auth *auth.Config
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func New(db *database.DB, authConfig *auth.Config) *Handler {
	return &Handler{
		db:   db,
		auth: authConfig,
	}
}

func (h *Handler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, ErrorResponse{Error: message})
}

func (h *Handler) getUserFromContext(ctx context.Context) (*models.User, error) {
	user, ok := ctx.Value("user").(*models.User)
	if !ok {
		return nil, fmt.Errorf("user not found in context")
	}
	return user, nil
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			h.writeError(w, http.StatusUnauthorized, "Missing authorization header")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			h.writeError(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		claims, err := h.auth.ValidateToken(tokenString)
		if err != nil {
			h.writeError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		user, err := h.db.GetUserByID(r.Context(), claims.UserID)
		if err != nil {
			if err == sql.ErrNoRows {
				h.writeError(w, http.StatusUnauthorized, "User not found")
				return
			}
			h.writeError(w, http.StatusInternalServerError, "Failed to get user")
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Email == "" || req.Password == "" || req.Name == "" {
		h.writeError(w, http.StatusBadRequest, "Email, password, and name are required")
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Name:         req.Name,
	}

	if err := h.db.CreateUser(r.Context(), user); err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			h.writeError(w, http.StatusConflict, "Email already exists")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	token, err := h.auth.GenerateToken(user)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	h.writeJSON(w, http.StatusCreated, models.AuthResponse{
		Token: token,
		User:  *user,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			h.writeError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	valid, err := auth.VerifyPassword(req.Password, user.PasswordHash)
	if err != nil || !valid {
		h.writeError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := h.auth.GenerateToken(user)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	h.writeJSON(w, http.StatusOK, models.AuthResponse{
		Token: token,
		User:  *user,
	})
}

func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	h.writeJSON(w, http.StatusOK, user)
}

func (h *Handler) CreateCurriculum(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	var req models.CurriculumRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		h.writeError(w, http.StatusBadRequest, "Name is required")
		return
	}

	curriculum := &models.Curriculum{
		UserID:      user.ID,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := h.db.CreateCurriculum(r.Context(), curriculum); err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to create curriculum")
		return
	}

	h.writeJSON(w, http.StatusCreated, curriculum)
}

func (h *Handler) GetCurricula(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	curricula, err := h.db.GetCurriculaByUserID(r.Context(), user.ID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get curricula")
		return
	}

	h.writeJSON(w, http.StatusOK, curricula)
}

func (h *Handler) GetCurriculum(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	curriculumID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid curriculum ID")
		return
	}

	curriculum, err := h.db.GetCurriculumByID(r.Context(), curriculumID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "Curriculum not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to get curriculum")
		return
	}

	if curriculum.UserID != user.ID {
		h.writeError(w, http.StatusForbidden, "Access denied")
		return
	}

	h.writeJSON(w, http.StatusOK, curriculum)
}

func (h *Handler) UpdateCurriculum(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	curriculumID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid curriculum ID")
		return
	}

	curriculum, err := h.db.GetCurriculumByID(r.Context(), curriculumID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "Curriculum not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to get curriculum")
		return
	}

	if curriculum.UserID != user.ID {
		h.writeError(w, http.StatusForbidden, "Access denied")
		return
	}

	var req models.CurriculumRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		h.writeError(w, http.StatusBadRequest, "Name is required")
		return
	}

	curriculum.Name = req.Name
	curriculum.Description = req.Description

	if err := h.db.UpdateCurriculum(r.Context(), curriculum); err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to update curriculum")
		return
	}

	h.writeJSON(w, http.StatusOK, curriculum)
}

func (h *Handler) DeleteCurriculum(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	curriculumID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid curriculum ID")
		return
	}

	curriculum, err := h.db.GetCurriculumByID(r.Context(), curriculumID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "Curriculum not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to get curriculum")
		return
	}

	if curriculum.UserID != user.ID {
		h.writeError(w, http.StatusForbidden, "Access denied")
		return
	}

	if err := h.db.DeleteCurriculum(r.Context(), curriculumID); err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to delete curriculum")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateProject(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	curriculumID, err := uuid.Parse(chi.URLParam(r, "curriculumId"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid curriculum ID")
		return
	}

	curriculum, err := h.db.GetCurriculumByID(r.Context(), curriculumID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "Curriculum not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to get curriculum")
		return
	}

	if curriculum.UserID != user.ID {
		h.writeError(w, http.StatusForbidden, "Access denied")
		return
	}

	var req models.ProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		h.writeError(w, http.StatusBadRequest, "Name is required")
		return
	}

	project := &models.Project{
		CurriculumID:       curriculumID,
		Identifier:         req.Identifier,
		Name:               req.Name,
		Description:        req.Description,
		LearningObjectives: req.LearningObjectives,
		EstimatedTime:      req.EstimatedTime,
		Prerequisites:      req.Prerequisites,
		ProjectType:        req.ProjectType,
		OrderIndex:         req.OrderIndex,
	}

	if err := h.db.CreateProject(r.Context(), project); err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to create project")
		return
	}

	h.writeJSON(w, http.StatusCreated, project)
}

func (h *Handler) GetProjects(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	curriculumID, err := uuid.Parse(chi.URLParam(r, "curriculumId"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid curriculum ID")
		return
	}

	curriculum, err := h.db.GetCurriculumByID(r.Context(), curriculumID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "Curriculum not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to get curriculum")
		return
	}

	if curriculum.UserID != user.ID {
		h.writeError(w, http.StatusForbidden, "Access denied")
		return
	}

	projects, err := h.db.GetProjectsByCurriculumID(r.Context(), curriculumID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get projects")
		return
	}

	h.writeJSON(w, http.StatusOK, projects)
}

func (h *Handler) GetProject(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	project, err := h.db.GetProjectByID(r.Context(), projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "Project not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to get project")
		return
	}

	curriculum, err := h.db.GetCurriculumByID(r.Context(), project.CurriculumID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get curriculum")
		return
	}

	if curriculum.UserID != user.ID {
		h.writeError(w, http.StatusForbidden, "Access denied")
		return
	}

	h.writeJSON(w, http.StatusOK, project)
}

func (h *Handler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	project, err := h.db.GetProjectByID(r.Context(), projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "Project not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to get project")
		return
	}

	curriculum, err := h.db.GetCurriculumByID(r.Context(), project.CurriculumID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get curriculum")
		return
	}

	if curriculum.UserID != user.ID {
		h.writeError(w, http.StatusForbidden, "Access denied")
		return
	}

	var req models.ProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		h.writeError(w, http.StatusBadRequest, "Name is required")
		return
	}

	project.Identifier = req.Identifier
	project.Name = req.Name
	project.Description = req.Description
	project.LearningObjectives = req.LearningObjectives
	project.EstimatedTime = req.EstimatedTime
	project.Prerequisites = req.Prerequisites
	project.ProjectType = req.ProjectType
	project.OrderIndex = req.OrderIndex

	if err := h.db.UpdateProject(r.Context(), project); err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to update project")
		return
	}

	h.writeJSON(w, http.StatusOK, project)
}

func (h *Handler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	project, err := h.db.GetProjectByID(r.Context(), projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "Project not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to get project")
		return
	}

	curriculum, err := h.db.GetCurriculumByID(r.Context(), project.CurriculumID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get curriculum")
		return
	}

	if curriculum.UserID != user.ID {
		h.writeError(w, http.StatusForbidden, "Access denied")
		return
	}

	if err := h.db.DeleteProject(r.Context(), projectID); err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to delete project")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetProjectProgress(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	project, err := h.db.GetProjectByID(r.Context(), projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "Project not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to get project")
		return
	}

	curriculum, err := h.db.GetCurriculumByID(r.Context(), project.CurriculumID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get curriculum")
		return
	}

	if curriculum.UserID != user.ID {
		h.writeError(w, http.StatusForbidden, "Access denied")
		return
	}

	progress, err := h.db.GetOrCreateProjectProgress(r.Context(), user.ID, projectID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get progress")
		return
	}

	h.writeJSON(w, http.StatusOK, progress)
}

func (h *Handler) UpdateProjectProgress(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	project, err := h.db.GetProjectByID(r.Context(), projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "Project not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to get project")
		return
	}

	curriculum, err := h.db.GetCurriculumByID(r.Context(), project.CurriculumID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get curriculum")
		return
	}

	if curriculum.UserID != user.ID {
		h.writeError(w, http.StatusForbidden, "Access denied")
		return
	}

	var req models.ProgressUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	progress, err := h.db.GetOrCreateProjectProgress(r.Context(), user.ID, projectID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get progress")
		return
	}

	progress.Status = req.Status
	progress.TimeSpentMinutes = req.TimeSpentMinutes

	if err := h.db.UpdateProjectProgress(r.Context(), progress); err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to update progress")
		return
	}

	h.writeJSON(w, http.StatusOK, progress)
}

func (h *Handler) CreateNote(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "projectId"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	project, err := h.db.GetProjectByID(r.Context(), projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "Project not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to get project")
		return
	}

	curriculum, err := h.db.GetCurriculumByID(r.Context(), project.CurriculumID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get curriculum")
		return
	}

	if curriculum.UserID != user.ID {
		h.writeError(w, http.StatusForbidden, "Access denied")
		return
	}

	var req models.NoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Content == "" {
		h.writeError(w, http.StatusBadRequest, "Content is required")
		return
	}

	note := &models.Note{
		UserID:    user.ID,
		ProjectID: projectID,
		Title:     req.Title,
		Content:   req.Content,
	}

	if err := h.db.CreateNote(r.Context(), note); err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to create note")
		return
	}

	h.writeJSON(w, http.StatusCreated, note)
}

func (h *Handler) GetNotes(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "projectId"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	project, err := h.db.GetProjectByID(r.Context(), projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "Project not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to get project")
		return
	}

	curriculum, err := h.db.GetCurriculumByID(r.Context(), project.CurriculumID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get curriculum")
		return
	}

	if curriculum.UserID != user.ID {
		h.writeError(w, http.StatusForbidden, "Access denied")
		return
	}

	notes, err := h.db.GetNotesByProjectID(r.Context(), user.ID, projectID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get notes")
		return
	}

	h.writeJSON(w, http.StatusOK, notes)
}

func (h *Handler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	noteID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	var req models.NoteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Content == "" {
		h.writeError(w, http.StatusBadRequest, "Content is required")
		return
	}

	note := &models.Note{
		ID:      noteID,
		UserID:  user.ID,
		Title:   req.Title,
		Content: req.Content,
	}

	if err := h.db.UpdateNote(r.Context(), note); err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to update note")
		return
	}

	h.writeJSON(w, http.StatusOK, note)
}

func (h *Handler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	noteID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	if err := h.db.DeleteNote(r.Context(), noteID, user.ID); err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to delete note")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateReflection(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "projectId"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	project, err := h.db.GetProjectByID(r.Context(), projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "Project not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to get project")
		return
	}

	curriculum, err := h.db.GetCurriculumByID(r.Context(), project.CurriculumID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get curriculum")
		return
	}

	if curriculum.UserID != user.ID {
		h.writeError(w, http.StatusForbidden, "Access denied")
		return
	}

	var req models.ReflectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Content == "" {
		h.writeError(w, http.StatusBadRequest, "Content is required")
		return
	}

	reflection := &models.Reflection{
		UserID:    user.ID,
		ProjectID: projectID,
		Content:   req.Content,
	}

	if err := h.db.CreateReflection(r.Context(), reflection); err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to create reflection")
		return
	}

	h.writeJSON(w, http.StatusCreated, reflection)
}

func (h *Handler) GetReflections(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "projectId"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	project, err := h.db.GetProjectByID(r.Context(), projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "Project not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to get project")
		return
	}

	curriculum, err := h.db.GetCurriculumByID(r.Context(), project.CurriculumID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get curriculum")
		return
	}

	if curriculum.UserID != user.ID {
		h.writeError(w, http.StatusForbidden, "Access denied")
		return
	}

	reflections, err := h.db.GetReflectionsByProjectID(r.Context(), user.ID, projectID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get reflections")
		return
	}

	h.writeJSON(w, http.StatusOK, reflections)
}

func (h *Handler) UpdateReflection(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	reflectionID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid reflection ID")
		return
	}

	var req models.ReflectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Content == "" {
		h.writeError(w, http.StatusBadRequest, "Content is required")
		return
	}

	reflection := &models.Reflection{
		ID:      reflectionID,
		UserID:  user.ID,
		Content: req.Content,
	}

	if err := h.db.UpdateReflection(r.Context(), reflection); err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to update reflection")
		return
	}

	h.writeJSON(w, http.StatusOK, reflection)
}

func (h *Handler) DeleteReflection(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	reflectionID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid reflection ID")
		return
	}

	if err := h.db.DeleteReflection(r.Context(), reflectionID, user.ID); err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to delete reflection")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) CreateTimeEntry(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "projectId"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	project, err := h.db.GetProjectByID(r.Context(), projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "Project not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to get project")
		return
	}

	curriculum, err := h.db.GetCurriculumByID(r.Context(), project.CurriculumID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get curriculum")
		return
	}

	if curriculum.UserID != user.ID {
		h.writeError(w, http.StatusForbidden, "Access denied")
		return
	}

	var req models.TimeEntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Minutes <= 0 {
		h.writeError(w, http.StatusBadRequest, "Minutes must be positive")
		return
	}

	entry := &models.TimeEntry{
		UserID:      user.ID,
		ProjectID:   projectID,
		Minutes:     req.Minutes,
		Description: req.Description,
	}

	if err := h.db.CreateTimeEntry(r.Context(), entry); err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to create time entry")
		return
	}

	h.writeJSON(w, http.StatusCreated, entry)
}

func (h *Handler) GetTimeEntries(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	projectID, err := uuid.Parse(chi.URLParam(r, "projectId"))
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	project, err := h.db.GetProjectByID(r.Context(), projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			h.writeError(w, http.StatusNotFound, "Project not found")
			return
		}
		h.writeError(w, http.StatusInternalServerError, "Failed to get project")
		return
	}

	curriculum, err := h.db.GetCurriculumByID(r.Context(), project.CurriculumID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get curriculum")
		return
	}

	if curriculum.UserID != user.ID {
		h.writeError(w, http.StatusForbidden, "Access denied")
		return
	}

	entries, err := h.db.GetTimeEntriesByProjectID(r.Context(), user.ID, projectID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get time entries")
		return
	}

	h.writeJSON(w, http.StatusOK, entries)
}

func (h *Handler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	user, err := h.getUserFromContext(r.Context())
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	analytics, err := h.db.GetAnalytics(r.Context(), user.ID)
	if err != nil {
		h.writeError(w, http.StatusInternalServerError, "Failed to get analytics")
		return
	}

	h.writeJSON(w, http.StatusOK, analytics)
}

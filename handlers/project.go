package handlers

import (
	"curriculum-tracker/middleware"
	"curriculum-tracker/models"
	"curriculum-tracker/services"
	"curriculum-tracker/utils"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ProjectHandler struct {
	projectService *services.ProjectService
	noteService    *services.NoteService
}

func NewProjectHandler(projectService *services.ProjectService, noteService *services.NoteService) *ProjectHandler {
	return &ProjectHandler{
		projectService: projectService,
		noteService:    noteService,
	}
}

func (h *ProjectHandler) CreateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	_, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	curriculumID, err := strconv.Atoi(vars["curriculumId"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid curriculum ID")
		return
	}

	var req models.CreateProjectRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.Name == "" || req.ProjectType == "" {
		utils.WriteError(w, http.StatusBadRequest, "Name and project type are required")
		return
	}

	project, err := h.projectService.CreateProject(curriculumID, req)
	if err != nil {
		log.Printf("error creating project: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create project")
		return
	}

	utils.WriteJSON(w, http.StatusCreated, project)
}

func (h *ProjectHandler) GetProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	project, err := h.projectService.GetProjectByID(userID, projectID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Project not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) UpdateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	var req models.UpdateProjectRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.Name == "" || req.ProjectType == "" {
		utils.WriteError(w, http.StatusBadRequest, "Name and project type are required")
		return
	}

	project, err := h.projectService.UpdateProject(userID, projectID, req)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Project not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) DeleteProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	err = h.projectService.DeleteProject(userID, projectID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Project not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Project deleted successfully"})
}

func (h *ProjectHandler) GetProjectNotes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	vars := mux.Vars(r)
	projectID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	notes, err := h.noteService.GetNotesByProjectID(userID, projectID)
	if err != nil {
		log.Printf("error getting notes: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch notes")
		return
	}

	utils.WriteJSON(w, http.StatusOK, notes)
}

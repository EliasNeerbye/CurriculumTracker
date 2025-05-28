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

type CurriculumHandler struct {
	curriculumService *services.CurriculumService
	projectService    *services.ProjectService
}

func NewCurriculumHandler(curriculumService *services.CurriculumService, projectService *services.ProjectService) *CurriculumHandler {
	return &CurriculumHandler{
		curriculumService: curriculumService,
		projectService:    projectService,
	}
}

func (h *CurriculumHandler) CreateCurriculum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.CreateCurriculumRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.Name == "" {
		utils.WriteError(w, http.StatusBadRequest, "Name is required")
		return
	}

	curriculum, err := h.curriculumService.CreateCurriculum(userID, req)
	if err != nil {
		log.Printf("error creating curriculum: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create curriculum")
		return
	}

	utils.WriteJSON(w, http.StatusCreated, curriculum)
}

func (h *CurriculumHandler) GetCurricula(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	curricula, err := h.curriculumService.GetCurriculumsByUserID(userID)
	if err != nil {
		log.Printf("error getting curricula: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch curricula")
		return
	}

	utils.WriteJSON(w, http.StatusOK, curricula)
}

func (h *CurriculumHandler) GetCurriculum(w http.ResponseWriter, r *http.Request) {
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
	curriculumID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid curriculum ID")
		return
	}

	curriculum, err := h.curriculumService.GetCurriculumByID(userID, curriculumID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Curriculum not found")
		return
	}

	projects, err := h.projectService.GetProjectsByCurriculumID(userID, curriculumID)
	if err != nil {
		log.Printf("error getting projects: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch projects")
		return
	}

	curriculum.Projects = projects
	utils.WriteJSON(w, http.StatusOK, curriculum)
}

func (h *CurriculumHandler) UpdateCurriculum(w http.ResponseWriter, r *http.Request) {
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
	curriculumID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid curriculum ID")
		return
	}

	var req models.UpdateCurriculumRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.Name == "" {
		utils.WriteError(w, http.StatusBadRequest, "Name is required")
		return
	}

	curriculum, err := h.curriculumService.UpdateCurriculum(userID, curriculumID, req)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Curriculum not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, curriculum)
}

func (h *CurriculumHandler) DeleteCurriculum(w http.ResponseWriter, r *http.Request) {
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
	curriculumID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid curriculum ID")
		return
	}

	err = h.curriculumService.DeleteCurriculum(userID, curriculumID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Curriculum not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Curriculum deleted successfully"})
}

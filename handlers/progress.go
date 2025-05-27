package handlers

import (
	"curriculum-tracker/middleware"
	"curriculum-tracker/models"
	"curriculum-tracker/services"
	"curriculum-tracker/utils"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ProgressHandler struct {
	progressService *services.ProgressService
}

func NewProgressHandler(progressService *services.ProgressService) *ProgressHandler {
	return &ProgressHandler{
		progressService: progressService,
	}
}

func (h *ProgressHandler) UpdateProgress(w http.ResponseWriter, r *http.Request) {
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
	projectID, err := strconv.Atoi(vars["projectId"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	var req models.UpdateProgressRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.Status == "" {
		utils.WriteError(w, http.StatusBadRequest, "Status is required")
		return
	}

	validStatuses := map[string]bool{
		models.StatusNotStarted: true,
		models.StatusInProgress: true,
		models.StatusCompleted:  true,
		models.StatusOnHold:     true,
		models.StatusAbandoned:  true,
	}

	if !validStatuses[req.Status] {
		utils.WriteError(w, http.StatusBadRequest, "Invalid status")
		return
	}

	if req.CompletionPercentage < 0 || req.CompletionPercentage > 100 {
		utils.WriteError(w, http.StatusBadRequest, "Completion percentage must be between 0 and 100")
		return
	}

	progress, err := h.progressService.UpdateProgress(userID, projectID, req)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to update progress")
		return
	}

	utils.WriteJSON(w, http.StatusOK, progress)
}

func (h *ProgressHandler) GetProgress(w http.ResponseWriter, r *http.Request) {
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
	projectID, err := strconv.Atoi(vars["projectId"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	progress, err := h.progressService.GetProgressByProjectID(userID, projectID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Progress not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, progress)
}

func (h *ProgressHandler) GetCurriculumProgress(w http.ResponseWriter, r *http.Request) {
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
	curriculumID, err := strconv.Atoi(vars["curriculumId"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid curriculum ID")
		return
	}

	progressList, err := h.progressService.GetProgressByCurriculumID(userID, curriculumID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch progress")
		return
	}

	utils.WriteJSON(w, http.StatusOK, progressList)
}

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

type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
}

func NewAnalyticsHandler(analyticsService *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
	}
}

func (h *AnalyticsHandler) CreateTimeEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req models.CreateTimeEntryRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.ProjectID == 0 || req.Minutes <= 0 {
		utils.WriteError(w, http.StatusBadRequest, "Project ID and positive minutes are required")
		return
	}

	if req.Date == "" {
		utils.WriteError(w, http.StatusBadRequest, "Date is required (YYYY-MM-DD format)")
		return
	}

	timeEntry, err := h.analyticsService.CreateTimeEntry(userID, req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Failed to create time entry: "+err.Error())
		return
	}

	utils.WriteJSON(w, http.StatusCreated, timeEntry)
}

func (h *AnalyticsHandler) GetProjectTimeEntries(w http.ResponseWriter, r *http.Request) {
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

	timeEntries, err := h.analyticsService.GetTimeEntriesByProjectID(userID, projectID)
	if err != nil {
		log.Printf("error getting time entries: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch time entries")
		return
	}

	utils.WriteJSON(w, http.StatusOK, timeEntries)
}

func (h *AnalyticsHandler) GetCurriculumTimeStats(w http.ResponseWriter, r *http.Request) {
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

	stats, err := h.analyticsService.GetTimeStatsByCurriculumID(userID, curriculumID)
	if err != nil {
		log.Printf("error getting time stats: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch time stats")
		return
	}

	utils.WriteJSON(w, http.StatusOK, stats)
}

func (h *AnalyticsHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		utils.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	stats, err := h.analyticsService.GetUserOverallStats(userID)
	if err != nil {
		log.Printf("error getting user stats: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "Failed to fetch user stats")
		return
	}

	utils.WriteJSON(w, http.StatusOK, stats)
}

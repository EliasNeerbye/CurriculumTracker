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

type NoteHandler struct {
	noteService *services.NoteService
}

func NewNoteHandler(noteService *services.NoteService) *NoteHandler {
	return &NoteHandler{
		noteService: noteService,
	}
}

func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var req models.CreateNoteRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.Content == "" {
		utils.WriteError(w, http.StatusBadRequest, "Content is required")
		return
	}

	if req.NoteType == "" {
		req.NoteType = models.NoteTypeNote
	}

	validTypes := map[string]bool{
		models.NoteTypeNote:       true,
		models.NoteTypeReflection: true,
		models.NoteTypeLearning:   true,
		models.NoteTypeQuestion:   true,
	}

	if !validTypes[req.NoteType] {
		utils.WriteError(w, http.StatusBadRequest, "Invalid note type")
		return
	}

	note, err := h.noteService.CreateNote(userID, projectID, req)
	if err != nil {
		log.Printf("error creating note: %v", err)
		utils.WriteError(w, http.StatusInternalServerError, "Failed to create note")
		return
	}

	utils.WriteJSON(w, http.StatusCreated, note)
}

func (h *NoteHandler) GetNote(w http.ResponseWriter, r *http.Request) {
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
	noteID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	note, err := h.noteService.GetNoteByID(userID, noteID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Note not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, note)
}

func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
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
	noteID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	var req models.UpdateNoteRequest
	if err := utils.ParseJSON(r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.Content == "" {
		utils.WriteError(w, http.StatusBadRequest, "Content is required")
		return
	}

	if req.NoteType == "" {
		req.NoteType = models.NoteTypeNote
	}

	validTypes := map[string]bool{
		models.NoteTypeNote:       true,
		models.NoteTypeReflection: true,
		models.NoteTypeLearning:   true,
		models.NoteTypeQuestion:   true,
	}

	if !validTypes[req.NoteType] {
		utils.WriteError(w, http.StatusBadRequest, "Invalid note type")
		return
	}

	note, err := h.noteService.UpdateNote(userID, noteID, req)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Note not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, note)
}

func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
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
	noteID, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "Invalid note ID")
		return
	}

	err = h.noteService.DeleteNote(userID, noteID)
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, "Note not found")
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "Note deleted successfully"})
}

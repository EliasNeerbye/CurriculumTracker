package routes

import (
	"curriculum-tracker/config"
	"curriculum-tracker/handlers"
	"curriculum-tracker/middleware"
	"curriculum-tracker/services"
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)

func Setup(db *sql.DB, cfg *config.Config) *mux.Router {
	authService := services.NewAuthService(db)
	curriculumService := services.NewCurriculumService(db)
	projectService := services.NewProjectService(db)
	progressService := services.NewProgressService(db)
	noteService := services.NewNoteService(db)
	analyticsService := services.NewAnalyticsService(db)

	authHandler := handlers.NewAuthHandler(authService, cfg)
	curriculumHandler := handlers.NewCurriculumHandler(curriculumService, projectService)
	projectHandler := handlers.NewProjectHandler(projectService, noteService)
	progressHandler := handlers.NewProgressHandler(progressService)
	noteHandler := handlers.NewNoteHandler(noteService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)

	router := mux.NewRouter()

	router.Use(middleware.CORS(cfg.AllowedOrigins))
	router.Use(middleware.Logging)

	api := router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/auth/register", authHandler.Register).Methods("POST", "OPTIONS")
	api.HandleFunc("/auth/login", authHandler.Login).Methods("POST", "OPTIONS")

	protected := api.NewRoute().Subrouter()
	protected.Use(middleware.Auth(cfg.JWTSecret))

	protected.HandleFunc("/auth/me", authHandler.Me).Methods("GET", "OPTIONS")

	protected.HandleFunc("/curricula", curriculumHandler.CreateCurriculum).Methods("POST", "OPTIONS")
	protected.HandleFunc("/curricula", curriculumHandler.GetCurricula).Methods("GET", "OPTIONS")
	protected.HandleFunc("/curricula/{id:[0-9]+}", curriculumHandler.GetCurriculum).Methods("GET", "OPTIONS")
	protected.HandleFunc("/curricula/{id:[0-9]+}", curriculumHandler.UpdateCurriculum).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/curricula/{id:[0-9]+}", curriculumHandler.DeleteCurriculum).Methods("DELETE", "OPTIONS")

	protected.HandleFunc("/curricula/{curriculumId:[0-9]+}/projects", projectHandler.CreateProject).Methods("POST", "OPTIONS")
	protected.HandleFunc("/projects/{id:[0-9]+}", projectHandler.GetProject).Methods("GET", "OPTIONS")
	protected.HandleFunc("/projects/{id:[0-9]+}", projectHandler.UpdateProject).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/projects/{id:[0-9]+}", projectHandler.DeleteProject).Methods("DELETE", "OPTIONS")
	protected.HandleFunc("/projects/{id:[0-9]+}/notes", projectHandler.GetProjectNotes).Methods("GET", "OPTIONS")

	protected.HandleFunc("/projects/{projectId:[0-9]+}/progress", progressHandler.UpdateProgress).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/projects/{projectId:[0-9]+}/progress", progressHandler.GetProgress).Methods("GET", "OPTIONS")
	protected.HandleFunc("/curricula/{curriculumId:[0-9]+}/progress", progressHandler.GetCurriculumProgress).Methods("GET", "OPTIONS")

	protected.HandleFunc("/projects/{projectId:[0-9]+}/notes", noteHandler.CreateNote).Methods("POST", "OPTIONS")
	protected.HandleFunc("/notes/{id:[0-9]+}", noteHandler.GetNote).Methods("GET", "OPTIONS")
	protected.HandleFunc("/notes/{id:[0-9]+}", noteHandler.UpdateNote).Methods("PUT", "OPTIONS")
	protected.HandleFunc("/notes/{id:[0-9]+}", noteHandler.DeleteNote).Methods("DELETE", "OPTIONS")

	protected.HandleFunc("/time-entries", analyticsHandler.CreateTimeEntry).Methods("POST", "OPTIONS")
	protected.HandleFunc("/projects/{projectId:[0-9]+}/time-entries", analyticsHandler.GetProjectTimeEntries).Methods("GET", "OPTIONS")
	protected.HandleFunc("/curricula/{curriculumId:[0-9]+}/time-stats", analyticsHandler.GetCurriculumTimeStats).Methods("GET", "OPTIONS")
	protected.HandleFunc("/analytics/user-stats", analyticsHandler.GetUserStats).Methods("GET", "OPTIONS")

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	return router
}

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

	router.Use(middleware.CORS)
	router.Use(middleware.Logging)

	api := router.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/auth/register", authHandler.Register)
	api.HandleFunc("/auth/login", authHandler.Login)

	protected := api.NewRoute().Subrouter()
	protected.Use(middleware.Auth(cfg.JWTSecret))

	protected.HandleFunc("/auth/me", authHandler.Me)

	protected.HandleFunc("/curricula", curriculumHandler.CreateCurriculum).Methods("OPTIONS", "POST")
	protected.HandleFunc("/curricula", curriculumHandler.GetCurricula).Methods("OPTIONS", "GET")
	protected.HandleFunc("/curricula/{id:[0-9]+}", curriculumHandler.GetCurriculum).Methods("OPTIONS", "GET")
	protected.HandleFunc("/curricula/{id:[0-9]+}", curriculumHandler.UpdateCurriculum).Methods("OPTIONS", "PUT")
	protected.HandleFunc("/curricula/{id:[0-9]+}", curriculumHandler.DeleteCurriculum).Methods("OPTIONS", "DELETE")

	protected.HandleFunc("/curricula/{curriculumId:[0-9]+}/projects", projectHandler.CreateProject).Methods("OPTIONS", "POST")
	protected.HandleFunc("/projects/{id:[0-9]+}", projectHandler.GetProject).Methods("OPTIONS", "GET")
	protected.HandleFunc("/projects/{id:[0-9]+}", projectHandler.UpdateProject).Methods("OPTIONS", "PUT")
	protected.HandleFunc("/projects/{id:[0-9]+}", projectHandler.DeleteProject).Methods("OPTIONS", "DELETE")
	protected.HandleFunc("/projects/{id:[0-9]+}/notes", projectHandler.GetProjectNotes).Methods("OPTIONS", "GET")

	protected.HandleFunc("/projects/{projectId:[0-9]+}/progress", progressHandler.UpdateProgress).Methods("OPTIONS", "PUT")
	protected.HandleFunc("/projects/{projectId:[0-9]+}/progress", progressHandler.GetProgress).Methods("OPTIONS", "GET")
	protected.HandleFunc("/curricula/{curriculumId:[0-9]+}/progress", progressHandler.GetCurriculumProgress).Methods("OPTIONS", "GET")

	protected.HandleFunc("/projects/{projectId:[0-9]+}/notes", noteHandler.CreateNote).Methods("OPTIONS", "POST")
	protected.HandleFunc("/notes/{id:[0-9]+}", noteHandler.GetNote).Methods("OPTIONS", "GET")
	protected.HandleFunc("/notes/{id:[0-9]+}", noteHandler.UpdateNote).Methods("OPTIONS", "PUT")
	protected.HandleFunc("/notes/{id:[0-9]+}", noteHandler.DeleteNote).Methods("OPTIONS", "DELETE")

	protected.HandleFunc("/time-entries", analyticsHandler.CreateTimeEntry).Methods("OPTIONS", "POST")
	protected.HandleFunc("/projects/{projectId:[0-9]+}/time-entries", analyticsHandler.GetProjectTimeEntries).Methods("OPTIONS", "GET")
	protected.HandleFunc("/curricula/{curriculumId:[0-9]+}/time-stats", analyticsHandler.GetCurriculumTimeStats).Methods("OPTIONS", "GET")
	protected.HandleFunc("/analytics/user-stats", analyticsHandler.GetUserStats).Methods("GET", "OPTIONS")

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	return router
}

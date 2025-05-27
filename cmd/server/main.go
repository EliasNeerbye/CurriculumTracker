package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"curriculum-tracker/internal/auth"
	"curriculum-tracker/internal/database"
	"curriculum-tracker/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		databaseURL = "postgres://postgres:password@localhost/curriculum_tracker?sslmode=disable"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-this-in-production"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	db, err := database.New(databaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	authConfig := auth.NewConfig(jwtSecret, 24*time.Hour)
	h := handlers.New(db, authConfig)

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Compress(5))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Route("/api", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)

		r.Route("/", func(r chi.Router) {
			r.Use(h.AuthMiddleware)

			r.Get("/profile", h.GetProfile)
			r.Get("/analytics", h.GetAnalytics)

			r.Route("/curricula", func(r chi.Router) {
				r.Get("/", h.GetCurricula)
				r.Post("/", h.CreateCurriculum)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", h.GetCurriculum)
					r.Put("/", h.UpdateCurriculum)
					r.Delete("/", h.DeleteCurriculum)
					r.Get("/projects", h.GetProjects)
					r.Post("/projects", h.CreateProject)
				})
			})

			r.Route("/projects", func(r chi.Router) {
				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", h.GetProject)
					r.Put("/", h.UpdateProject)
					r.Delete("/", h.DeleteProject)
					r.Get("/progress", h.GetProjectProgress)
					r.Put("/progress", h.UpdateProjectProgress)

					r.Route("/notes", func(r chi.Router) {
						r.Get("/", h.GetNotes)
						r.Post("/", h.CreateNote)
					})

					r.Route("/reflections", func(r chi.Router) {
						r.Get("/", h.GetReflections)
						r.Post("/", h.CreateReflection)
					})

					r.Route("/time-entries", func(r chi.Router) {
						r.Get("/", h.GetTimeEntries)
						r.Post("/", h.CreateTimeEntry)
					})
				})
			})

			r.Route("/notes", func(r chi.Router) {
				r.Route("/{id}", func(r chi.Router) {
					r.Put("/", h.UpdateNote)
					r.Delete("/", h.DeleteNote)
				})
			})

			r.Route("/reflections", func(r chi.Router) {
				r.Route("/{id}", func(r chi.Router) {
					r.Put("/", h.UpdateReflection)
					r.Delete("/", h.DeleteReflection)
				})
			})
		})
	})

	fileServer := http.FileServer(http.Dir("./web/"))
	r.Handle("/*", fileServer)

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

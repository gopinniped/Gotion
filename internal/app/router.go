package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/gopinniped/gotion/internal/transport/http/handler"
)

func NewRouter(userHandler *handler.UserHandler, taskHandler *handler.TaskHandler, authMW func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Route("/auth", func(r chi.Router) {
			r.With(chimw.AllowContentType("application/json")).Post("/register", userHandler.RegisterUser)
			r.With(chimw.AllowContentType("application/json")).Post("/login", userHandler.LoginUser)
		})

		r.Group(func(r chi.Router) {
			r.Use(authMW)

			r.Route("/users", func(r chi.Router) {
				r.Get("/{id}", userHandler.GetByID)
				r.With(chimw.AllowContentType("application/json")).Patch("/{id}", userHandler.UpdateUser)
				r.Delete("/{id}", userHandler.DeleteUser)
			})

			r.Route("/tasks", func(r chi.Router) {
				r.Get("/", taskHandler.GetMyTasks)
				r.With(chimw.AllowContentType("application/json")).Post("/", taskHandler.CreateTask)
				r.Get("/{id}", taskHandler.GetByID)
				r.With(chimw.AllowContentType("application/json")).Patch("/{id}", taskHandler.UpdateTask)
				r.Delete("/{id}", taskHandler.DeleteTask)
			})
		})
	})

	return r
}

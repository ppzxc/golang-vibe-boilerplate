package httphandler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	apptodo "github.com/ppzxc/golang-vibe-boilerplate/internal/app/todo"
)

// NewRouter creates and configures the Chi router with all routes and middleware.
// Per ppzxc RESTful Guidelines: no /v1/ prefix, Api-Version header for versioning.
func NewRouter(todoSvc *apptodo.Service) http.Handler {
	r := chi.NewRouter()

	// Global middleware (order matters)
	r.Use(Recovery)
	r.Use(RequestID)
	r.Use(Logger)
	r.Use(Versioning("2026-04-04"))
	r.Use(chimiddleware.Compress(5))

	todoHandler := NewTodoHandler(todoSvc)

	// Todo routes
	r.Route("/todos", func(r chi.Router) {
		r.Get("/", todoHandler.List)
		r.Post("/", todoHandler.Create)
		r.Route("/{todoId}", func(r chi.Router) {
			r.Get("/", todoHandler.Get)
			r.Patch("/", todoHandler.Update)
			r.Delete("/", todoHandler.Delete)
			// Custom action: POST /todos/{todoId}:complete
			// Chi uses pattern matching, colon in path is literal
		})
		// Register colon-action routes at parent level
		r.Post("/{todoId}:complete", todoHandler.Complete)
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	return r
}

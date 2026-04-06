package httphandler

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

// NewRouter creates and configures the Chi router with all routes and middleware.
func NewRouter(todoSvc TodoService) http.Handler {
	r := chi.NewRouter()

	// Global middleware (order matters)
	r.Use(Recovery)
	r.Use(RequestID)
	r.Use(Logger)
	r.Use(Versioning("2026-04-04"))
	r.Use(chimiddleware.Compress(5))

	todoHandler := NewTodoHandler(todoSvc)

	// Register generated handlers
	HandlerFromMux(todoHandler, r)

	// Documentation
	r.Get("/docs", SwaggerUI)
	r.Get("/openapi.yaml", OpenAPISpec)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
			slog.Warn("failed to write health response", "error", err)
		}
	})

	return r
}

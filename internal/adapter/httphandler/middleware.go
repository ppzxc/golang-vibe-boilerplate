// Package httphandler provides the HTTP inbound adapter using Chi router.
package httphandler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/ppzxc/golang-vibe-boilerplate/pkg/httputil"
	"github.com/ppzxc/golang-vibe-boilerplate/pkg/problemdetail"
)

// RequestID middleware ensures every request has a Request-Id.
// If the client sends one, it is adopted; otherwise a new UUID v4 is generated.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := httputil.RequestIDFromRequest(r)
		ctx := httputil.WithRequestID(r.Context(), id)
		w.Header().Set(httputil.HeaderRequestID, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Logger middleware logs each request with duration and status using slog.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)
		slog.InfoContext(r.Context(), "request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"duration_ms", time.Since(start).Milliseconds(),
			"request_id", httputil.RequestIDFromContext(r.Context()),
		)
	})
}

// Recovery middleware recovers from panics and returns a 500 ProblemDetail.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				slog.ErrorContext(r.Context(), "panic recovered", "panic", rec)
				id := httputil.RequestIDFromContext(r.Context())
				problemdetail.InternalServerError(r.URL.Path, id).Write(w)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Versioning middleware reads the Api-Version header and sets it on the response.
// Per ppzxc RESTful Guidelines: header-based versioning, no /v1/... prefix.
func Versioning(currentVersion string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			version := r.Header.Get("Api-Version")
			if version == "" {
				version = currentVersion
			}
			w.Header().Set("Api-Version", version)
			next.ServeHTTP(w, r)
		})
	}
}

// responseWriter wraps http.ResponseWriter to capture the status code.
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

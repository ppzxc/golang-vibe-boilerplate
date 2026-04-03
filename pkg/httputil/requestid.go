// Package httputil provides common HTTP utilities.
package httputil

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// HeaderRequestID is the request/response header name for request IDs.
// Per ppzxc RESTful Guidelines: no X- prefix.
const HeaderRequestID = "Request-Id"

type contextKey string

const requestIDKey contextKey = "requestID"

// NewRequestID generates a new UUID v4 request ID.
func NewRequestID() string {
	return uuid.New().String()
}

// WithRequestID stores the request ID in the context.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// RequestIDFromContext retrieves the request ID from the context.
// Returns empty string if not found.
func RequestIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(requestIDKey).(string)
	return id
}

// RequestIDFromRequest retrieves the Request-Id header or generates a new one.
func RequestIDFromRequest(r *http.Request) string {
	if id := r.Header.Get(HeaderRequestID); id != "" {
		return id
	}
	return NewRequestID()
}

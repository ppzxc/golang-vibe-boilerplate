// Package problemdetail implements RFC 9457 Problem Details for HTTP APIs.
package problemdetail

import (
	"encoding/json"
	"net/http"
)

// Problem represents an RFC 9457 Problem Details object.
type Problem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
	TraceID  string `json:"traceId,omitempty"`
}

// ContentType is the MIME type for Problem Details responses.
const ContentType = "application/problem+json"

// New creates a new Problem with the given status and title.
func New(status int, title string) *Problem {
	return &Problem{
		Type:   "about:blank",
		Title:  title,
		Status: status,
	}
}

// WithDetail adds a detail message to the Problem.
func (p *Problem) WithDetail(detail string) *Problem {
	p.Detail = detail
	return p
}

// WithInstance adds the request instance URI to the Problem.
func (p *Problem) WithInstance(instance string) *Problem {
	p.Instance = instance
	return p
}

// WithTraceID adds a trace ID (matching Request-Id header) to the Problem.
func (p *Problem) WithTraceID(traceID string) *Problem {
	p.TraceID = traceID
	return p
}

// WithType sets the problem type URI.
func (p *Problem) WithType(problemType string) *Problem {
	p.Type = problemType
	return p
}

// Write writes the Problem as JSON to the given ResponseWriter.
func (p *Problem) Write(w http.ResponseWriter) {
	w.Header().Set("Content-Type", ContentType)
	w.WriteHeader(p.Status)
	_ = json.NewEncoder(w).Encode(p)
}

// NotFound returns a 404 Problem.
func NotFound(instance, traceID string) *Problem {
	return New(http.StatusNotFound, "Not Found").
		WithInstance(instance).
		WithTraceID(traceID)
}

// BadRequest returns a 400 Problem.
func BadRequest(detail, instance, traceID string) *Problem {
	return New(http.StatusBadRequest, "Bad Request").
		WithDetail(detail).
		WithInstance(instance).
		WithTraceID(traceID)
}

// UnprocessableEntity returns a 422 Problem.
func UnprocessableEntity(detail, instance, traceID string) *Problem {
	return New(http.StatusUnprocessableEntity, "Unprocessable Entity").
		WithDetail(detail).
		WithInstance(instance).
		WithTraceID(traceID)
}

// InternalServerError returns a 500 Problem.
func InternalServerError(instance, traceID string) *Problem {
	return New(http.StatusInternalServerError, "Internal Server Error").
		WithInstance(instance).
		WithTraceID(traceID)
}

// Conflict returns a 409 Problem.
func Conflict(detail, instance, traceID string) *Problem {
	return New(http.StatusConflict, "Conflict").
		WithDetail(detail).
		WithInstance(instance).
		WithTraceID(traceID)
}

package httphandler

import (
	"errors"
	"net/http"

	domain "github.com/ppzxc/golang-vibe-boilerplate/internal/domain/todo"
	"github.com/ppzxc/golang-vibe-boilerplate/pkg/httputil"
	"github.com/ppzxc/golang-vibe-boilerplate/pkg/problemdetail"
)

// Problem is a helper to write a problem detail response.
func Problem(w http.ResponseWriter, r *http.Request, status int, title, detail string) {
	id := httputil.RequestIDFromContext(r.Context())
	instance := r.URL.Path
	p := problemdetail.New(status, title).
		WithDetail(detail).
		WithInstance(instance).
		WithTraceID(id)
	p.Write(w)
}

// writeError maps domain errors to HTTP Problem Details responses.
func writeError(w http.ResponseWriter, r *http.Request, err error) {
	id := httputil.RequestIDFromContext(r.Context())
	instance := r.URL.Path

	switch {
	case errors.Is(err, domain.ErrNotFound):
		problemdetail.NotFound(instance, id).Write(w)
	case errors.Is(err, domain.ErrAlreadyDone):
		problemdetail.UnprocessableEntity(err.Error(), instance, id).Write(w)
	case errors.Is(err, domain.ErrTitleRequired):
		problemdetail.BadRequest(err.Error(), instance, id).Write(w)
	default:
		problemdetail.InternalServerError(instance, id).Write(w)
	}
}

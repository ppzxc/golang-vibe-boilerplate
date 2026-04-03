package httphandler

import (
	"errors"
	"net/http"

	apptodo "github.com/ppzxc/golang-vibe-boilerplate/internal/app/todo"
	domain "github.com/ppzxc/golang-vibe-boilerplate/internal/domain/todo"
	"github.com/ppzxc/golang-vibe-boilerplate/pkg/httputil"
	"github.com/ppzxc/golang-vibe-boilerplate/pkg/problemdetail"
)

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
	case errors.Is(err, apptodo.ErrInvalidJSON):
		problemdetail.BadRequest(err.Error(), instance, id).Write(w)
	case errors.Is(err, apptodo.ErrInvalidPageSize):
		problemdetail.BadRequest(err.Error(), instance, id).Write(w)
	default:
		problemdetail.InternalServerError(instance, id).Write(w)
	}
}

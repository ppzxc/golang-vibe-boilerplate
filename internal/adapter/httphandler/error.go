package httphandler

import (
	"errors"
	"net/http"

	domain "github.com/ppzxc/golang-vibe-boilerplate/internal/domain/todo"
	"github.com/ppzxc/golang-vibe-boilerplate/pkg/problemdetail"
)

func (h *TodoHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	if errors.Is(err, domain.ErrNotFound) {
		problemdetail.New(http.StatusNotFound, "Todo Not Found").WithDetail(err.Error()).Write(w)
		return
	}
	if errors.Is(err, domain.ErrTitleRequired) || errors.Is(err, domain.ErrAlreadyDone) {
		problemdetail.New(http.StatusBadRequest, "Bad Request").WithDetail(err.Error()).Write(w)
		return
	}
	problemdetail.New(http.StatusInternalServerError, "Internal Server Error").WithDetail(err.Error()).Write(w)
}

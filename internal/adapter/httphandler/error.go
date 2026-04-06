package httphandler

import (
	"errors"
	"net/http"

	domain "github.com/ppzxc/golang-vibe-boilerplate/internal/domain/todo"
	"github.com/ppzxc/golang-vibe-boilerplate/pkg/problemdetail"
)

func (h *TodoHandler) handleError(w http.ResponseWriter, r *http.Request, err error) {
	var pd problemdetail.ProblemDetail

	if errors.Is(err, domain.ErrNotFound) {
		pd = problemdetail.New(http.StatusNotFound, "Todo not found", err.Error())
	} else if errors.Is(err, domain.ErrInvalidInput) {
		pd = problemdetail.New(http.StatusBadRequest, "Invalid input", err.Error())
	} else {
		pd = problemdetail.New(http.StatusInternalServerError, "Internal server error", "An unexpected error occurred")
	}

	pd.Write(w, r)
}

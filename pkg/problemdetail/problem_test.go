package problemdetail_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ppzxc/golang-vibe-boilerplate/pkg/problemdetail"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	p := problemdetail.New(http.StatusNotFound, "Not Found")
	assert.Equal(t, http.StatusNotFound, p.Status)
	assert.Equal(t, "Not Found", p.Title)
	assert.Equal(t, "about:blank", p.Type)
}

func TestProblem_Write(t *testing.T) {
	w := httptest.NewRecorder()
	p := problemdetail.NotFound("/todos/999", "trace-123")
	p.Write(w)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, problemdetail.ContentType, w.Header().Get("Content-Type"))

	var got problemdetail.Problem
	require.NoError(t, json.NewDecoder(w.Body).Decode(&got))
	assert.Equal(t, http.StatusNotFound, got.Status)
	assert.Equal(t, "/todos/999", got.Instance)
	assert.Equal(t, "trace-123", got.TraceID)
}

func TestProblem_Chaining(t *testing.T) {
	p := problemdetail.New(http.StatusBadRequest, "Bad Request").
		WithDetail("field is required").
		WithInstance("/todos").
		WithTraceID("abc-123").
		WithType("https://example.com/errors/validation")

	assert.Equal(t, "field is required", p.Detail)
	assert.Equal(t, "/todos", p.Instance)
	assert.Equal(t, "abc-123", p.TraceID)
	assert.Equal(t, "https://example.com/errors/validation", p.Type)
}

package pagination_test

import (
	"net/http/httptest"
	"testing"

	"github.com/ppzxc/golang-vibe-boilerplate/pkg/pagination"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRequest_Defaults(t *testing.T) {
	r := httptest.NewRequest("GET", "/todos", nil)
	req, err := pagination.ParseRequest(r)
	require.NoError(t, err)
	assert.Equal(t, pagination.DefaultPageSize, req.PageSize)
	assert.Empty(t, req.PageToken)
}

func TestParseRequest_CustomPageSize(t *testing.T) {
	r := httptest.NewRequest("GET", "/todos?pageSize=50", nil)
	req, err := pagination.ParseRequest(r)
	require.NoError(t, err)
	assert.Equal(t, 50, req.PageSize)
}

func TestParseRequest_MaxPageSize(t *testing.T) {
	r := httptest.NewRequest("GET", "/todos?pageSize=999", nil)
	req, err := pagination.ParseRequest(r)
	require.NoError(t, err)
	assert.Equal(t, pagination.MaxPageSize, req.PageSize)
}

func TestParseRequest_InvalidPageSize(t *testing.T) {
	r := httptest.NewRequest("GET", "/todos?pageSize=0", nil)
	_, err := pagination.ParseRequest(r)
	assert.Error(t, err)
}

func TestEncodeDecodeToken(t *testing.T) {
	original := "cursor-value-123"
	encoded := pagination.EncodeToken(original)
	assert.NotEqual(t, original, encoded)

	decoded, err := pagination.DecodeToken(encoded)
	require.NoError(t, err)
	assert.Equal(t, original, decoded)
}

func TestResponse_SetHeaders_WithNextToken(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/todos?pageSize=20", nil)

	resp := &pagination.Response{
		TotalCount: 100,
		NextToken:  pagination.EncodeToken("cursor-50"),
		BaseURL:    "http://localhost/todos",
	}
	resp.SetHeaders(w, r)

	assert.Equal(t, "100", w.Header().Get("Total-Count"))
	assert.Contains(t, w.Header().Get("Link"), `rel="next"`)
}

func TestResponse_SetHeaders_NoNextToken(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/todos", nil)

	resp := &pagination.Response{
		TotalCount: 5,
		BaseURL:    "http://localhost/todos",
	}
	resp.SetHeaders(w, r)

	assert.Equal(t, "5", w.Header().Get("Total-Count"))
	assert.Empty(t, w.Header().Get("Link"))
}

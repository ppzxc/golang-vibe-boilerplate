package httputil_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/ppzxc/golang-vibe-boilerplate/pkg/httputil"
	"github.com/stretchr/testify/assert"
)

func TestNewRequestID(t *testing.T) {
	id1 := httputil.NewRequestID()
	id2 := httputil.NewRequestID()
	assert.NotEmpty(t, id1)
	assert.NotEqual(t, id1, id2)
	assert.Len(t, id1, 36) // UUID v4 format
}

func TestWithRequestID_RoundTrip(t *testing.T) {
	ctx := context.Background()
	ctx = httputil.WithRequestID(ctx, "test-id-123")
	got := httputil.RequestIDFromContext(ctx)
	assert.Equal(t, "test-id-123", got)
}

func TestRequestIDFromContext_Missing(t *testing.T) {
	ctx := context.Background()
	got := httputil.RequestIDFromContext(ctx)
	assert.Empty(t, got)
}

func TestRequestIDFromRequest_FromHeader(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set(httputil.HeaderRequestID, "existing-id")
	got := httputil.RequestIDFromRequest(r)
	assert.Equal(t, "existing-id", got)
}

func TestRequestIDFromRequest_Generated(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	got := httputil.RequestIDFromRequest(r)
	assert.NotEmpty(t, got)
	assert.Len(t, got, 36) // UUID v4
}

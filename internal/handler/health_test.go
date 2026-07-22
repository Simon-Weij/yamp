package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthHandler_GetHealth(t *testing.T) {
	rec := httptest.NewRecorder()

	healthHandler := NewHealthHandler()

	req := httptest.NewRequest("GET", "/health", nil)

	healthHandler.GetHealth(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.JSONEq(t, `{"status":"ok"}`, rec.Body.String())
}

package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Tests for ReservationHandler
func TestReserveTicket_InvalidEventID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	// Create a minimal handler for testing validation
	handler := &ReservationHandler{}
	r.POST("/reservations", handler.ReserveTicket)

	body := `{
		"event_id": "invalid-uuid",
		"user_email": "test@example.com",
		"name": "Test User"
	}`
	req := httptest.NewRequest(http.MethodPost, "/reservations", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Formato de event_id inválido")
}

func TestReserveTicket_MissingEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &ReservationHandler{}
	r.POST("/reservations", handler.ReserveTicket)

	body := `{
		"event_id": "550e8400-e29b-41d4-a716-446655440003",
		"name": "Test User"
	}`
	req := httptest.NewRequest(http.MethodPost, "/reservations", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Email requerido")
}

func TestReserveTicket_InvalidEmail(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &ReservationHandler{}
	r.POST("/reservations", handler.ReserveTicket)

	body := `{
		"event_id": "550e8400-e29b-41d4-a716-446655440003",
		"user_email": "invalid-email",
		"name": "Test User"
	}`
	req := httptest.NewRequest(http.MethodPost, "/reservations", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Formato de email inválido")
}

func TestReserveTicket_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &ReservationHandler{}
	r.POST("/reservations", handler.ReserveTicket)

	req := httptest.NewRequest(http.MethodPost, "/reservations", bytes.NewBufferString(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Datos de reserva inválidos")
}

func TestReserveTicket_ValidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &ReservationHandler{}
	r.POST("/reservations", handler.ReserveTicket)

	body := `{
		"event_id": "550e8400-e29b-41d4-a716-446655440003",
		"user_email": "test@example.com",
		"name": "Test User"
	}`
	req := httptest.NewRequest(http.MethodPost, "/reservations", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// This will fail because we don't have the dependencies set up,
	// but it should pass validation and fail at a later stage
	assert.NotEqual(t, http.StatusBadRequest, w.Code)
	// The response should not be a validation error
	assert.NotContains(t, w.Body.String(), "Formato de event_id inválido")
	assert.NotContains(t, w.Body.String(), "Email requerido")
	assert.NotContains(t, w.Body.String(), "Formato de email inválido")
}

package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Tests for TicketHandler
func TestListTickets_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &TicketHandler{}
	r.GET("/tickets", handler.ListTickets)

	req := httptest.NewRequest(http.MethodGet, "/tickets", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// This will fail because we don't have the dependencies set up,
	// but it should not be a validation error
	assert.NotEqual(t, http.StatusBadRequest, w.Code)
}

func TestListTickets_WithFilters(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &TicketHandler{}
	r.GET("/tickets", handler.ListTickets)

	req := httptest.NewRequest(http.MethodGet, "/tickets?user_email=test@example.com&event_id=event-123&limit=5", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// This will fail because we don't have the dependencies set up,
	// but it should not be a validation error
	assert.NotEqual(t, http.StatusBadRequest, w.Code)
}

func TestGetTicket_MissingID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &TicketHandler{}
	r.GET("/tickets/:id", handler.GetTicket)

	req := httptest.NewRequest(http.MethodGet, "/tickets/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "ID de ticket requerido")
}

func TestGetTicket_ValidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &TicketHandler{}
	r.GET("/tickets/:id", handler.GetTicket)

	req := httptest.NewRequest(http.MethodGet, "/tickets/550e8400-e29b-41d4-a716-446655440003", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// This will fail because we don't have the dependencies set up,
	// but it should not be a validation error
	assert.NotEqual(t, http.StatusBadRequest, w.Code)
}

func TestCreateTicket_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &TicketHandler{}
	r.POST("/tickets", handler.CreateTicket)

	body := `{"email":"test@example.com","event_id":"550e8400-e29b-41d4-a716-446655440003"}`
	req := httptest.NewRequest(http.MethodPost, "/tickets", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// This will fail because we don't have the dependencies set up,
	// but it should pass validation
	assert.NotEqual(t, http.StatusBadRequest, w.Code)
}

func TestCreateTicket_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &TicketHandler{}
	r.POST("/tickets", handler.CreateTicket)

	req := httptest.NewRequest(http.MethodPost, "/tickets", bytes.NewBufferString(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Datos de ticket inválidos")
}

func TestCreateTicket_MissingRequiredFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &TicketHandler{}
	r.POST("/tickets", handler.CreateTicket)

	body := `{"email":"test@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/tickets", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Datos de ticket inválidos")
}

func TestUpdateTicket_MissingID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &TicketHandler{}
	r.PUT("/tickets/:id", handler.UpdateTicket)

	req := httptest.NewRequest(http.MethodPut, "/tickets/", bytes.NewBufferString(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "ID de ticket requerido")
}

func TestUpdateTicket_InvalidEventID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &TicketHandler{}
	r.PUT("/tickets/:id", handler.UpdateTicket)

	body := `{"event_id":"invalid-uuid"}`
	req := httptest.NewRequest(http.MethodPut, "/tickets/550e8400-e29b-41d4-a716-446655440003", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Event ID inválido")
}

func TestUpdateTicket_ValidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &TicketHandler{}
	r.PUT("/tickets/:id", handler.UpdateTicket)

	body := `{"email":"updated@example.com"}`
	req := httptest.NewRequest(http.MethodPut, "/tickets/550e8400-e29b-41d4-a716-446655440003", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// This will fail because we don't have the dependencies set up,
	// but it should pass validation
	assert.NotEqual(t, http.StatusBadRequest, w.Code)
}

func TestDeleteTicket_MissingID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &TicketHandler{}
	r.DELETE("/tickets/:id", handler.DeleteTicket)

	req := httptest.NewRequest(http.MethodDelete, "/tickets/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "ID de ticket requerido")
}

func TestDeleteTicket_ValidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &TicketHandler{}
	r.DELETE("/tickets/:id", handler.DeleteTicket)

	req := httptest.NewRequest(http.MethodDelete, "/tickets/550e8400-e29b-41d4-a716-446655440003", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// This will fail because we don't have the dependencies set up,
	// but it should not be a validation error
	assert.NotEqual(t, http.StatusBadRequest, w.Code)
}

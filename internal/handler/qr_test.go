package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Tests for QRHandler
func TestGetTicketQR_MissingID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &QRHandler{}
	r.GET("/tickets/:id/qr", handler.GetTicketQR)

	req := httptest.NewRequest(http.MethodGet, "/tickets//qr", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "ID de ticket requerido")
}

func TestGetTicketQR_ValidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &QRHandler{}
	r.GET("/tickets/:id/qr", handler.GetTicketQR)

	req := httptest.NewRequest(http.MethodGet, "/tickets/550e8400-e29b-41d4-a716-446655440003/qr", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// This will fail because we don't have the dependencies set up,
	// but it should not be a validation error
	assert.NotEqual(t, http.StatusBadRequest, w.Code)
}

func TestGetTicketQRFromS3_MissingID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &QRHandler{}
	r.GET("/tickets/:id/qr-s3", handler.GetTicketQRFromS3)

	req := httptest.NewRequest(http.MethodGet, "/tickets//qr-s3", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "ID de ticket requerido")
}

func TestGetTicketQRFromS3_ValidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &QRHandler{}
	r.GET("/tickets/:id/qr-s3", handler.GetTicketQRFromS3)

	req := httptest.NewRequest(http.MethodGet, "/tickets/550e8400-e29b-41d4-a716-446655440003/qr-s3", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// This will fail because we don't have the dependencies set up,
	// but it should not be a validation error
	assert.NotEqual(t, http.StatusBadRequest, w.Code)
}

func TestValidateQR_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &QRHandler{}
	r.POST("/qr/validate", handler.ValidateQR)

	req := httptest.NewRequest(http.MethodPost, "/qr/validate", bytes.NewBufferString(`invalid json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Contenido QR requerido")
}

func TestValidateQR_MissingQRContent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &QRHandler{}
	r.POST("/qr/validate", handler.ValidateQR)

	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/qr/validate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Contenido QR requerido")
}

func TestValidateQR_ValidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &QRHandler{}
	r.POST("/qr/validate", handler.ValidateQR)

	body := `{"qr_content":"TICKET:550e8400-e29b-41d4-a716-446655440003|EMAIL:test@example.com|CODE:TKT-12345678"}`
	req := httptest.NewRequest(http.MethodPost, "/qr/validate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// This will fail because we don't have the dependencies set up,
	// but it should pass validation
	assert.NotEqual(t, http.StatusBadRequest, w.Code)
}

func TestGenerateQRForTicket_MissingID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &QRHandler{}
	r.POST("/tickets/:id/qr", handler.GenerateQRForTicket)

	req := httptest.NewRequest(http.MethodPost, "/tickets//qr", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "ID de ticket requerido")
}

func TestGenerateQRForTicket_ValidID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	handler := &QRHandler{}
	r.POST("/tickets/:id/qr", handler.GenerateQRForTicket)

	req := httptest.NewRequest(http.MethodPost, "/tickets/550e8400-e29b-41d4-a716-446655440003/qr", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// This will fail because we don't have the dependencies set up,
	// but it should not be a validation error
	assert.NotEqual(t, http.StatusBadRequest, w.Code)
}

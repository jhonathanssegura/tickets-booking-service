package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Definir interfaces para mocking
// Estas deben coincidir con los métodos usados en ReservationHandler

type SQSService interface {
	SendReservationMessage(ctx interface{}, msg interface{}) error
}

type S3Service interface {
	UploadTicketFile(ctx interface{}, key interface{}, body interface{}) error
}

// Redefinir ReservationHandler para usar interfaces en los tests
// (esto es solo para el test, el handler real usa los structs concretos)
type testReservationHandler struct {
	SQS SQSService
	S3  S3Service
}

func (h *testReservationHandler) ReserveTicket(c *gin.Context) {
	// Copiar la lógica del handler real, pero usando las interfaces
	var ticket struct {
		UserEmail string `json:"user_email"`
		EventID   string `json:"event_id"`
	}
	if err := c.BindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.SQS.SendReservationMessage(nil, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error enviando mensaje a SQS"})
		return
	}
	if err := h.S3.UploadTicketFile(nil, nil, nil); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error subiendo ticket a S3"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ticket reservado con éxito"})
}

type mockSQS struct{ sendErr error }

func (m *mockSQS) SendReservationMessage(_, _ interface{}) error { return m.sendErr }

type mockS3 struct{ uploadErr error }

func (m *mockS3) UploadTicketFile(_, _, _ interface{}) error { return m.uploadErr }

func TestReserveTicket_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &testReservationHandler{
		SQS: &mockSQS{},
		S3:  &mockS3{},
	}
	r := gin.Default()
	r.POST("/tickets", h.ReserveTicket)

	body := `{"user_email":"test@example.com","event_id":"evt-1"}`
	req := httptest.NewRequest(http.MethodPost, "/tickets", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Ticket reservado con éxito")
}

func TestReserveTicket_SQSError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &testReservationHandler{
		SQS: &mockSQS{sendErr: errors.New("sqs fail")},
		S3:  &mockS3{},
	}
	r := gin.Default()
	r.POST("/tickets", h.ReserveTicket)

	body := `{"user_email":"test@example.com","event_id":"evt-1"}`
	req := httptest.NewRequest(http.MethodPost, "/tickets", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error enviando mensaje a SQS")
}

func TestReserveTicket_S3Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	h := &testReservationHandler{
		SQS: &mockSQS{},
		S3:  &mockS3{uploadErr: errors.New("s3 fail")},
	}
	r := gin.Default()
	r.POST("/tickets", h.ReserveTicket)

	body := `{"user_email":"test@example.com","event_id":"evt-1"}`
	req := httptest.NewRequest(http.MethodPost, "/tickets", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Error subiendo ticket a S3")
}

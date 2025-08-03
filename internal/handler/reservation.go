package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jhonathanssegura/ticket-reservation/internal/db"
	"github.com/jhonathanssegura/ticket-reservation/internal/model"
	"github.com/jhonathanssegura/ticket-reservation/internal/queue"
	"github.com/jhonathanssegura/ticket-reservation/internal/service"
	"github.com/jhonathanssegura/ticket-reservation/internal/storage"
)

type ReservationHandler struct {
	SQS *queue.SQSClient
	S3  *storage.S3Client
	DB  *db.DynamoClient
	QR  *service.QRService
}

func NewReservationHandler(sqs *queue.SQSClient, s3 *storage.S3Client, db *db.DynamoClient) *ReservationHandler {
	return &ReservationHandler{
		SQS: sqs,
		S3:  s3,
		DB:  db,
		QR:  service.NewQRService(),
	}
}

func (h *ReservationHandler) ReserveTicket(c *gin.Context) {
	var req struct {
		EventID   string `json:"event_id" binding:"required"`
		UserID    string `json:"user_id"`
		Email     string `json:"email"`
		UserEmail string `json:"user_email"`
		Name      string `json:"name"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Datos de reserva inválidos",
			"details": err.Error(),
			"expected_format": map[string]string{
				"event_id":   "UUID válido (ej: 550e8400-e29b-41d4-a716-446655440003)",
				"user_id":    "UUID válido (opcional, se genera automáticamente si no se proporciona)",
				"email":      "Email válido (opcional si se proporciona user_email)",
				"user_email": "Email válido (opcional si se proporciona email)",
				"name":       "Nombre del usuario (opcional, se usa 'Usuario Anónimo' por defecto)",
			},
		})
		return
	}

	eventID, err := uuid.Parse(req.EventID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":           "Formato de event_id inválido",
			"details":         fmt.Sprintf("El event_id '%s' no es un UUID válido", req.EventID),
			"expected_format": "UUID válido (ej: 550e8400-e29b-41d4-a716-446655440003)",
		})
		return
	}

	var userID uuid.UUID
	if req.UserID != "" {
		userID, err = uuid.Parse(req.UserID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":           "Formato de user_id inválido",
				"details":         fmt.Sprintf("El user_id '%s' no es un UUID válido", req.UserID),
				"expected_format": "UUID válido (ej: 550e8400-e29b-41d4-a716-446655440003)",
			})
			return
		}
	} else {
		userID = uuid.New() // Generate new UUID if not provided
	}

	// Handle email - use user_email if email is not provided
	userEmail := req.Email
	if userEmail == "" {
		userEmail = req.UserEmail
	}
	if userEmail == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Email requerido",
			"details": "Debe proporcionar un email válido usando 'email' o 'user_email'",
			"example": map[string]string{
				"email":      "usuario@ejemplo.com",
				"user_email": "usuario@ejemplo.com",
			},
		})
		return
	}

	if !strings.Contains(userEmail, "@") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":           "Formato de email inválido",
			"details":         fmt.Sprintf("El email '%s' no tiene un formato válido", userEmail),
			"expected_format": "usuario@dominio.com",
		})
		return
	}

	userName := req.Name
	if userName == "" {
		userName = "Usuario Anónimo"
	}

	// Generate UUID for ticket
	ticketID := uuid.New()
	now := time.Now()

	ticket := model.Ticket{
		ID:         ticketID,
		EventID:    eventID,
		UserID:     userID,
		Email:      userEmail,
		Name:       userName,
		TicketCode: fmt.Sprintf("TKT-%s", ticketID.String()[:8]),
		Status:     model.TicketStatusReserved,
		Price:      0.0, // This should be calculated
		ReservedAt: now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	// Generate QR code for the ticket
	qrData, err := h.QR.GenerateTicketQRPNG(ticket.ID, ticket.Email, ticket.TicketCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error generando código QR",
			"details": err.Error(),
		})
		return
	}

	qrS3Key := fmt.Sprintf("qrcodes/%s.png", ticket.ID)
	if err := h.S3.UploadTicketFile(context.Background(), qrS3Key, bytes.NewReader(qrData)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error subiendo código QR a S3",
			"details": err.Error(),
		})
		return
	}

	ticketContent := fmt.Sprintf(`TICKET INFORMATION
		==================
		Ticket ID: %s
		Event ID: %s
		User: %s (%s)
		Ticket Code: %s
		Status: %s
		Price: $%.2f
		Reserved At: %s
		QR Code: %s
		`,
		ticket.ID, ticket.EventID, ticket.Name, ticket.Email,
		ticket.TicketCode, ticket.Status, ticket.Price,
		ticket.ReservedAt.Format("2006-01-02 15:04:05"), qrS3Key)

	ticketS3Key := fmt.Sprintf("tickets/%s.txt", ticket.ID)
	if err := h.S3.UploadTicketFile(context.Background(), ticketS3Key, bytes.NewReader([]byte(ticketContent))); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Error subiendo archivo de ticket a S3",
			"details": err.Error(),
		})
		return
	}

	if err := h.DB.SaveTicket(ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Error guardando ticket en base de datos",
			"details":   err.Error(),
			"ticket_id": ticket.ID.String(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Ticket reservado con éxito",
		"ticket_id":   ticket.ID,
		"ticket_file": ticketS3Key,
		"qr_code":     qrS3Key,
		"ticket_info": map[string]interface{}{
			"event_id":    ticket.EventID,
			"user_id":     ticket.UserID,
			"email":       ticket.Email,
			"name":        ticket.Name,
			"ticket_code": ticket.TicketCode,
			"status":      ticket.Status,
			"price":       ticket.Price,
			"reserved_at": ticket.ReservedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

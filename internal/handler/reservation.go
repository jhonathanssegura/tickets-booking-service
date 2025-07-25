package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jhonathanssegura/ticket-reservation/internal/model"
	"github.com/jhonathanssegura/ticket-reservation/internal/queue"
	"github.com/jhonathanssegura/ticket-reservation/internal/storage"
)

type ReservationHandler struct {
	SQS *queue.SQSClient
	S3  *storage.S3Client
}

func NewReservationHandler(sqs *queue.SQSClient, s3 *storage.S3Client) *ReservationHandler {
	return &ReservationHandler{SQS: sqs, S3: s3}
}

func (h *ReservationHandler) ReserveTicket(c *gin.Context) {
	var ticket model.Ticket
	if err := c.BindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket.ID = fmt.Sprintf("tkt-%d", time.Now().UnixNano())
	ticket.ReservedAt = time.Now().Format(time.RFC3339)

	// Enviar mensaje a SQS
	sqsMsg := queue.TicketReservationMessage{
		ReservationID: ticket.ID,
		UserID:        ticket.UserEmail,
		EventID:       ticket.EventID,
		NumTickets:    1, // Suponiendo 1 ticket por reserva
	}
	if err := h.SQS.SendReservationMessage(context.Background(), sqsMsg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error enviando mensaje a SQS", "details": err.Error()})
		return
	}

	// Simular generación de PDF/ticket y subirlo a S3
	fileContent := []byte(fmt.Sprintf("Ticket ID: %s\nUser: %s\nEvent: %s\nFecha: %s", ticket.ID, ticket.UserEmail, ticket.EventID, ticket.ReservedAt))
	s3Key := fmt.Sprintf("tickets/%s.txt", ticket.ID)
	if err := h.S3.UploadTicketFile(context.Background(), s3Key, bytes.NewReader(fileContent)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error subiendo ticket a S3", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket reservado con éxito", "ticket_id": ticket.ID, "s3_key": s3Key})
}

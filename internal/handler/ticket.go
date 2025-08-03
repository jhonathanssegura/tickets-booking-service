package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jhonathanssegura/ticket-reservation/internal/db"
	"github.com/jhonathanssegura/ticket-reservation/internal/model"
)

type TicketHandler struct {
	DB *db.DynamoClient
}

func NewTicketHandler(db *db.DynamoClient) *TicketHandler {
	return &TicketHandler{DB: db}
}

func (h *TicketHandler) ListTickets(c *gin.Context) {
	userEmail := c.Query("user_email")
	eventID := c.Query("event_id")
	limitStr := c.Query("limit")
	limit := 10 // default limit
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	tickets, err := h.DB.GetTickets(userEmail, eventID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo tickets", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"tickets": tickets,
		"count":   len(tickets),
		"limit":   limit,
	})
}

func (h *TicketHandler) GetTicket(c *gin.Context) {
	ticketID := c.Param("id")
	if ticketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de ticket requerido"})
		return
	}

	ticket, err := h.DB.GetTicketByID(ticketID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ticket no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo ticket", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"ticket": ticket})
}

func (h *TicketHandler) CreateTicket(c *gin.Context) {
	var ticketData struct {
		Email   string `json:"email" binding:"required"`
		EventID string `json:"event_id" binding:"required"`
	}

	if err := c.BindJSON(&ticketData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de ticket inválidos", "details": err.Error()})
		return
	}

	ticketID := uuid.New()
	now := time.Now()

	ticket := &model.Ticket{
		ID:         ticketID,
		EventID:    uuid.MustParse(ticketData.EventID),
		UserID:     uuid.New(),
		Email:      ticketData.Email,
		Name:       "User Name",
		TicketCode: fmt.Sprintf("TKT-%s", ticketID.String()[:8]),
		Status:     model.TicketStatusReserved,
		Price:      0.0,
		ReservedAt: now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := h.DB.SaveTicket(*ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creando ticket", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Ticket creado con éxito",
		"ticket":  ticket,
	})
}

func (h *TicketHandler) UpdateTicket(c *gin.Context) {
	ticketID := c.Param("id")
	if ticketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de ticket requerido"})
		return
	}

	existingTicket, err := h.DB.GetTicketByID(ticketID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ticket no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error obteniendo ticket", "details": err.Error()})
		return
	}

	var updateData struct {
		Email   string `json:"email"`
		EventID string `json:"event_id"`
	}

	if err := c.BindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos de actualización inválidos", "details": err.Error()})
		return
	}

	if updateData.Email != "" {
		existingTicket.Email = updateData.Email
	}
	if updateData.EventID != "" {
		eventID, err := uuid.Parse(updateData.EventID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Event ID inválido", "details": err.Error()})
			return
		}
		existingTicket.EventID = eventID
	}

	if err := h.DB.SaveTicket(*existingTicket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error actualizando ticket", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ticket actualizado con éxito",
		"ticket":  existingTicket,
	})
}

func (h *TicketHandler) DeleteTicket(c *gin.Context) {
	ticketID := c.Param("id")
	if ticketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de ticket requerido"})
		return
	}

	_, err := h.DB.GetTicketByID(ticketID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ticket no encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error verificando ticket", "details": err.Error()})
		return
	}

	if err := h.DB.DeleteTicket(ticketID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error eliminando ticket", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Ticket eliminado con éxito"})
}

func generateTicketID() string {
	return "TICKET-" + strconv.FormatInt(time.Now().UnixNano(), 10)
}

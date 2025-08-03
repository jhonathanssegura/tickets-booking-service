package handler

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jhonathanssegura/ticket-reservation/internal/db"
	"github.com/jhonathanssegura/ticket-reservation/internal/service"
	"github.com/jhonathanssegura/ticket-reservation/internal/storage"
)

type QRHandler struct {
	DB *db.DynamoClient
	S3 *storage.S3Client
	QR *service.QRService
}

func NewQRHandler(db *db.DynamoClient, s3 *storage.S3Client) *QRHandler {
	return &QRHandler{
		DB: db,
		S3: s3,
		QR: service.NewQRService(),
	}
}

func (h *QRHandler) GetTicketQR(c *gin.Context) {
	ticketID := c.Param("id")
	if ticketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de ticket requerido"})
		return
	}

	ticket, err := h.DB.GetTicketByID(ticketID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket no encontrado"})
		return
	}

	qrData, err := h.QR.GenerateTicketQRPNG(ticket.ID, ticket.Email, ticket.TicketCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generando código QR", "details": err.Error()})
		return
	}

	c.Header("Content-Type", "image/png")
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=qr-%s.png", ticket.ID))
	c.Data(http.StatusOK, "image/png", qrData)
}

func (h *QRHandler) GetTicketQRFromS3(c *gin.Context) {
	ticketID := c.Param("id")
	if ticketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de ticket requerido"})
		return
	}

	_, err := h.DB.GetTicketByID(ticketID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket no encontrado"})
		return
	}

	qrS3Key := fmt.Sprintf("qrcodes/%s.png", ticketID)

	qrReader, err := h.S3.DownloadTicketFile(context.Background(), qrS3Key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Código QR no encontrado en S3"})
		return
	}
	defer qrReader.Close()

	qrData, err := io.ReadAll(qrReader)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error leyendo archivo QR"})
		return
	}

	c.Header("Content-Type", "image/png")
	c.Header("Content-Disposition", fmt.Sprintf("inline; filename=qr-%s.png", ticketID))
	c.Data(http.StatusOK, "image/png", qrData)
}

func (h *QRHandler) ValidateQR(c *gin.Context) {
	var req struct {
		QRContent string `json:"qr_content" binding:"required"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Contenido QR requerido"})
		return
	}

	isValid, err := h.QR.ValidateQRContent(req.QRContent)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Formato QR inválido", "details": err.Error()})
		return
	}

	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Código QR inválido"})
		return
	}

	parts := strings.Split(req.QRContent, "|")
	ticketInfo := make(map[string]string)

	for _, part := range parts {
		if strings.Contains(part, ":") {
			keyValue := strings.SplitN(part, ":", 2)
			if len(keyValue) == 2 {
				ticketInfo[keyValue[0]] = keyValue[1]
			}
		}
	}

	ticketID, exists := ticketInfo["TICKET"]
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de ticket no encontrado en QR"})
		return
	}

	ticket, err := h.DB.GetTicketByID(ticketID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket no encontrado en base de datos"})
		return
	}

	expectedQRContent := h.QR.GenerateQRContent(ticket.ID, ticket.Email, ticket.TicketCode, "")
	if req.QRContent != expectedQRContent {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Código QR no coincide con ticket"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":   true,
		"ticket":  ticket,
		"message": "Código QR válido",
	})
}

func (h *QRHandler) GenerateQRForTicket(c *gin.Context) {
	ticketID := c.Param("id")
	if ticketID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID de ticket requerido"})
		return
	}

	ticket, err := h.DB.GetTicketByID(ticketID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket no encontrado"})
		return
	}

	qrData, err := h.QR.GenerateTicketQRPNG(ticket.ID, ticket.Email, ticket.TicketCode)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generando código QR", "details": err.Error()})
		return
	}

	qrS3Key := fmt.Sprintf("qrcodes/%s.png", ticket.ID)
	if err := h.S3.UploadTicketFile(context.Background(), qrS3Key, bytes.NewReader(qrData)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error subiendo código QR a S3", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Código QR generado y subido exitosamente",
		"qr_code":   qrS3Key,
		"ticket_id": ticket.ID,
	})
}

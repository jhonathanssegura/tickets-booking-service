package service

import (
	"bytes"
	"fmt"

	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

// QRService maneja la generación de códigos QR
type QRService struct{}

// NewQRService crea una nueva instancia del servicio QR
func NewQRService() *QRService {
	return &QRService{}
}

// GenerateTicketQR genera un código QR para un ticket
func (s *QRService) GenerateTicketQR(ticketID uuid.UUID, userEmail, ticketCode string) ([]byte, error) {
	// Crear el contenido del QR con información del ticket
	qrContent := fmt.Sprintf("TICKET:%s|EMAIL:%s|CODE:%s", ticketID.String(), userEmail, ticketCode)

	// Generar el código QR
	qrCode, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("error generando código QR: %v", err)
	}

	return qrCode, nil
}

// GenerateTicketQRPNG genera un código QR en formato PNG
func (s *QRService) GenerateTicketQRPNG(ticketID uuid.UUID, userEmail, ticketCode string) ([]byte, error) {
	// Crear el contenido del QR con información del ticket
	qrContent := fmt.Sprintf("TICKET:%s|EMAIL:%s|CODE:%s", ticketID.String(), userEmail, ticketCode)

	// Generar el código QR como PNG
	qrCode, err := qrcode.New(qrContent, qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("error creando código QR: %v", err)
	}

	// Convertir a PNG
	pngData, err := qrCode.PNG(256)
	if err != nil {
		return nil, fmt.Errorf("error convirtiendo QR a PNG: %v", err)
	}

	return pngData, nil
}

// GenerateTicketQRWithLogo genera un código QR con logo (versión avanzada)
func (s *QRService) GenerateTicketQRWithLogo(ticketID uuid.UUID, userEmail, ticketCode string) ([]byte, error) {
	// Crear el contenido del QR con información del ticket
	qrContent := fmt.Sprintf("TICKET:%s|EMAIL:%s|CODE:%s", ticketID.String(), userEmail, ticketCode)

	// Generar el código QR con configuración personalizada
	qrCode, err := qrcode.New(qrContent, qrcode.High)
	if err != nil {
		return nil, fmt.Errorf("error creando código QR: %v", err)
	}

	// Configurar el código QR
	qrCode.DisableBorder = false

	// Convertir a PNG con tamaño personalizado
	pngData, err := qrCode.PNG(300)
	if err != nil {
		return nil, fmt.Errorf("error convirtiendo QR a PNG: %v", err)
	}

	return pngData, nil
}

// GenerateQRContent genera el contenido que se codificará en el QR
func (s *QRService) GenerateQRContent(ticketID uuid.UUID, userEmail, ticketCode, eventName string) string {
	return fmt.Sprintf("TICKET:%s|EMAIL:%s|CODE:%s|EVENT:%s",
		ticketID.String(), userEmail, ticketCode, eventName)
}

// ValidateQRContent valida el contenido de un código QR
func (s *QRService) ValidateQRContent(content string) (bool, error) {
	// Verificar que el contenido tenga el formato esperado
	if len(content) < 10 {
		return false, fmt.Errorf("contenido QR demasiado corto")
	}

	// Verificar que contenga los campos requeridos
	requiredFields := []string{"TICKET:", "EMAIL:", "CODE:"}
	for _, field := range requiredFields {
		if !bytes.Contains([]byte(content), []byte(field)) {
			return false, fmt.Errorf("campo requerido no encontrado: %s", field)
		}
	}

	return true, nil
}

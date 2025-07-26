package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// TicketReservationMessage representa el mensaje que se enviará a la cola SQS
// para procesar una reserva de ticket.
type TicketReservationMessage struct {
	ReservationID string `json:"reservation_id"`
	UserID        string `json:"user_id"`
	EventID       string `json:"event_id"`
	NumTickets    int    `json:"num_tickets"`
}

// SQSClient envuelve el cliente de AWS SQS y la URL de la cola.
type SQSClient struct {
	Client   *sqs.Client
	QueueURL string
}

// SendReservationMessage envía un mensaje de reserva de ticket a la cola SQS.
func (s *SQSClient) SendReservationMessage(ctx context.Context, msg TicketReservationMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("error marshaling SQS message: %w", err)
	}

	_, err = s.Client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.QueueURL),
		MessageBody: aws.String(string(body)),
	})
	if err != nil {
		return fmt.Errorf("error sending SQS message: %w", err)
	}
	return nil
}

// ReceiveReservationMessages recibe mensajes de reserva de ticket desde la cola SQS.
func (s *SQSClient) ReceiveReservationMessages(ctx context.Context, maxMessages int32) ([]TicketReservationMessage, error) {
	resp, err := s.Client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.QueueURL),
		MaxNumberOfMessages: maxMessages,
		WaitTimeSeconds:     10,
	})
	if err != nil {
		return nil, fmt.Errorf("error receiving SQS messages: %w", err)
	}

	var messages []TicketReservationMessage
	for _, m := range resp.Messages {
		var msg TicketReservationMessage
		if err := json.Unmarshal([]byte(*m.Body), &msg); err == nil {
			messages = append(messages, msg)
		}
	}
	return messages, nil
}

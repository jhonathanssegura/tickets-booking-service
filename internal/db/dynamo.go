package db

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/jhonathanssegura/ticket-reservation/internal/model"
)

type DynamoClient struct {
	Client *dynamodb.Client
}

func (d *DynamoClient) SaveTicket(ticket model.Ticket) error {
	fmt.Printf("Guardando ticket: ID=%s, EventID=%s, UserID=%s, Email=%s\n",
		ticket.ID.String(), ticket.EventID.String(), ticket.UserID.String(), ticket.Email)

	item := map[string]types.AttributeValue{
		"id":          &types.AttributeValueMemberS{Value: ticket.ID.String()},
		"event_id":    &types.AttributeValueMemberS{Value: ticket.EventID.String()},
		"user_id":     &types.AttributeValueMemberS{Value: ticket.UserID.String()},
		"email":       &types.AttributeValueMemberS{Value: ticket.Email},
		"name":        &types.AttributeValueMemberS{Value: ticket.Name},
		"ticket_code": &types.AttributeValueMemberS{Value: ticket.TicketCode},
		"status":      &types.AttributeValueMemberS{Value: ticket.Status},
		"price":       &types.AttributeValueMemberN{Value: fmt.Sprintf("%.2f", ticket.Price)},
		"reserved_at": &types.AttributeValueMemberS{Value: ticket.ReservedAt.Format(time.RFC3339)},
		"created_at":  &types.AttributeValueMemberS{Value: ticket.CreatedAt.Format(time.RFC3339)},
		"updated_at":  &types.AttributeValueMemberS{Value: ticket.UpdatedAt.Format(time.RFC3339)},
	}

	if ticket.CheckedInAt != nil {
		item["checked_in_at"] = &types.AttributeValueMemberS{Value: ticket.CheckedInAt.Format(time.RFC3339)}
	}

	if ticket.CheckedInBy != nil {
		item["checked_in_by"] = &types.AttributeValueMemberS{Value: ticket.CheckedInBy.String()}
	}

	_, err := d.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("tickets"),
		Item:      item,
	})

	if err != nil {
		var errorMsg string
		switch {
		case strings.Contains(err.Error(), "ResourceNotFoundException"):
			errorMsg = "La tabla 'tickets' no existe en DynamoDB. Verifique que LocalStack esté ejecutándose y la tabla haya sido creada."
		case strings.Contains(err.Error(), "RequestCanceled"):
			errorMsg = "Error de conexión con DynamoDB. Verifique que LocalStack esté ejecutándose en http://localhost:4566."
		case strings.Contains(err.Error(), "ConditionalCheckFailedException"):
			errorMsg = "El ticket ya existe en la base de datos."
		default:
			errorMsg = fmt.Sprintf("Error guardando ticket en DynamoDB: %v", err)
		}
		return fmt.Errorf(errorMsg)
	}

	return nil
}

func (d *DynamoClient) GetTicketByID(ticketID string) (*model.Ticket, error) {
	result, err := d.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("tickets"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: ticketID},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, errors.New("ticket not found")
	}

	ticket, err := d.unmarshalTicket(result.Item)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (d *DynamoClient) GetTickets(userEmail, eventID string, limit int) ([]model.Ticket, error) {
	scanInput := &dynamodb.ScanInput{
		TableName: aws.String("tickets"),
		Limit:     aws.Int32(int32(limit)),
	}

	if userEmail != "" || eventID != "" {
		filterExpressions := []string{}
		expressionAttributeNames := make(map[string]string)
		expressionAttributeValues := make(map[string]types.AttributeValue)

		if userEmail != "" {
			filterExpressions = append(filterExpressions, "#email = :email")
			expressionAttributeNames["#email"] = "email"
			expressionAttributeValues[":email"] = &types.AttributeValueMemberS{Value: userEmail}
		}

		if eventID != "" {
			filterExpressions = append(filterExpressions, "#event_id = :event_id")
			expressionAttributeNames["#event_id"] = "event_id"
			eventUUID, err := uuid.Parse(eventID)
			if err != nil {
				return nil, fmt.Errorf("invalid event ID format: %v", err)
			}
			expressionAttributeValues[":event_id"] = &types.AttributeValueMemberS{Value: eventUUID.String()}
		}

		scanInput.FilterExpression = aws.String(strings.Join(filterExpressions, " AND "))
		scanInput.ExpressionAttributeNames = expressionAttributeNames
		scanInput.ExpressionAttributeValues = expressionAttributeValues
	}

	result, err := d.Client.Scan(context.TODO(), scanInput)
	if err != nil {
		return nil, err
	}

	var tickets []model.Ticket
	for _, item := range result.Items {
		ticket, err := d.unmarshalTicket(item)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, *ticket)
	}

	return tickets, nil
}

func (d *DynamoClient) DeleteTicket(ticketID string) error {
	_, err := d.Client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String("tickets"),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{Value: ticketID},
		},
	})
	return err
}

func (d *DynamoClient) unmarshalTicket(item map[string]types.AttributeValue) (*model.Ticket, error) {
	ticket := &model.Ticket{}

	if idVal, ok := item["id"].(*types.AttributeValueMemberS); ok {
		id, err := uuid.Parse(idVal.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid ticket ID: %v", err)
		}
		ticket.ID = id
	}

	if eventIDVal, ok := item["event_id"].(*types.AttributeValueMemberS); ok {
		eventID, err := uuid.Parse(eventIDVal.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid event ID: %v", err)
		}
		ticket.EventID = eventID
	}

	if userIDVal, ok := item["user_id"].(*types.AttributeValueMemberS); ok {
		userID, err := uuid.Parse(userIDVal.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid user ID: %v", err)
		}
		ticket.UserID = userID
	}

	if emailVal, ok := item["email"].(*types.AttributeValueMemberS); ok {
		ticket.Email = emailVal.Value
	}

	if nameVal, ok := item["name"].(*types.AttributeValueMemberS); ok {
		ticket.Name = nameVal.Value
	}

	if ticketCodeVal, ok := item["ticket_code"].(*types.AttributeValueMemberS); ok {
		ticket.TicketCode = ticketCodeVal.Value
	}

	if statusVal, ok := item["status"].(*types.AttributeValueMemberS); ok {
		ticket.Status = statusVal.Value
	}

	if priceVal, ok := item["price"].(*types.AttributeValueMemberN); ok {
		price, err := strconv.ParseFloat(priceVal.Value, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid price: %v", err)
		}
		ticket.Price = price
	}

	if reservedAtVal, ok := item["reserved_at"].(*types.AttributeValueMemberS); ok {
		reservedAt, err := time.Parse(time.RFC3339, reservedAtVal.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid reserved_at time: %v", err)
		}
		ticket.ReservedAt = reservedAt
	}

	if createdAtVal, ok := item["created_at"].(*types.AttributeValueMemberS); ok {
		createdAt, err := time.Parse(time.RFC3339, createdAtVal.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid created_at time: %v", err)
		}
		ticket.CreatedAt = createdAt
	}

	if updatedAtVal, ok := item["updated_at"].(*types.AttributeValueMemberS); ok {
		updatedAt, err := time.Parse(time.RFC3339, updatedAtVal.Value)
		if err != nil {
			return nil, fmt.Errorf("invalid updated_at time: %v", err)
		}
		ticket.UpdatedAt = updatedAt
	}

	return ticket, nil
}

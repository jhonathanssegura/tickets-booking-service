package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/jhonathanssegura/ticket-reservation/internal/awsconfig"
	"github.com/jhonathanssegura/ticket-reservation/internal/model"
)

func main() {
	// Cargar configuración AWS
	cfg, err := awsconfig.LoadAWSConfig()
	if err != nil {
		log.Fatalf("Error cargando configuración AWS: %v", err)
	}

	// Crear cliente DynamoDB
	dynamoClient := dynamodb.NewFromConfig(cfg)

	// Generar UUIDs para eventos
	eventIDs := map[string]uuid.UUID{
		"evt-concierto-rock":   uuid.MustParse("550e8400-e29b-41d4-a716-446655440001"),
		"evt-teatro-clasico":   uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"),
		"evt-deportes-futbol":  uuid.MustParse("550e8400-e29b-41d4-a716-446655440003"),
		"evt-cine-estreno":     uuid.MustParse("550e8400-e29b-41d4-a716-446655440004"),
		"evt-conferencia-tech": uuid.MustParse("550e8400-e29b-41d4-a716-446655440005"),
	}

	// Datos de prueba
	tickets := []model.Ticket{
		{
			ID:         uuid.MustParse("550e8400-e29b-41d4-a716-446655440101"),
			EventID:    eventIDs["evt-concierto-rock"],
			UserID:     uuid.New(),
			Email:      "juan.perez@example.com",
			Name:       "Juan Pérez",
			TicketCode: "TKT-001",
			Status:     model.TicketStatusReserved,
			Price:      75.00,
			ReservedAt: time.Now().Add(-24 * time.Hour),
			CreatedAt:  time.Now().Add(-24 * time.Hour),
			UpdatedAt:  time.Now().Add(-24 * time.Hour),
		},
		{
			ID:         uuid.MustParse("550e8400-e29b-41d4-a716-446655440102"),
			EventID:    eventIDs["evt-concierto-rock"],
			UserID:     uuid.New(),
			Email:      "maria.garcia@example.com",
			Name:       "María García",
			TicketCode: "TKT-002",
			Status:     model.TicketStatusReserved,
			Price:      75.00,
			ReservedAt: time.Now().Add(-12 * time.Hour),
			CreatedAt:  time.Now().Add(-12 * time.Hour),
			UpdatedAt:  time.Now().Add(-12 * time.Hour),
		},
		{
			ID:         uuid.MustParse("550e8400-e29b-41d4-a716-446655440103"),
			EventID:    eventIDs["evt-teatro-clasico"],
			UserID:     uuid.New(),
			Email:      "carlos.rodriguez@example.com",
			Name:       "Carlos Rodríguez",
			TicketCode: "TKT-003",
			Status:     model.TicketStatusReserved,
			Price:      45.00,
			ReservedAt: time.Now().Add(-6 * time.Hour),
			CreatedAt:  time.Now().Add(-6 * time.Hour),
			UpdatedAt:  time.Now().Add(-6 * time.Hour),
		},
		{
			ID:         uuid.MustParse("550e8400-e29b-41d4-a716-446655440104"),
			EventID:    eventIDs["evt-teatro-clasico"],
			UserID:     uuid.New(),
			Email:      "ana.lopez@example.com",
			Name:       "Ana López",
			TicketCode: "TKT-004",
			Status:     model.TicketStatusReserved,
			Price:      45.00,
			ReservedAt: time.Now().Add(-3 * time.Hour),
			CreatedAt:  time.Now().Add(-3 * time.Hour),
			UpdatedAt:  time.Now().Add(-3 * time.Hour),
		},
		{
			ID:         uuid.MustParse("550e8400-e29b-41d4-a716-446655440105"),
			EventID:    eventIDs["evt-deportes-futbol"],
			UserID:     uuid.New(),
			Email:      "juan.perez@example.com",
			Name:       "Juan Pérez",
			TicketCode: "TKT-005",
			Status:     model.TicketStatusReserved,
			Price:      30.00,
			ReservedAt: time.Now().Add(-1 * time.Hour),
			CreatedAt:  time.Now().Add(-1 * time.Hour),
			UpdatedAt:  time.Now().Add(-1 * time.Hour),
		},
		{
			ID:         uuid.MustParse("550e8400-e29b-41d4-a716-446655440106"),
			EventID:    eventIDs["evt-deportes-futbol"],
			UserID:     uuid.New(),
			Email:      "lucia.martinez@example.com",
			Name:       "Lucía Martínez",
			TicketCode: "TKT-006",
			Status:     model.TicketStatusReserved,
			Price:      30.00,
			ReservedAt: time.Now().Add(-30 * time.Minute),
			CreatedAt:  time.Now().Add(-30 * time.Minute),
			UpdatedAt:  time.Now().Add(-30 * time.Minute),
		},
		{
			ID:         uuid.MustParse("550e8400-e29b-41d4-a716-446655440107"),
			EventID:    eventIDs["evt-cine-estreno"],
			UserID:     uuid.New(),
			Email:      "pedro.sanchez@example.com",
			Name:       "Pedro Sánchez",
			TicketCode: "TKT-007",
			Status:     model.TicketStatusReserved,
			Price:      12.00,
			ReservedAt: time.Now().Add(-15 * time.Minute),
			CreatedAt:  time.Now().Add(-15 * time.Minute),
			UpdatedAt:  time.Now().Add(-15 * time.Minute),
		},
		{
			ID:         uuid.MustParse("550e8400-e29b-41d4-a716-446655440108"),
			EventID:    eventIDs["evt-cine-estreno"],
			UserID:     uuid.New(),
			Email:      "sofia.hernandez@example.com",
			Name:       "Sofía Hernández",
			TicketCode: "TKT-008",
			Status:     model.TicketStatusReserved,
			Price:      12.00,
			ReservedAt: time.Now().Add(-5 * time.Minute),
			CreatedAt:  time.Now().Add(-5 * time.Minute),
			UpdatedAt:  time.Now().Add(-5 * time.Minute),
		},
		{
			ID:         uuid.MustParse("550e8400-e29b-41d4-a716-446655440109"),
			EventID:    eventIDs["evt-conferencia-tech"],
			UserID:     uuid.New(),
			Email:      "roberto.diaz@example.com",
			Name:       "Roberto Díaz",
			TicketCode: "TKT-009",
			Status:     model.TicketStatusReserved,
			Price:      150.00,
			ReservedAt: time.Now().Add(-2 * time.Minute),
			CreatedAt:  time.Now().Add(-2 * time.Minute),
			UpdatedAt:  time.Now().Add(-2 * time.Minute),
		},
		{
			ID:         uuid.MustParse("550e8400-e29b-41d4-a716-446655440110"),
			EventID:    eventIDs["evt-conferencia-tech"],
			UserID:     uuid.New(),
			Email:      "elena.morales@example.com",
			Name:       "Elena Morales",
			TicketCode: "TKT-010",
			Status:     model.TicketStatusReserved,
			Price:      150.00,
			ReservedAt: time.Now(),
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		},
	}

	fmt.Println("🌱 Cargando datos de prueba en DynamoDB...")
	fmt.Printf("📊 Insertando %d tickets...\n", len(tickets))

	// Insertar tickets en DynamoDB
	for i, ticket := range tickets {
		// Convert UUID fields to strings for DynamoDB
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

		_, err = dynamoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
			TableName: aws.String("tickets"),
			Item:      item,
		})

		if err != nil {
			log.Printf("Error insertando ticket %d: %v", i+1, err)
		} else {
			fmt.Printf("✅ Ticket %s insertado correctamente\n", ticket.TicketCode)
		}
	}

	fmt.Println("\n🎉 Datos de prueba cargados exitosamente!")
	fmt.Println("\n📋 Resumen de datos cargados:")
	fmt.Println("   • 10 tickets de prueba")
	fmt.Println("   • 5 eventos diferentes:")
	fmt.Println("     - evt-concierto-rock (2 tickets) - $75.00")
	fmt.Println("     - evt-teatro-clasico (2 tickets) - $45.00")
	fmt.Println("     - evt-deportes-futbol (2 tickets) - $30.00")
	fmt.Println("     - evt-cine-estreno (2 tickets) - $12.00")
	fmt.Println("     - evt-conferencia-tech (2 tickets) - $150.00")
	fmt.Println("   • 10 usuarios diferentes")

	fmt.Println("\n🧪 Pruebas que puedes realizar:")
	fmt.Println("1. Listar todos los tickets:")
	fmt.Println("   curl -X GET http://localhost:8080/api/tickets")
	fmt.Println("\n2. Filtrar por evento:")
	fmt.Println("   curl -X GET 'http://localhost:8080/api/tickets?event_id=550e8400-e29b-41d4-a716-446655440001'")
	fmt.Println("\n3. Filtrar por usuario:")
	fmt.Println("   curl -X GET 'http://localhost:8080/api/tickets?user_email=juan.perez@example.com'")
	fmt.Println("\n4. Obtener un ticket específico:")
	fmt.Println("   curl -X GET http://localhost:8080/api/tickets/550e8400-e29b-41d4-a716-446655440101")
	fmt.Println("\n5. Crear un nuevo ticket:")
	fmt.Println("   curl -X POST http://localhost:8080/api/tickets \\")
	fmt.Println("     -H 'Content-Type: application/json' \\")
	fmt.Println("     -d '{\"user_email\":\"nuevo@example.com\",\"event_id\":\"550e8400-e29b-41d4-a716-446655440001\"}'")
}

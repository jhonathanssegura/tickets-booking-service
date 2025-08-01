package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gin-gonic/gin"
	"github.com/jhonathanssegura/ticket-reservation/internal/awsconfig"
	"github.com/jhonathanssegura/ticket-reservation/internal/db"
	"github.com/jhonathanssegura/ticket-reservation/internal/handler"
	"github.com/jhonathanssegura/ticket-reservation/internal/queue"
	"github.com/jhonathanssegura/ticket-reservation/internal/storage"
)

func main() {
	cfg, err := awsconfig.LoadAWSConfig()
	if err != nil {
		log.Fatalf("Error cargando configuración AWS: %v", err)
	}

	queueURL := "http://localhost:4566/000000000000/ticket-queue"
	bucketName := "ticket-bucket"

	s3Client := s3.NewFromConfig(cfg, func(o *s3.Options) { o.UsePathStyle = true })

	log.Println("Verificando bucket S3...")
	_, err = s3Client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		log.Printf("Bucket S3 '%s' no existe, creándolo...", bucketName)
		_, err = s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})
		if err != nil {
			log.Fatalf("Error creando bucket S3: %v", err)
		}
		log.Printf("Bucket S3 '%s' creado exitosamente", bucketName)
	} else {
		log.Printf("Bucket S3 '%s' ya existe", bucketName)
	}

	sqsClient := &queue.SQSClient{
		Client:   sqs.NewFromConfig(cfg),
		QueueURL: queueURL,
	}

	storageClient := &storage.S3Client{
		Client:     s3Client,
		BucketName: bucketName,
	}

	dynamoClient := &db.DynamoClient{
		Client: dynamodb.NewFromConfig(cfg),
	}

	handlerReserva := handler.NewReservationHandler(sqsClient, storageClient, dynamoClient)
	handlerTicket := handler.NewTicketHandler(dynamoClient)
	handlerQR := handler.NewQRHandler(dynamoClient, storageClient)

	r := gin.Default()

	api := r.Group("/api")
	{
		// Ticket management endpoints
		api.GET("/tickets", handlerTicket.ListTickets)
		api.GET("/tickets/:id", handlerTicket.GetTicket)
		api.POST("/tickets", handlerTicket.CreateTicket)
		api.PUT("/tickets/:id", handlerTicket.UpdateTicket)
		api.DELETE("/tickets/:id", handlerTicket.DeleteTicket)
		// Reservation endpoint
		api.POST("/reservations", handlerReserva.ReserveTicket)
		// QR code endpoints
		api.GET("/tickets/:id/qr", handlerQR.GetTicketQR)
		api.GET("/tickets/:id/qr-s3", handlerQR.GetTicketQRFromS3)
		api.POST("/qr/validate", handlerQR.ValidateQR)
		api.POST("/tickets/:id/qr", handlerQR.GenerateQRForTicket)
	}

	log.Println("🚀 Iniciando servidor en puerto 8080...")
	r.Run(":8080")
}

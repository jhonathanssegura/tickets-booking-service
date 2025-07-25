package main

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/gin-gonic/gin"
	"github.com/jhonathanssegura/ticket-reservation/internal/awsconfig"
	"github.com/jhonathanssegura/ticket-reservation/internal/handler"
	"github.com/jhonathanssegura/ticket-reservation/internal/queue"
	"github.com/jhonathanssegura/ticket-reservation/internal/storage"
)

func main() {
	// Cargar configuración AWS
	cfg, err := awsconfig.LoadAWSConfig()
	if err != nil {
		log.Fatalf("Error cargando configuración AWS: %v", err)
	}

	// Inicializar clientes SQS y S3
	queueURL := "http://localhost:4566/000000000000/ticket-queue" // Cambia esto por tu URL real
	bucketName := "ticket-bucket"                                 // Cambia esto por tu bucket real

	sqsClient := &queue.SQSClient{
		Client:   sqs.NewFromConfig(cfg),
		QueueURL: queueURL,
	}
	s3Client := &storage.S3Client{
		Client:     s3.NewFromConfig(cfg),
		BucketName: bucketName,
	}

	handlerReserva := handler.NewReservationHandler(sqsClient, s3Client)

	r := gin.Default()
	r.POST("/tickets", handlerReserva.ReserveTicket)
	r.Run(":8080")
}

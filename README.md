# Ticket Booking
App to Booking tickets for event with Golang

## Tech Stack

### LocalStack
* DynamoDB: Almacenamiento de las reservas.
* SQS: Procesamiento asíncrono para las confirmaciones o eventos.
* S3: Almacenamiento de archivos tipo ticket PDF.
* ECS:
* RDS:
* API Gateway:

### Golang
* Gin

## Levantar LocalStack con Docker Compose

    docker-compose up -d

Verificar que LocalStack está corriendo:

    docker-compose logs localhost

## Crear recursos en LocalStack (SQS y S3)

Ejecutar el script aws-config.sh para crear el bucket S3 y la cola SQS necesaria.

    bash aws-config.sh

## Levantar la API de Go

Instalar dependencias

    go mod tidy

Ejecutar la API

    go run cmd/main.go

## Probar el Flujo

Se puede probar el endpoint con *curl* utilizando la instrucción

    curl -X POST http://localhost:8080/tickets \
        -H "Content-Type: application/json" \
        -d '{"user_email":"test@example.com","event_id":"evt-1"}'

Verificar en LocalStack

Ver mensajes en SQS:

    aws --endpoint-url=http://localhost:4566 sqs receive-message --queue-url http://localhost:4566/000000000000/ticket-queue

Ver archvios en S3:

    aws --endpoint-url=http://localhost:4566 s3 ls s3://ticket-bucket/tickets/
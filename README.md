# Ticket Booking API

API para gestionar tickets y reservas de eventos, integrando AWS SQS, S3 y DynamoDB mediante LocalStack.

## Tech Stack

### LocalStack
* DynamoDB: Almacenamiento de tickets y reservas
* SQS (Simple Queue Service): Procesamiento asíncrono para confirmaciones
* S3 (Simple Storage Service): Almacenamiento de archivos tipo ticket PDF y códigos QR

### Golang
* Gin: Framework web
* AWS SDK v2: Integración con servicios AWS

## Levantar LocalStack con Docker Compose

```bash
docker-compose up -d
```

Verificar que LocalStack está corriendo:

```bash
docker-compose logs localstack
```

## Crear recursos en LocalStack

Ejecutar el script aws-config.sh para crear el bucket S3, la cola SQS y la tabla DynamoDB necesaria:

```bash
bash aws-config.sh
```

## Levantar la API de Go

Instalar dependencias:

```bash
go mod tidy
```

Cargar la data de prueba
```bash
go run scripts/seed-data.go
```

Ejecutar la API:

```bash
go run cmd/main.go
```

## Verificar en LocalStack

### Ver mensajes en SQS:
```bash
aws --endpoint-url=http://localhost:4566 sqs receive-message --queue-url http://localhost:4566/000000000000/ticket-queue
```

### Ver archivos en S3:
```bash
aws --endpoint-url=http://localhost:4566 s3 ls s3://ticket-bucket/tickets/
```

### Ver tickets en DynamoDB:
```bash
aws --endpoint-url=http://localhost:4566 dynamodb scan --table-name tickets
```

## Documentación

- **Swagger**: Disponible en `docs/swagger.yaml`
- **Postman Collection**: Disponible en `docs/ticket-booking.postman_collection.json`

## Estructura del Proyecto

```
ticket-booking/
├── cmd/
│   └── main.go              # Punto de entrada de la aplicación
├── internal/
│   ├── awsconfig/           # Configuración de AWS
│   ├── db/                  # Cliente de DynamoDB
│   ├── handler/             # Handlers HTTP
│   ├── model/               # Modelos de datos
│   ├── queue/               # Cliente de SQS
│   └── storage/             # Cliente de S3
```
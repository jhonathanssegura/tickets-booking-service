#!/bin/bash
set -euo pipefail

AWS_ENDPOINT="--endpoint-url=http://localhost:4566"

echo "🚀 Configurando recursos AWS en LocalStack..."

# Verificar que LocalStack esté ejecutándose
echo "📦 Verificando LocalStack..."
if ! docker ps | grep -q localstack; then
    echo "❌ LocalStack no está ejecutándose. Iniciando..."
    docker-compose up -d
    echo "⏳ Esperando que LocalStack esté listo..."
    sleep 10
else
    echo "✅ LocalStack está ejecutándose"
fi

# Crear tabla DynamoDB solo si no existe
echo "🗄️ Configurando tabla DynamoDB..."
table_exists=$(aws $AWS_ENDPOINT dynamodb list-tables 2>/dev/null | grep 'tickets' || true)
if [ -z "$table_exists" ]; then
  echo "📝 Creando tabla DynamoDB 'tickets'..."
  aws $AWS_ENDPOINT dynamodb create-table \
    --table-name tickets \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5
  echo "✅ Tabla DynamoDB 'tickets' creada exitosamente"
else
  echo "✅ La tabla DynamoDB 'tickets' ya existe."
fi

# Crear bucket S3 solo si no existe
echo "☁️ Configurando bucket S3..."
# Intentar listar el bucket específico
if aws $AWS_ENDPOINT s3 ls s3://ticket-bucket 2>/dev/null; then
  echo "✅ El bucket S3 'ticket-bucket' ya existe."
else
  echo "📝 Creando bucket S3 'ticket-bucket'..."
  aws $AWS_ENDPOINT s3 mb s3://ticket-bucket
  echo "✅ Bucket S3 'ticket-bucket' creado exitosamente"
fi

# Crear cola SQS solo si no existe
echo "📬 Configurando cola SQS..."
queue_exists=$(aws $AWS_ENDPOINT sqs list-queues 2>/dev/null | grep 'ticket-queue' || true)
if [ -z "$queue_exists" ]; then
  echo "📝 Creando cola SQS 'ticket-queue'..."
  aws $AWS_ENDPOINT sqs create-queue --queue-name ticket-queue
  echo "✅ Cola SQS 'ticket-queue' creada exitosamente"
else
  echo "✅ La cola SQS 'ticket-queue' ya existe."
fi

# Verificar configuración
echo "🔍 Verificando configuración..."
echo "📊 Tablas DynamoDB:"
aws $AWS_ENDPOINT dynamodb list-tables 2>/dev/null || echo "❌ Error listando tablas DynamoDB"

echo "📊 Buckets S3:"
aws $AWS_ENDPOINT s3 ls 2>/dev/null || echo "❌ Error listando buckets S3"

echo "📊 Colas SQS:"
aws $AWS_ENDPOINT sqs list-queues 2>/dev/null || echo "❌ Error listando colas SQS"

echo "🎉 Configuración completada!"
echo "💡 Ahora puedes ejecutar: go run cmd/main.go"

#!/bin/bash
set -euo pipefail

AWS_ENDPOINT="--endpoint-url=http://localhost:4566"

# Crear tabla DynamoDB solo si no existe
table_exists=$(aws $AWS_ENDPOINT dynamodb list-tables | grep 'tickets' || true)
if [ -z "$table_exists" ]; then
  echo "Creando tabla DynamoDB 'tickets'..."
  aws $AWS_ENDPOINT dynamodb create-table \
    --table-name tickets \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5
else
  echo "La tabla DynamoDB 'tickets' ya existe."
fi

# Crear cola SQS solo si no existe
queue_exists=$(aws $AWS_ENDPOINT sqs list-queues | grep 'ticket-queue' || true)
if [ -z "$queue_exists" ]; then
  echo "Creando cola SQS 'ticket-queue'..."
  aws $AWS_ENDPOINT sqs create-queue --queue-name ticket-queue
else
  echo "La cola SQS 'ticket-queue' ya existe."
fi

# Crear bucket S3 solo si no existe
bucket_exists=$(aws $AWS_ENDPOINT s3api list-buckets | grep 'ticket-bucket' || true)
if [ -z "$bucket_exists" ]; then
  echo "Creando bucket S3 'ticket-bucket'..."
  aws $AWS_ENDPOINT s3 mb s3://ticket-bucket
else
  echo "El bucket S3 'ticket-bucket' ya existe."
fi

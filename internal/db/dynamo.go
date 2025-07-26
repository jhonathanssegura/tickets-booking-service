package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jhonathanssegura/ticket-reservation/internal/model"
)

type DynamoClient struct {
	Client *dynamodb.Client
}

func (d *DynamoClient) SaveTicket(ticket model.Ticket) error {
	av, err := attributevalue.MarshalMap(ticket)
	if err != nil {
		return err
	}

	_, err = d.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("tickets"),
		Item:      av,
	})
	return err
}

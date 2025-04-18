package awsclient

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

//go:generate mockery --name DynamoClient
type DynamoClient interface {
	GetItem(ctx context.Context, parameters *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, parameters *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

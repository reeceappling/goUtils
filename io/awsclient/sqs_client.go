package awsclient

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

var SqsClientKey = "sqs-client-key"

//go:generate mockery --name SqsClient
type SqsClient interface {
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
}

func GetSQSClient(ctx context.Context) (context.Context, SqsClient, error) {
	if existingClient, ok := ctx.Value(SqsClientKey).(SqsClient); ok {
		return ctx, existingClient, nil
	}
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, nil, err
	}
	client := sqs.NewFromConfig(cfg)
	return context.WithValue(ctx, SqsClientKey, client), client, nil
}

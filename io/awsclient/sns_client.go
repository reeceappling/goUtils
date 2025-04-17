package awsclient

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
)

var SnsClientKey = "sns-client-key"

//go:generate mockery --name SnsClient
type SnsClient interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

func GetSNSClient(ctx context.Context) (context.Context, SnsClient, error) {
	if existingClient, ok := ctx.Value(SnsClientKey).(SnsClient); ok {
		return ctx, existingClient, nil
	}
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, nil, err
	}
	client := sns.NewFromConfig(cfg)
	return context.WithValue(ctx, SnsClientKey, client), client, nil
}

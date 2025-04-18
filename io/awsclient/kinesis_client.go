package awsclient

import (
	"context"
	AwsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/firehose"
)

const FirehoseClientKey = "firehose-client-key"

type FirehoseClient interface {
	PutRecord(ctx context.Context, params *firehose.PutRecordInput, optFns ...func(*firehose.Options)) (*firehose.PutRecordOutput, error)
}

func GetFirehoseClient(ctx context.Context) (FirehoseClient, error) {
	if client, ok := ctx.Value(FirehoseClientKey).(FirehoseClient); ok {
		return client, nil
	}
	config, err := AwsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	return firehose.NewFromConfig(config), nil
}

func SetFirehoseClient(ctx context.Context, client FirehoseClient) context.Context {
	return context.WithValue(ctx, FirehoseClientKey, client)
}

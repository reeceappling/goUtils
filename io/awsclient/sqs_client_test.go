package awsclient

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockSqsClient struct {
}

func (mock mockSqsClient) ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
	return nil, errors.New("not implemented")
}
func (mock mockSqsClient) DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
	return nil, errors.New("not implemented")
}

func TestDefaultSqsClient(t *testing.T) {

	mockClient := mockSqsClient{}
	contextWithoutSqsClient := context.Background()
	contextWithSqsClient := context.WithValue(contextWithoutSqsClient, SqsClientKey, mockClient)

	t.Run("load already existing client", func(t *testing.T) {
		context, client, error := GetSQSClient(contextWithSqsClient)

		assert.Nil(t, error)
		assert.Equal(t, mockClient, client)
		assert.Equal(t, contextWithSqsClient, context)
	})
}

package awsclient

import (
	"context"
	"encoding/base64"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

const KmsClientKey = "kms-client-key"

type KmsClient interface {
	Decrypt(ctx context.Context, params *kms.DecryptInput, optFns ...func(*kms.Options)) (*kms.DecryptOutput, error)
}

func GetKMSClient(ctx context.Context) (context.Context, KmsClient, error) {
	if existingClient, ok := ctx.Value(KmsClientKey).(KmsClient); ok {
		return ctx, existingClient, nil
	}
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, nil, err
	}
	client := kms.NewFromConfig(cfg)
	return context.WithValue(ctx, KmsClientKey, client), client, nil
}

func Decrypt(ctx context.Context, encryptedString string) (decryptedString string, err error) {
	bs, err := base64.StdEncoding.DecodeString(encryptedString)
	_, kmsClient, err := GetKMSClient(ctx)
	if err != nil {
		return "", err
	}
	res, err := kmsClient.Decrypt(ctx, &kms.DecryptInput{
		CiphertextBlob: bs,
	})
	if err != nil {
		return "", err
	}
	return string(res.Plaintext), nil
}

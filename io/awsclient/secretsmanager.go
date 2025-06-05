package awsclient

import (
	"context"
	"errors"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/reeceappling/goUtils/v2/logging"
	"github.com/reeceappling/goUtils/v2/this"
	"github.com/reeceappling/goUtils/v2/utils"
	ctxUtils "github.com/reeceappling/goUtils/v2/utils/context"
	"os"
	"path"
)

//go:generate mockery --name SecretsManagerClient
type SecretsManagerClient interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

var secretsManagerClient SecretsManagerClient

func GetSecretsManagerClient(ctx context.Context) (SecretsManagerClient, error) {
	var err error
	if secretsManagerClient != nil {
		return secretsManagerClient, nil
	}
	secretsManagerClient, err = NewLocalFirstSecretsManagerClient(ctx)
	if err == nil {
		return secretsManagerClient, nil
	}
	secretsManagerClient, err = NewCloudSecretsManagerClient(ctx)
	return secretsManagerClient, err
}

func SetSecretsManagerClient(client SecretsManagerClient) {
	secretsManagerClient = client
}

func NewCloudSecretsManagerClient(ctx context.Context) (SecretsManagerClient, error) {
	secretsCfg, err := awsConfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	secretsManagerClient = secretsmanager.NewFromConfig(secretsCfg)
	return secretsManagerClient, nil
}

type CachingSecretsManagerClient struct {
	directory   string
	cloudClient SecretsManagerClient
}

func NewLocalFirstSecretsManagerClient(ctx context.Context) (CachingSecretsManagerClient, error) {
	directory, err := os.Getwd()
	if err != nil {
		directory = this.Dir()
	}

	environment := ctxUtils.GetStringFromContext(ctx, ctxUtils.Environment)

	if environment == "" {
		return CachingSecretsManagerClient{}, errors.New("environment not set")
	}

	cloudClient, err := NewCloudSecretsManagerClient(ctx)

	return CachingSecretsManagerClient{
		directory:   path.Join(directory, "persistence/secrets", environment),
		cloudClient: cloudClient,
	}, err
}

func (csm CachingSecretsManagerClient) GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	if params != nil && params.SecretId != nil {
		fileContents, err := os.ReadFile(path.Join(csm.directory, *params.SecretId))
		if err == nil {
			return &secretsmanager.GetSecretValueOutput{
				SecretBinary: fileContents,
				SecretString: utils.Pointer(string(fileContents)),
			}, nil
		}
	}
	output, err := csm.cloudClient.GetSecretValue(ctx, params, optFns...)
	if err != nil {
		return output, err
	}
	if params != nil && params.SecretId != nil {
		err = os.MkdirAll(path.Dir(path.Join(csm.directory, *params.SecretId)), 0777) //nolint:gosec
		if err != nil {
			logging.GetSugaredLogger(ctx).Errorw("Failed to cache secret value", "err", err)
		} else {
			err = os.WriteFile(path.Join(csm.directory, *params.SecretId), []byte(*output.SecretString), 0777) //nolint:gosec
			if err != nil {
				logging.GetSugaredLogger(ctx).Errorw("Failed to cache secret value", "err", err)
			}
		}
	}
	return output, nil
}

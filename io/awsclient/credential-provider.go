package awsclient

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"strings"
	"time"
)

var cfg aws.Config
var credentials aws.Credentials

func init() {
	var err error
	cfg, err = config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("default config load failure: " + err.Error())
	}

	credentials, err = refreshCredentials(context.TODO())
	if err != nil {
		panic("credentials refresh failure: " + err.Error())
	}
}

func refreshCredentials(ctx context.Context) (aws.Credentials, error) {
	var err error
	credentials, err = cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return aws.Credentials{}, err
	}

	if strings.HasPrefix(credentials.Source, "SharedConfigCredentials:") {
		credentials.CanExpire = true
		credentials.Expires = time.Now().Add(15 * time.Minute)
	}

	return credentials, nil
}

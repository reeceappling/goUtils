package awsclient

import (
	"context"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//go:generate mockery --name S3Client
type S3Client interface {
	ListObjectsV2(context.Context, *s3.ListObjectsV2Input, ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	GetObject(context.Context, *s3.GetObjectInput, ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	PutObject(context.Context, *s3.PutObjectInput, ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	DeleteObject(context.Context, *s3.DeleteObjectInput, ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	HeadObject(ctx context.Context, params *s3.HeadObjectInput, optFns ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
}

var s3Client S3Client

func GetS3Client() S3Client {
	if s3Client == nil {
		if err := SetupWithDefault(context.Background()); err != nil {
			panic(err)
		}
	}
	return s3Client
}

func SetS3Client(client S3Client) {
	s3Client = client
}

const AwsRegion = "us-east-1"

func SetupWithDefault(ctx context.Context) error {
	return Setup(ctx, defaultClientCfg, nil, func(options *s3.Options) {})
}

func Setup(ctx context.Context, options ClientConfig, customEndptResolver aws.EndpointResolverWithOptions, usePathStyle func(options2 *s3.Options)) error {
	if s3Client != nil {
		return nil
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(AwsRegion),
		config.WithHTTPClient(http.DefaultClient),
		config.WithEndpointResolverWithOptions(customEndptResolver), // TODO: ?????
		config.WithRetryer(func() aws.Retryer {
			return retry.NewStandard(func(options *retry.StandardOptions) {
				options.MaxAttempts = 10
				options.MaxBackoff = 500 * time.Millisecond
			})
		}),
	)

	if err != nil {
		return err
	}
	s3Client = NewS3Client(s3.NewFromConfig(cfg, usePathStyle))

	clientConfig = options
	return nil
}

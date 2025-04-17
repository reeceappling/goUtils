package awsclient

import (
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/reeceappling/goUtils/v2/logging"
	"github.com/reeceappling/goUtils/v2/noCommit_local"
	"github.com/reeceappling/goUtils/v2/this"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/reeceappling/goUtils/v2/errorreference"
)

type CloudS3Client struct {
	client *s3.Client
}

var _ S3Client = &CloudS3Client{}

func NewS3Client(client *s3.Client) S3Client {
	if local.IsLocal() {
		directory, err := os.Getwd()
		if err != nil {
			directory = this.Dir()
		}
		return NewLocalFirstS3Client(client, path.Join(directory, "persistence/s3"))
	}
	return NewCloudS3Client(client)
}

func NewCloudS3Client(client *s3.Client) CloudS3Client {
	return CloudS3Client{client: client}
}

func (adapter CloudS3Client) ListObjectsV2(
	ctx context.Context,
	inp *s3.ListObjectsV2Input,
	options ...func(*s3.Options),
) (*s3.ListObjectsV2Output, error) {
	response, err := adapter.client.ListObjectsV2(ctx, inp, options...)
	return response, StandardizeError(ctx, err)
}

func (adapter CloudS3Client) GetObject(
	ctx context.Context,
	input *s3.GetObjectInput,
	options ...func(*s3.Options),
) (*s3.GetObjectOutput, error) {
	response, err := adapter.client.GetObject(ctx, input, options...)
	return response, StandardizeError(ctx, err)
}

func (adapter CloudS3Client) PutObject(
	ctx context.Context,
	input *s3.PutObjectInput,
	options ...func(*s3.Options),
) (*s3.PutObjectOutput, error) {
	response, err := adapter.client.PutObject(ctx, input, options...)
	return response, StandardizeError(ctx, err)
}

func (adapter CloudS3Client) DeleteObject(
	ctx context.Context,
	input *s3.DeleteObjectInput,
	options ...func(*s3.Options),
) (*s3.DeleteObjectOutput, error) {
	response, err := adapter.client.DeleteObject(ctx, input, options...)
	return response, StandardizeError(ctx, err)
}

func (adapter CloudS3Client) HeadObject(
	ctx context.Context,
	input *s3.HeadObjectInput,
	options ...func(*s3.Options),
) (*s3.HeadObjectOutput, error) {
	response, err := adapter.client.HeadObject(ctx, input, options...)
	return response, StandardizeError(ctx, err)
}

func StandardizeError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}

	var errNoSuchKey *types.NoSuchKey
	var errNotFound *types.NotFound
	if errors.As(err, &errNoSuchKey) || errors.As(err, &errNotFound) {
		return errorreference.ErrorNotFound
	}

	var re s3.ResponseError
	if errors.As(err, &re) {
		log := logging.GetSugaredLogger(ctx)
		log.Errorw("aws error", "aws.hostId", re.ServiceHostID(), "aws.requestId", re.ServiceRequestID(), "err", re)
	}

	lowerErr := strings.ToLower(err.Error())

	if strings.Contains(lowerErr, "the specified bucket does not exist") {
		return ErrorUndefinedS3Bucket
	}

	result := retry.IsErrorThrottles(retry.DefaultThrottles).IsErrorThrottle(err)
	if result == aws.TrueTernary {
		return errorreference.ErrorSlowDown
	}

	// the preferred way to catch throttles doesn't appear to work very well,
	// so here's the backup net
	if strings.Contains(lowerErr, "retry quota exceeded") {
		return errorreference.ErrorSlowDown
	}
	if strings.Contains(lowerErr, "503") {
		return errorreference.ErrorSlowDown
	}

	if strings.Contains(lowerErr, "input member bucket must not be empty") {
		return ErrorUndefinedS3Bucket
	}

	if strings.Contains(lowerErr, "input member key must not be empty") {
		return ErrorUndefinedS3Key
	}

	return err
}

var (
	ErrorUndefinedS3Bucket = errors.New("bucket undefined")
	ErrorUndefinedS3Key    = errors.New("key undefined")
)

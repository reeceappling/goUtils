package s3

import (
	"bytes"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/reeceappling/goUtils/v2/errorreference"
	"github.com/reeceappling/goUtils/v2/io/awsclient"
)

type S3FileWriter struct {
	Bucket string
}

func (writer *S3FileWriter) Put(ctx context.Context, path string, data []byte) error {
	clientConfig := awsclient.GetClientConfig()
	client := awsclient.GetS3Client()
	var errs error
	for range clientConfig.MaxPutRetries {
		if _, err := client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: &writer.Bucket,
			Key:    &path,
			Body:   bytes.NewReader(data),
		}); err != nil {
			errs = errors.Join(errs, err)
		} else {
			return nil
		}
	}
	return errs
}

func (writer *S3FileWriter) Delete(ctx context.Context, path string) error {
	clientConfig := awsclient.GetClientConfig()
	client := awsclient.GetS3Client()
	var errs error
	for range clientConfig.MaxPutRetries {
		if _, err := client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: &writer.Bucket,
			Key:    &path,
		}); err != nil {
			if errors.Is(err, errorreference.ErrorNotFound) {
				return err
			}
			errs = errors.Join(errs, err)
		} else {
			return nil
		}
	}
	return errs
}

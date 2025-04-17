package s3

import (
	"context"
	"io"
)

//go:generate mockery --name BucketStreamer
type BucketStreamer interface {
	Construct(Bucket string)
	ReadStreaming(ctx context.Context, path string) (output io.ReadCloser, contentLength int64, err error)
}

type S3BucketStreamer struct {
	S3FileReader
}

func (streamer *S3BucketStreamer) Construct(Bucket string) {
	streamer.S3FileReader = S3FileReader{Bucket: Bucket}
}

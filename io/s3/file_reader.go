package s3

import (
	"context"
	"errors"
	"github.com/reeceappling/goUtils/v2/errorreference"
	"github.com/reeceappling/goUtils/v2/io/awsclient"
	recover2 "github.com/reeceappling/goUtils/v2/recover"
	"github.com/reeceappling/goUtils/v2/utils"
	"github.com/reeceappling/goUtils/v2/utils/channels"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	goio "io"
	"time"
)

type S3FileReader struct {
	Bucket string
}

// provide an unpopulated s3 file reader.
// bucket argument is required and will error later if missing.
func NewFileReader(args ...string) *S3FileReader {
	return &S3FileReader{Bucket: args[0]}
}

func (reader *S3FileReader) Read(ctx context.Context, path string) (output []byte, err error) {
	firstChan := backgroundRead(ctx, path, reader)

	select {
	case res := <-firstChan:
		return res.Item, res.Err
	case <-time.After(2 * time.Second): // slightly less arbitrary, basically everything completes before this
		return lazyRace(ctx, path, reader, firstChan)
	}
}

func lazyRace(ctx context.Context, path string, reader *S3FileReader, firstChan <-chan utils.ErrAnd[[]byte]) (output []byte, err error) {
	secondChan := backgroundRead(ctx, path, reader)
	select {
	case res := <-firstChan:
		if res.Err != nil {
			return handleBackupChan(res.Err, secondChan)
		}
		return res.Item, nil
	case res := <-secondChan:
		if res.Err != nil {
			return handleBackupChan(res.Err, firstChan)
		}
		return res.Item, nil
	}
}
func backgroundRead(ctx context.Context, path string, reader *S3FileReader) <-chan utils.ErrAnd[[]byte] {
	out := make(chan utils.ErrAnd[[]byte], 1)
	go func() {
		defer close(out)
		r, contentLength, err := reader.ReadStreaming(ctx, path)
		if err != nil {
			out <- utils.ErrAnd[[]byte]{Err: err}
			return
		}
		defer r.Close()
		output := make([]byte, contentLength)
		_, err = goio.ReadFull(r, output)
		out <- utils.ErrAnd[[]byte]{
			Item: output,
			Err:  err,
		}
	}()
	return out
}

func handleBackupChan(err error, c <-chan utils.ErrAnd[[]byte]) ([]byte, error) {
	other := <-c
	if other.Err != nil {
		return nil, errors.Join(err, other.Err)
	}
	return other.Item, nil
}

func (reader *S3FileReader) ReadStreaming(ctx context.Context, path string) (output goio.ReadCloser, contentLength int64, err error) {
	clientConfig := awsclient.GetClientConfig()
	client := awsclient.GetS3Client()

	bucket := reader.Bucket

	var res *s3.GetObjectOutput
	for i := 0; i < clientConfig.MaxReadRetries; i++ {
		res, err = client.GetObject(
			ctx,
			&s3.GetObjectInput{Bucket: &bucket, Key: &path},
		)

		if err == nil { // success. no other tests needed
			return res.Body, *res.ContentLength, nil
		}

		if errors.Is(err, context.Canceled) {
			return nil, 0, err // the request is aborted
		}

		if errors.Is(err, errorreference.ErrorNotFound) {
			return nil, 0, err
		}

		if !errors.Is(err, errorreference.ErrorSlowDown) {
			return nil, 0, err // unexpected/unhandled, catastrophic error
		} else {
			// TODO: log.Sugar().Warn("ErrorSlowDown from s3")
		}

		time.Sleep(utils.Jitter()) // retry after delay
	}

	return nil, 0, err
}

func (reader *S3FileReader) List(ctx context.Context, path string) (list []string, err error) {
	clientConfig := awsclient.GetClientConfig()
	client := awsclient.GetS3Client()

	for i := 0; i < clientConfig.MaxListRetries; i++ {
		list = []string{}
		paginator := s3.NewListObjectsV2Paginator(
			client,
			&s3.ListObjectsV2Input{Bucket: &reader.Bucket, Prefix: &path},
		)

		var page *s3.ListObjectsV2Output
		for paginator.HasMorePages() {
			page, err = paginator.NextPage(ctx)
			if err == nil { // success case
				for _, item := range page.Contents { // aggregate results
					list = append(list, *item.Key)
				}
				continue
			}

			if errors.Is(err, context.Canceled) {
				// TODO: log.Sugar().Debugw("Canceled S3 List")
				return // the request is aborted
			}

			if errors.Is(err, errorreference.ErrorSlowDown) {
				break // return to outer loop to try again
			} // else fallthrough

			// TODO: log.Sugar().Errorw("Unable to do S3 List", "error", err)
			return // unexpected/unhandled, catastrophic error
		}
		if err == nil {
			return // complete without error
		}

		time.Sleep(utils.Jitter()) // retry after delay
	}

	return
}

func (reader *S3FileReader) RaceRead(ctx context.Context, path string) ([]byte, error) {
	defaultConcurrentReads := 2
	return reader.RaceReadN(ctx, path, defaultConcurrentReads)
}

type s3Data struct {
	data []byte
	err  error
}

func (reader *S3FileReader) RaceReadN(ctx context.Context, path string, concurrentReads int) ([]byte, error) {
	ctx, cancel := context.WithCancel(ctx)

	chans := make([]<-chan s3Data, concurrentReads)
	for i := 0; i < concurrentReads; i++ {
		chans[i] = reader.readerProducer(ctx, path)
	}
	channel := channels.Multiplex(chans)
	defer func() {
		cancel()
		channels.Drain(channel)
	}()

	var err error
	for i := 0; i < concurrentReads; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case res := <-channel:
			if res.err == nil {
				return res.data, res.err
			}
			var nsk *types.NoSuchKey
			if errors.As(err, &nsk) {
				return nil, nsk
			}
			err = errors.Join(err, res.err)
		}
	}

	return nil, errors.Join(errors.New("all concurrent reads failed"), err)
}

func (reader *S3FileReader) readerProducer(ctx context.Context, path string) <-chan s3Data {
	output := make(chan s3Data)

	go func() {
		defer close(output)
		defer func() {
			if err := recover2.HandleRecoverAndLog(ctx, recover()); err != nil { // TODO: EW
				output <- s3Data{data: nil, err: err}
			}
		}()
		bytes, err := reader.Read(ctx, path)

		output <- s3Data{bytes, err}
	}()

	return output
}

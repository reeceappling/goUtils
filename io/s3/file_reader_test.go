package s3

import (
	"context"
	"errors"
	"github.com/reeceappling/goUtils/v2/io/awsclient/mocks"
	"github.com/reeceappling/goUtils/v2/utils"
	"github.com/stretchr/testify/mock"
	goio "io"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/reeceappling/goUtils/v2/errorreference"
	"github.com/reeceappling/goUtils/v2/io/awsclient"
	"github.com/reeceappling/goUtils/v2/logging"
	"github.com/stretchr/testify/assert"
)

// TODO: ALL TESTS

func isObjectPath(bucketName string, objectPath string) func(input *s3.GetObjectInput) bool {
	return func(input *s3.GetObjectInput) bool {
		return *input.Bucket == bucketName && *input.Key == objectPath
	}
}
func isListPath(bucketName string, prefix string, continuationToken *string) func(input *s3.ListObjectsV2Input) bool {
	return func(input *s3.ListObjectsV2Input) bool {
		return *input.Bucket == bucketName && *input.Prefix == prefix &&
			((continuationToken == nil && input.ContinuationToken == nil) ||
				(continuationToken != nil && input.ContinuationToken != nil && *input.ContinuationToken == *continuationToken))
	}
}

func getDefaultS3TestPath() string {
	return "path" // TODO: fix
}

func TestS3FileReader(t *testing.T) {
	log := logging.LoggerFactoryFor("test") // TODO: ok?
	ctx := logging.SetLogger(context.Background(), log)
	mockS3Client := mocks.NewS3Client(t)
	awsclient.SetS3Client(mockS3Client)
	//bucket := "s3-file-reader-test-bucket"

	t.Run("fails without bucket", func(t *testing.T) {
		mockS3Client.On("GetObject", mock.Anything, mock.MatchedBy(isObjectPath("", "path"))).Return(
			nil,
			awsclient.ErrorUndefinedS3Bucket,
		)

		_, err := (&S3FileReader{""}).Read(ctx, "path")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, awsclient.ErrorUndefinedS3Bucket))
	})

	t.Run("fails without path", func(t *testing.T) {
		mockS3Client.On("GetObject", mock.Anything, mock.MatchedBy(isObjectPath("bucket", ""))).Return(
			nil,
			awsclient.ErrorUndefinedS3Key,
		)

		_, err := (&S3FileReader{"bucket"}).Read(ctx, "")
		assert.Error(t, err)
		assert.True(t, errors.Is(err, awsclient.ErrorUndefinedS3Key))
	})

	//t.Run("reads a file from s3", func(t *testing.T) { // TODO: fix
	//	path := addDefaultS3TestAvroPath(addDefaultS3TestQuadKeyPath(getDefaultS3TestPath()))
	//	reader := strings.NewReader("test file")
	//	mockS3Client.On("GetObject", mock.Anything, mock.MatchedBy(isObjectPath(bucket, path))).Return(
	//		&s3.GetObjectOutput{Body: goio.NopCloser(reader), ContentLength: utils.Pointer(reader.Size())},
	//		nil,
	//	)
	//	data, err := (&S3FileReader{Bucket: bucket}).Read(ctx, path)
	//	assert.NoError(t, err)
	//	actual := sha256.Sum256(data)
	//	expected := "9a30a503b2862c51c3c5acd7fbce2f1f784cf4658ccf8e87d5023a90c21c0714"
	//	assert.Equal(t, expected, hex.EncodeToString(actual[:]))
	//})

	// only errors for permissions or systematic failures // TODO: FIX
	//t.Run("listing an invalid location will return an empty list", func(t *testing.T) {
	//	mockS3Client.On("ListObjectsV2", mock.Anything, mock.MatchedBy(isListPath(bucket, "invalid", nil))).Return(
	//		&s3.ListObjectsV2Output{Contents: []types.Object{}},
	//		nil,
	//	)
	//	data, err := (&S3FileReader{Bucket: bucket}).List(ctx, "invalid")
	//	assert.NoError(t, err)
	//	assert.Equal(t, []string{}, data)
	//})

	// TODO: fix
	//t.Run("listing a valid location produces a list of readers for found objects", func(t *testing.T) {
	//	root := addDefaultS3TestQuadKeyPath(getDefaultS3TestPath())
	//	listPath := root + "/list-test-path"
	//	mockS3Client.On("ListObjectsV2", mock.Anything, mock.MatchedBy(isListPath(bucket, listPath, nil))).Return(
	//		&s3.ListObjectsV2Output{Contents: []types.Object{
	//			{
	//				Key: utils.Pointer(listPath + "/subsession-batch-0.avro"),
	//			},
	//			{
	//				Key: utils.Pointer(listPath + "/subsession-batch-1.avro"),
	//			},
	//		}},
	//		nil,
	//	)
	//	data, err := (&S3FileReader{Bucket: bucket}).List(ctx, listPath)
	//	assert.NoError(t, err)
	//
	//	expected := []string{
	//		listPath + "/subsession-batch-0.avro",
	//		listPath + "/subsession-batch-1.avro",
	//	}
	//	assert.Equal(t, expected, data)
	//})

	// TODO: fix
	//t.Run("listing a large location will aggregate list over multiple pages", func(t *testing.T) {
	//	root := "rootDir" // TODO: fixme
	//
	//	s3ListPageSize := 1000
	//
	//	contents := make([]types.Object, s3ListPageSize+10)
	//	for i := range contents {
	//		contents[i].Key = utils.Pointer(root + strconv.Itoa(i) + ".txt")
	//	}
	//	mockS3Client.On("ListObjectsV2", mock.Anything, mock.MatchedBy(isListPath(bucket, root, nil))).Return(
	//		&s3.ListObjectsV2Output{Contents: contents[0:s3ListPageSize], NextContinuationToken: utils.Pointer("continue"), IsTruncated: aws.Bool(true)},
	//		nil,
	//	)
	//	mockS3Client.On("ListObjectsV2", mock.Anything, mock.MatchedBy(isListPath(bucket, root, utils.Pointer("continue")))).Return(
	//		&s3.ListObjectsV2Output{Contents: contents[s3ListPageSize:]},
	//		nil,
	//	)
	//
	//	data, err := (&S3FileReader{Bucket: bucket}).List(ctx, root)
	//	assert.NoError(t, err)
	//	assert.Equal(t, s3ListPageSize+10, len(data))
	//})
}

// AWS_PROFILE=SCUD_QUAL /usr/bin/time go test -benchmem -run=^$ -bench ^BenchmarkS3FileReader_Read -benchtime=100x
// Baseline                100         169888017 ns/op         7882115 Bytes/op        692 allocs/op
// Timeout 1 sec           100         174604980 ns/op         7930726 Bytes/op        673 allocs/op
//func BenchmarkS3FileReader_Read(b *testing.B) { // TODO: fix
//	logger := &logging.Logger{Logger: zap.NewNop()}
//	ctx := logging.SetLogger(context.Background(), logger)
//	bucket := "aws-bucket-name"
//	b.ResetTimer()
//	for i := 0; i < b.N; i++ {
//		data, err := (&S3FileReader{Bucket: bucket}).Read(ctx, addDefaultS3TestAvroPath(addDefaultS3TestQuadKeyPath(getDefaultS3TestPath()))) // TODO: fix
//		assert.NoError(b, err)
//		actual := sha256.Sum256(data)
//		expected := "a8b4377a469122a2b5a2ff6ac98a32eaa8ec6b87ff1f3d7cda8393c2323326d6"
//		assert.Equal(b, expected, hex.EncodeToString(actual[:]))
//	}
//}

// break out test for edge cases
func TestS3FileReaderRead(t *testing.T) {
	t.Run("repeated throttle errors will return slow down", func(t *testing.T) {
		callCount := 0
		client := &MockS3Client{
			MockGetObject: func(ctx2 context.Context, input *s3.GetObjectInput, f ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				callCount++
				return nil, errorreference.ErrorSlowDown
			},
		}
		awsclient.SetS3Client(client)
		clientConfig := awsclient.GetClientConfig()

		reader := NewFileReader("my-s3-bucket")

		_, err := reader.Read(context.Background(), "")
		assert.Error(t, err)
		// Allow for lazyRead to be triggered if S3 read is taking too long due to jitter and retries
		assert.True(t, clientConfig.MaxReadRetries == callCount || clientConfig.MaxReadRetries*2 == callCount)
		assert.True(t, errors.Is(err, errorreference.ErrorSlowDown))
	})

	t.Run("Should return a cancelled context error for reads if the parent context is cancelled", func(t *testing.T) {
		i := 0
		s3Client := &MockS3Client{
			MockGetObject: func(ctx2 context.Context, input *s3.GetObjectInput, f ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
				time.Sleep(time.Second)
				select {
				case <-ctx2.Done():
					return nil, context.Canceled
				default:
					i++
					reader := strings.NewReader("Hello World")
					return &s3.GetObjectOutput{Body: goio.NopCloser(reader), ContentLength: utils.Pointer(reader.Size())}, nil
				}
			},
		}
		awsclient.SetS3Client(s3Client)
		clientConfig := awsclient.GetClientConfig()
		clientConfig.MaxReadRetries = 1
		clientConfig.ReadTimeoutDuration = 10 * time.Second
		awsclient.SetClientConfig(clientConfig)

		reader := NewFileReader("my-s3-bucket")

		ctx, cancelFunc := context.WithCancel(context.Background())
		out := make(chan error)
		go func() {
			defer close(out)
			_, err := reader.Read(ctx, "")
			out <- err
		}()
		cancelFunc()
		err := <-out

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context canceled")
		assert.Equal(t, 0, i)
	})
}

// break out test for edge cases
func TestS3FileReaderList(t *testing.T) {
	t.Run("an unexpected error will abort", func(t *testing.T) {
		errorref := errors.New("apocalypse")
		callCount := 0
		client := &MockS3Client{
			MockListObjectsV2: func(ctx2 context.Context, input *s3.ListObjectsV2Input, f ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
				callCount++
				return nil, errorref
			},
		}
		awsclient.SetS3Client(client)

		reader := NewFileReader("my-s3-bucket")

		_, err := reader.List(context.Background(), "")
		assert.Error(t, err)
		assert.Equal(t, 1, callCount)
		assert.True(t, errors.Is(err, errorref))
	})

	t.Run("repeated throttle errors will return slow down error", func(t *testing.T) {
		callCount := 0
		client := &MockS3Client{
			MockListObjectsV2: func(ctx2 context.Context, input *s3.ListObjectsV2Input, f ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
				callCount++
				return nil, errorreference.ErrorSlowDown
			},
		}
		awsclient.SetS3Client(client)
		clientConfig := awsclient.GetClientConfig()

		reader := NewFileReader("my-s3-bucket")

		_, err := reader.List(context.Background(), "")
		assert.Error(t, err)
		assert.Equal(t, clientConfig.MaxListRetries, callCount)
		assert.True(t, errors.Is(err, errorreference.ErrorSlowDown))
	})

	t.Run("Should list multiple times if times out", func(t *testing.T) {
		i := 0
		client := &MockS3Client{
			MockListObjectsV2: func(ctx2 context.Context, input *s3.ListObjectsV2Input, f ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
				i++
				if i < 4 {
					return nil, errorreference.ErrorSlowDown
				}
				key := "/dirName/" //nolint:goconst
				return &s3.ListObjectsV2Output{IsTruncated: utils.Pointer(false), Contents: []types.Object{{Key: &key}}}, nil
			},
		}
		awsclient.SetS3Client(client)
		clientConfig := awsclient.GetClientConfig()
		clientConfig.MaxListRetries = 4
		clientConfig.ListTimeoutDuration = time.Nanosecond
		awsclient.SetClientConfig(clientConfig)

		reader := NewFileReader("my-s3-bucket")

		list, err := reader.List(context.Background(), "")
		assert.NoError(t, err)
		assert.Equal(t, 4, i)
		assert.Equal(t, 1, len(list))
		assert.Equal(t, "/dirName/", list[0])
	})

	t.Run("Should not list multiple times if the first call goes fast enough", func(t *testing.T) {
		i := 0
		client := &MockS3Client{
			MockListObjectsV2: func(ctx2 context.Context, input *s3.ListObjectsV2Input, f ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
				i++
				key := "/dirName/"
				return &s3.ListObjectsV2Output{IsTruncated: utils.Pointer(false), Contents: []types.Object{{Key: &key}}}, nil
			},
		}
		awsclient.SetS3Client(client)
		clientConfig := awsclient.GetClientConfig()
		clientConfig.MaxListRetries = 4
		clientConfig.ListTimeoutDuration = time.Second
		awsclient.SetClientConfig(clientConfig)

		reader := NewFileReader("my-s3-bucket")

		list, err := reader.List(context.Background(), "")
		assert.NoError(t, err)
		assert.Equal(t, 1, i)
		assert.Equal(t, 1, len(list))
		assert.Equal(t, "/dirName/", list[0])
	})

	t.Run("Should return a cancelled context error for lists if the parent context is cancelled", func(t *testing.T) {
		i := 0
		s3Client := &MockS3Client{
			MockListObjectsV2: func(ctx2 context.Context, input *s3.ListObjectsV2Input, f ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
				i++
				if i == 1 {
					return nil, context.Canceled
				}
				key := "/dirName/"
				return &s3.ListObjectsV2Output{IsTruncated: utils.Pointer(false), Contents: []types.Object{{Key: &key}}}, nil
			},
		}
		awsclient.SetS3Client(s3Client)
		clientConfig := awsclient.GetClientConfig()
		clientConfig.MaxReadRetries = 4
		clientConfig.ReadTimeoutDuration = 10 * time.Second
		awsclient.SetClientConfig(clientConfig)

		reader := NewFileReader("my-s3-bucket")

		ctx, cancelFunc := context.WithCancel(context.Background())
		cancelFunc()
		_, err := reader.List(ctx, "")

		assert.Error(t, err)
		assert.ErrorIs(t, err, context.Canceled)
		assert.Equal(t, 1, i)
	})
}

func TestRaceRead(t *testing.T) {
	s3Client := &MockS3Client{
		MockGetObject: func(ctx context.Context, input *s3.GetObjectInput, f ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
			// give it enough time to hit the other side of the select
			// as well as simulate a real network call
			// closing the Done channel in the context can happen asynchronously
			// according to https://pkg.go.dev/context#Context
			time.Sleep(time.Second)
			return nil, errors.New("other error")
		},
	}
	awsclient.SetS3Client(s3Client)
	clientConfig := awsclient.GetClientConfig()
	clientConfig.MaxReadRetries = 4
	clientConfig.ReadTimeoutDuration = 10 * time.Second

	t.Run("returns error from context if it has been cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		reader := NewFileReader("my-s3-bucket")

		_, err := reader.RaceRead(ctx, "")
		assert.ErrorIs(t, err, context.Canceled)
	})
}

type MockS3Client struct {
	MockListObjectsV2 func(context.Context, *s3.ListObjectsV2Input, ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	MockGetObject     func(context.Context, *s3.GetObjectInput, ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	MockPutObject     func(context.Context, *s3.PutObjectInput, ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	MockDeleteObject  func(context.Context, *s3.DeleteObjectInput, ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

func (client *MockS3Client) ListObjectsV2(ctx context.Context, input *s3.ListObjectsV2Input, options ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	if client.MockListObjectsV2 != nil {
		return client.MockListObjectsV2(ctx, input, options...)
	}
	return nil, nil
}

func (client *MockS3Client) GetObject(ctx context.Context, input *s3.GetObjectInput, options ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	if client.MockGetObject != nil {
		return client.MockGetObject(ctx, input, options...)
	}
	return nil, nil
}

func (client *MockS3Client) PutObject(ctx context.Context, input *s3.PutObjectInput, options ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	if client.MockPutObject != nil {
		return client.MockPutObject(ctx, input, options...)
	}
	return nil, nil
}

func (client *MockS3Client) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput, options ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	if client.MockDeleteObject != nil {
		return client.MockDeleteObject(ctx, input, options...)
	}
	return nil, nil
}

func (client *MockS3Client) HeadObject(ctx context.Context, input *s3.HeadObjectInput, options ...func(options2 *s3.Options)) (*s3.HeadObjectOutput, error) {
	return &s3.HeadObjectOutput{ContentLength: utils.Pointer(int64(2))}, nil
}

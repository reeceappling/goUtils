package test

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/reeceappling/goUtils/v2/io/awsclient"
	"github.com/reeceappling/goUtils/v2/utils"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/reeceappling/goUtils/v2/errorreference"
)

// BadReadCloser meets the criteria for an io.ReadCloser
type BadReadCloser struct {
	BadReader
}

func (rc BadReadCloser) Close() error {
	return nil
}

// BadReader meets the criteria for an io.Reader()
type BadReader struct{}

func (r BadReader) Read(b []byte) (n int, err error) {
	return 0, errors.New("mock read failure")
}

// FakeHttpWriter can be used in place of an http writer (only when the var line below is present)
type FakeHttpWriter struct {
	Headers    http.Header
	Bytes      []byte
	StatusCode int
}

var _ http.ResponseWriter = &FakeHttpWriter{} // Forces pointer to struct to implement interface

func NewFakeHttpWriter() *FakeHttpWriter {
	return &FakeHttpWriter{Bytes: []byte{}, StatusCode: 200, Headers: http.Header{}}
}

func (w *FakeHttpWriter) Header() http.Header {
	return w.Headers
}

func (w *FakeHttpWriter) Write(b []byte) (int, error) {
	w.Bytes = make([]byte, len(b))
	copy(w.Bytes, b) // Must be copied because writers should not retain their input
	return 0, nil
}

func (w *FakeHttpWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
}

type EmptyMemcachedClient struct{}

func (cache EmptyMemcachedClient) Get(key string) (item *memcache.Item, err error) {
	return nil, errors.New("empty cache miss")
}
func (cache EmptyMemcachedClient) Set(item *memcache.Item) error {
	return nil
}

func SetEmptyMemcachedClient(ctx context.Context) context.Context {
	return context.WithValue(ctx, awsclient.MemcachedClientKey, EmptyMemcachedClient{})
}

func SetLocalAWSClient(ctx context.Context) context.Context {
	return SetLocalAWSClientWithRedirects(ctx, map[string]string{})
}

func SetLocalAWSClientWithRedirects(ctx context.Context, bucketRedirect map[string]string) context.Context {
	awsclient.SetS3Client(&LocalS3Client{bucketRedirect})
	return SetEmptyMemcachedClient(ctx)
}

type LocalS3Client struct {
	BucketRedirect map[string]string
}

func (tClient LocalS3Client) getRedirect(bucketName string) string {
	if redirect := tClient.BucketRedirect[bucketName]; redirect != "" {
		return redirect
	}
	return bucketName
}

func (tClient LocalS3Client) ListObjectsV2(ctx context.Context, input *s3.ListObjectsV2Input, options ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	bucketName := tClient.getRedirect(*input.Bucket)
	prefixPath, walkPath := bucketName, bucketName
	if strings.Contains(prefixPath, "*") {
		return nil, errors.New("invalid bucket name pattern")
	}
	if input.Prefix != nil {
		prefixPath = path.Join(bucketName, *input.Prefix)
		walkPath = path.Join(bucketName, path.Dir(*input.Prefix))
		if strings.Contains(*input.Prefix, "*") {
			return nil, errors.New("invalid prefix pattern")
		}
	}
	var files []string
	err := filepath.WalkDir(walkPath, func(pathStr string, d fs.DirEntry, err error) error {
		if d.IsDir() || strings.HasPrefix(path.Base(pathStr), ".") {
			return nil
		}
		foundFiles, err := filepath.Glob(path.Join(path.Dir(pathStr), "*"))
		if err != nil {
			return err
		}
		for _, foundFile := range foundFiles {
			if !strings.HasPrefix(path.Base(foundFile), ".") {
				files = append(files, foundFile)
			}
		}
		return filepath.SkipDir
	})
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return nil, errorreference.ErrorNotFound
		}
		return nil, err
	}

	now := time.Now()
	fileNameObjects := make([]types.Object, 0, len(files))
	bucketFolder := bucketName
	if !strings.HasSuffix(bucketFolder, string(os.PathSeparator)) {
		bucketFolder += string(os.PathSeparator)
	}
	for _, file := range files {
		if !strings.HasPrefix(file, prefixPath) {
			continue
		}
		fileKey := strings.TrimPrefix(file, bucketFolder)
		fileNameObjects = append(fileNameObjects, types.Object{Key: &fileKey, LastModified: &now})
	}

	output := &s3.ListObjectsV2Output{
		Prefix:   input.Prefix,
		Contents: fileNameObjects,
	}
	return output, nil
}

func (tClient LocalS3Client) GetObject(ctx context.Context, input *s3.GetObjectInput, options ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	itemPath := path.Join(tClient.getRedirect(*input.Bucket), *input.Key)
	fileContents, err := os.ReadFile(itemPath) //nolint:gosec
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return nil, errorreference.ErrorNotFound
		}
		return nil, err
	}
	now := time.Now()
	output := &s3.GetObjectOutput{
		LastModified:  &now,
		Body:          io.NopCloser(bytes.NewReader(fileContents)),
		ContentLength: utils.Pointer(int64(len(fileContents))),
	}
	return output, nil
}

func (LocalS3Client) PutObject(context.Context, *s3.PutObjectInput, ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	return nil, errors.New("not implemented")
}

func (tClient LocalS3Client) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput, opts ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	return nil, errors.New("not implemented")
}

func (tClient LocalS3Client) HeadObject(ctx context.Context, input *s3.HeadObjectInput, opts ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	itemPath := path.Join(tClient.getRedirect(*input.Bucket), *input.Key)
	fi, err := os.Stat(itemPath)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return nil, errorreference.ErrorNotFound
		}
		return nil, err
	}
	contentLength := fi.Size()
	return &s3.HeadObjectOutput{ContentLength: &contentLength}, nil
}

type TestSettings struct {
	RunTests      bool
	RunAcceptance bool
}

func RegisterTests(containsTests bool, containsBenchmarks bool) TestSettings {
	runAcceptance := false

	if !flag.Parsed() {
		acceptance := flag.Bool("acceptance", false, "only performed for acceptance")
		flag.Parse()
		runAcceptance = *acceptance
	} else {
		acceptance := flag.Lookup("acceptance")
		runAcceptance = acceptance != nil && acceptance.Value.String() == "true"
	}

	settings := TestSettings{
		RunTests:      true,
		RunAcceptance: runAcceptance,
	}
	runPattern := flag.Lookup("test.run")
	if containsTests && runPattern != nil && runPattern.Value.String() != "^$" {
		return settings
	}

	benchPattern := flag.Lookup("test.bench")
	if containsBenchmarks && benchPattern != nil && benchPattern.Value.String() != "" && benchPattern.Value.String() != "^$" {
		settings.RunAcceptance = true
		return settings
	}

	settings.RunTests = false
	return settings
}

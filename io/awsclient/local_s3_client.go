package awsclient

import (
	"bytes"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/reeceappling/goUtils/v2/errorreference"
	"github.com/reeceappling/goUtils/v2/utils"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

type LocalS3Client struct {
	defaultDirectory        string
	bucketDirectoryRedirect map[string]string
}

var _ S3Client = LocalS3Client{}

func (lc LocalS3Client) getRedirect(bucketName string) string {
	if redirect := lc.bucketDirectoryRedirect[bucketName]; redirect != "" {
		return redirect
	}
	return path.Join(lc.defaultDirectory, bucketName)
}

func (lc LocalS3Client) ListObjectsV2(ctx context.Context, input *s3.ListObjectsV2Input, options ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	bucketName := lc.getRedirect(*input.Bucket)
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
		if err != nil {
			return err
		}
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
		if _, ok := err.(*os.PathError); ok {
			return nil, errorreference.ErrorNotFound
		}
		return nil, err
	}

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
		fileNameObjects = append(fileNameObjects, types.Object{Key: &fileKey, LastModified: utils.Pointer(time.Now())})
	}

	output := &s3.ListObjectsV2Output{
		Prefix:   input.Prefix,
		Contents: fileNameObjects,
	}
	return output, nil
}

func (lc LocalS3Client) GetObject(ctx context.Context, input *s3.GetObjectInput, options ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	itemPath := path.Join(lc.getRedirect(*input.Bucket), *input.Key)
	fileContents, err := os.ReadFile(itemPath) //nolint:gosec
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
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

func (lc LocalS3Client) PutObject(ctx context.Context, input *s3.PutObjectInput, opts ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	itemPath := path.Join(lc.getRedirect(*input.Bucket), *input.Key)
	contents, err := io.ReadAll(input.Body)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(path.Dir(itemPath), 0777) //nolint:gosec
	if err != nil {
		return nil, err
	}
	err = os.WriteFile(itemPath, contents, 0777) //nolint:gosec
	if err != nil {
		return nil, err
	}
	return &s3.PutObjectOutput{}, nil
}

func (lc LocalS3Client) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput, opts ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	itemPath := path.Join(lc.getRedirect(*input.Bucket), *input.Key)
	err := os.Remove(itemPath)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			return nil, errorreference.ErrorNotFound
		}
		return nil, err
	}
	return &s3.DeleteObjectOutput{}, nil
}

func (lc LocalS3Client) HeadObject(ctx context.Context, input *s3.HeadObjectInput, opts ...func(*s3.Options)) (*s3.HeadObjectOutput, error) {
	itemPath := path.Join(lc.getRedirect(*input.Bucket), *input.Key)
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

func NewLocalS3Client(directory string) LocalS3Client {
	return LocalS3Client{defaultDirectory: directory}
}

type LocalFirstS3Client struct {
	cloudClient CloudS3Client
	localClient LocalS3Client
}

var _ S3Client = LocalFirstS3Client{}

func NewLocalFirstS3Client(client *s3.Client, directory string) LocalFirstS3Client {
	return LocalFirstS3Client{
		cloudClient: NewCloudS3Client(client),
		localClient: NewLocalS3Client(directory),
	}
}

func (lfs3 LocalFirstS3Client) ListObjectsV2(
	ctx context.Context,
	input *s3.ListObjectsV2Input,
	options ...func(*s3.Options),
) (*s3.ListObjectsV2Output, error) {
	if response, err := lfs3.localClient.ListObjectsV2(ctx, input, options...); err == nil {
		return response, err
	}
	response, err := lfs3.cloudClient.ListObjectsV2(ctx, input, options...)
	return response, StandardizeError(ctx, err)
}

func (lfs3 LocalFirstS3Client) GetObject(
	ctx context.Context,
	input *s3.GetObjectInput,
	options ...func(*s3.Options),
) (*s3.GetObjectOutput, error) {
	if response, err := lfs3.localClient.GetObject(ctx, input, options...); err == nil {
		return response, err
	}
	response, err := lfs3.cloudClient.GetObject(ctx, input, options...)
	return response, StandardizeError(ctx, err)
}

func (lfs3 LocalFirstS3Client) PutObject(
	ctx context.Context,
	input *s3.PutObjectInput,
	options ...func(*s3.Options),
) (*s3.PutObjectOutput, error) {
	return lfs3.localClient.PutObject(ctx, input, options...)
}

func (lfs3 LocalFirstS3Client) DeleteObject(
	ctx context.Context,
	input *s3.DeleteObjectInput,
	options ...func(*s3.Options),
) (*s3.DeleteObjectOutput, error) {
	return lfs3.localClient.DeleteObject(ctx, input, options...)
}

func (lfs3 LocalFirstS3Client) HeadObject(
	ctx context.Context,
	input *s3.HeadObjectInput,
	options ...func(*s3.Options),
) (*s3.HeadObjectOutput, error) {
	if response, err := lfs3.localClient.HeadObject(ctx, input, options...); err == nil {
		return response, err
	}
	response, err := lfs3.cloudClient.HeadObject(ctx, input, options...)
	return response, StandardizeError(ctx, err)
}

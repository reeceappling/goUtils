package io

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
)

type FileWriter interface {
	Put(ctx context.Context, path string, data []byte) error
	Delete(ctx context.Context, path string) error
}

// Interface checks
var (
	_ FileWriter = inMemoryFileWriter{}
)

// inMemoryFileWriter stores files as a map of filePath to fileBytes
type inMemoryFileWriter map[string][]byte

func NewInMemoryFileWriter() inMemoryFileWriter {
	return map[string][]byte{}
}
func (fw inMemoryFileWriter) Put(ctx context.Context, path string, data []byte) error {
	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)
	fw[path] = dataCopy
	return nil
}
func (fw inMemoryFileWriter) Delete(ctx context.Context, path string) error {
	delete(fw, path)
	return nil
}

// GetZippedBytes processes all the writer's internally contained files into zipped bytes
func (fw inMemoryFileWriter) GetZippedBytes() ([]byte, error) {
	buff := bytes.NewBuffer([]byte{})
	err := fw.WriteZippedBytes(buff)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

// WriteZippedBytes skips the middleman of GetZippedBytes.
//
// WriteZippedBytes should not be used directly with an http.ResponseWriter.
//
// For an http.ResponseWriter, use GetZippedBytes instead, then pass the bytes to the http.ResponseWriter
// if no error occurred.
func (fw inMemoryFileWriter) WriteZippedBytes(w io.Writer) error {
	zipper := zip.NewWriter(w)
	defer zipper.Close()
	for filePath, data := range fw {
		fileZipper, err := zipper.Create(filePath)
		if err != nil {
			return err
		}
		_, err = io.Copy(fileZipper, bytes.NewReader(data))
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO: per-file iterator?

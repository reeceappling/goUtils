package io

import (
	"context"
)

//go:generate mockery --name FileReader
type FileReader interface {
	List(context.Context, string) ([]string, error)
	Read(context.Context, string) ([]byte, error)
	RaceRead(ctx context.Context, path string) ([]byte, error)
}

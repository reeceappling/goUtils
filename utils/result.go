package utils

import "errors"

var NotFound = errors.New("not found")

type Result[T any] struct {
	Item *T
	Err  error
}

func SuccessfulResult[T any](t T) Result[T] {
	return Result[T]{Item: &t}
}

func ErroredResult[T any](err error) Result[T] {
	return Result[T]{Err: err}
}

func ResultFrom[T any](t T, err error) Result[T] {
	if err == nil {
		return Result[T]{Item: &t}
	} else {
		return Result[T]{Err: err}
	}
}

type ErrAnd[T any] struct { // TODO come up with better name
	Item T
	Err  error
}

func ErrAndT[T any](t T) ErrAnd[T] {
	return ErrAnd[T]{
		Item: t,
		Err:  nil,
	}
}

func TandErr[T any](t T, err error) ErrAnd[T] {
	return ErrAnd[T]{
		Item: t,
		Err:  err,
	}
}

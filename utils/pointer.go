package utils

import (
	slices2 "github.com/reeceappling/goUtils/v2/utils/slices"
	"reflect"
)

func Pointer[T any](t T) *T {
	return &t
}

func Default[T any](t *T, defaultValue T) T {
	if t != nil {
		return *t
	}
	return defaultValue
}

func PtrContains[T comparable](t *T, value T) bool {
	if t == nil {
		return false
	}
	return *t == value
}

func CountNotNil(ptrs ...any) int {
	count := 0
	for _, ptr := range ptrs {
		if !reflect.ValueOf(ptr).IsNil() {
			count++
		}
	}
	return count
}

func OneIsDefined(ptrs ...any) bool {
	return CountNotNil(ptrs...) == 1
}

func AllAreNil(ptrs ...any) bool {
	return CountNotNil(ptrs...) == 0
}

func NoneAreNil(ptrs ...any) bool {
	return CountNotNil(ptrs...) == len(ptrs)
}

func NonNil[T any](ptrs []*T) []T {
	keepFunc := func(p *T) bool {
		return p != nil
	}
	return slices2.Map(slices2.Filter(ptrs, keepFunc), func(i *T) T {
		return *i
	})
}

func IsPointer[T any](v T) bool {
	return reflect.ValueOf(v).Kind() == reflect.Ptr
}

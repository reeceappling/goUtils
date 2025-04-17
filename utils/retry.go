package utils

import (
	"math/rand"
	"time"
)

func JitterFunc(min *int, maxMinusMin *int) func() time.Duration {
	minMs, maxExtraMs := Default(min, 500), Default(maxMinusMin, 1000)
	return func() time.Duration {
		return jitter(minMs, maxExtraMs)
	}
}
func Jitter() time.Duration {
	return jitter(500, 1000)
}
func jitter(min, maxExtra int) time.Duration {
	return time.Duration(min+rand.Intn(maxExtra)) * time.Millisecond //nolint:gosec
}
func Retry[T any](f func() (T, error)) (t T, err error) {
	return RetryTimes(10, f)
}

func RetryTimes[T any](numRetries int, f func() (T, error)) (t T, err error) {
	for i := 0; i < numRetries; i++ {
		t, err = f()
		if err != nil {
			time.Sleep(Jitter())
			continue
		}
		return
	}
	return
}

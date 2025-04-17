package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

type BasicallyAFloat interface {
	~float64
}

func AnyNanEqual(t *testing.T, expected, actual any) {
	if e, eOk := expected.(float64); eOk {
		if a, aOk := actual.(float64); aOk {
			NaNEqual(t, e, a)
			return
		}
	}
	assert.Equal(t, expected, actual)
}

func NaNEqual[T BasicallyAFloat](t *testing.T, expected, actual T) {
	if math.IsNaN(float64(expected)) {
		if math.IsNaN(float64(actual)) {
			return
		} else {
			t.Logf("floats were not equal: expected=%f actual=%f", float64(expected), float64(actual))
		}
	}
	assert.Equal(t, expected, actual)
}
func SliceNaNEqual[T BasicallyAFloat](t *testing.T, expected []T, actual []T) {
	require.Equal(t, len(expected), len(actual), fmt.Sprintf("lens not equal len(expected)=%d len(actual)=%d", len(expected), len(actual)))
	for i := range expected {
		NaNEqual(t, expected[i], actual[i])
	}
}

package httpUtils

import "testing"

func TestMultiWriter(t *testing.T) {
	t.Run("httpMultiResponseChildWriter", func(t *testing.T) {
		t.Run("Write", func(t *testing.T) {
			// TODO: THIS
		})
	})
	t.Run("*HttpMultiResponseWriter", func(t *testing.T) {
		t.Run("registerDuplicate", func(t *testing.T) {
			// TODO: THIS
		})
		t.Run("Header", func(t *testing.T) {
			// TODO: ?
		})
		t.Run("finalizeHeaders", func(t *testing.T) {
			// TODO: THIS
		})
		t.Run("Write", func(t *testing.T) {
			// TODO: THIS
		})
		t.Run("WriteHeader", func(t *testing.T) {
			// TODO: THIS
		})
	})
}

func TestMultiwriterPool(t *testing.T) {
	// TODO: DuplicateRequestPoolMiddleware
}

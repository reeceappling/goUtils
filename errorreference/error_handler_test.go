package errorreference

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	testErrorStr    = "testError"
	testBadStatCode = 999
)

func TestErrorHandler(t *testing.T) {
	t.Run("Error handler setError() works as intended", func(t *testing.T) {
		handler := &ErrorHandler{}
		handler.SetError(testErrorStr, testBadStatCode)
		assert.Nil(t, handler.Err, "error was never set and should be nil")
		assert.Equal(t, testErrorStr, handler.Message, "message should be correctly written")
		assert.Equal(t, testBadStatCode, handler.StatusCode, "status code should be correctly written")
	})
}

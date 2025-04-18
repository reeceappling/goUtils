package context

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO: hate this file

func TestGetStringFromContext(t *testing.T) {
	t.Run("given a key of a property in the context, returns the value", func(t *testing.T) {
		k, v := "someKey", "someValue"
		ctx := SetStringInContext(context.Background(), k, v)

		result := GetStringFromContext(ctx, k)

		assert.Equal(t, v, result)
	})

	t.Run("given a key of a property not in the context, returns an empty string", func(t *testing.T) {
		ctx := context.Background()

		result := GetStringFromContext(ctx, "not a thing that exists")

		assert.Equal(t, "", result)
	})

	t.Run("given a key of a property in the context, returns empty string if that property is not a string", func(t *testing.T) {
		k, v := "someKey", 123
		ctx := context.WithValue(context.Background(), k, v)

		result := GetStringFromContext(ctx, k)

		assert.Equal(t, "", result)
	})
}

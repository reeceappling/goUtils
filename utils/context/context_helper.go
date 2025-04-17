package context

import (
	"context"
	"reflect"
)

// Create a new type definition to protect against context namespace collisions
type contextKey string

// global constants of properties we add to context
const (
	ClientApiKey      = "clientApiKey"
	ApiKey            = "apiKey"
	OriginEndpoint    = "originEndpoint"
	AppNameContextKey = "clusterApplicationName"
	ParentRequestId   = "parentRequestId"
	Environment       = "environment"
)

func GetStringFromContext(ctx context.Context, key string) string {
	value := ""
	result := ctx.Value(contextKey(key))
	if v, ok := result.(string); ok {
		value = v
	}
	return value
}

func SetValueInContext(ctx context.Context, key string, value any) context.Context { // TODO: test?
	return context.WithValue(ctx, contextKey(key), value)
}

func SetStringInContext(ctx context.Context, key string, value string) context.Context {
	return SetValueInContext(ctx, key, value)
}

func GetValueFromContext[T any](ctx context.Context, key string) T { // TODO: TEST
	result := ctx.Value(contextKey(key))
	if v, ok := result.(T); ok {
		return v
	}
	// If nonexistent, create empty version
	outPtr := reflect.New(reflect.TypeFor[T]())
	outVal := outPtr.Elem() // TODO: panics?
	value, ok := outVal.Interface().(T)
	if !ok {
		// TODO: SOMETHING HERE
	}
	return value
}

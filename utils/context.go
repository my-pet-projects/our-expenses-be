package utils

import (
	"context"
)

type contextKey string

const (
	// ContextKeyCorrelationID is a correlation id context key.
	ContextKeyCorrelationID = contextKey("correlationId")
)

func (c contextKey) String() string {
	return string(c)
}

// SetContextStringValue puts a value into the context.
func SetContextStringValue(ctx context.Context, key contextKey, value string) context.Context {
	return context.WithValue(ctx, key, value)
}

// GetContextStringValue gets values from the context.
func GetContextStringValue(ctx context.Context, key contextKey) string {
	value, ok := ctx.Value(key).(string)
	if !ok {
		return ""
	}
	return value
}

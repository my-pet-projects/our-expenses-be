package utils

import (
	"context"
)

type contextKey string

var (
	contextKeyAuthtoken     = contextKey("auth-token")
	contextKeyAnother       = contextKey("another")
	ContextKeyCorrelationID = contextKey("correlationId")
)

func (c contextKey) String() string {
	return string(c)
}

func SetContextStringValue(ctx context.Context, key contextKey, value string) context.Context {
	return context.WithValue(ctx, key, value)
}

func GetContextStringValue(ctx context.Context, key contextKey) string {
	value, ok := ctx.Value(key).(string)
	if !ok {
		return ""
	}
	return value
}

package logger

import (
	"context"
	"log/slog"
)

type contextKey struct{}

func ContextWith(ctx context.Context, l *slog.Logger) context.Context {
	if l == nil {
		return ctx
	}

	return context.WithValue(ctx, contextKey{}, l)
}

func FromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(contextKey{}).(*slog.Logger); ok {
		return l
	}

	return slog.Default()
}

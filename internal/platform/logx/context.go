package logx

import (
	"context"
	"log/slog"
)

type ctxKey struct{}

// WithRequestID returns a new context carrying a request id.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, ctxKey{}, id)
}

// RequestID returns the request id stored in ctx, or "" if there is none.
func RequestID(ctx context.Context) string {
	id, _ := ctx.Value(ctxKey{}).(string)
	return id
}

// FromContext returns a logger enriched with context fields (request id, etc.).
func FromContext(ctx context.Context) *slog.Logger {
	l := slog.Default()
	if id := RequestID(ctx); id != "" {
		l = l.With("req_id", id)
	}
	return l
}

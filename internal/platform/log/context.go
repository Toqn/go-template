package log

import (
	"context"
	"log/slog"
)

type ctxKey string

const keyReqID ctxKey = "req_id"

// WithRequestID returns a new context carrying a request id.
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyReqID, id)
}

// FromContext returns a logger enriched with context fields (request id, etc.).
func FromContext(ctx context.Context) *slog.Logger {
	l := slog.Default()
	if v := ctx.Value(keyReqID); v != nil {
		l = l.With("req_id", v)
	}
	return l
}

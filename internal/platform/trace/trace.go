// Package trace is a placeholder for distributed tracing.
package trace

import "context"

// Start begins a span with the given name and returns the span context and
// a function that ends it. It is a no-op for now; the signature is shaped for
// future OpenTelemetry wiring, so keep it stable.
func Start(ctx context.Context, _ string) (context.Context, func(err error)) {
	return ctx, func(error) {}
}

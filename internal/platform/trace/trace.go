package trace

import "context"

// Start returns a context carrying a span id and an end function.
// No-op for now; shape is compatible with future OpenTelemetry wiring.
func Start(ctx context.Context, name string) (context.Context, func(err error)) {
	// cheap placeholder; keep signature stable
	end := func(err error) { _ = err }
	return ctx, end
}

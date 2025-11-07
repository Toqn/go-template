package log

import (
	"log/slog"
	"os"
	"strings"
	"time"
)

// New returns a structured JSON logger with level from string.
func New(level string) *slog.Logger {
	var lv slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lv = slog.LevelDebug
	case "warn":
		lv = slog.LevelWarn
	case "error":
		lv = slog.LevelError
	default:
		lv = slog.LevelInfo
	}
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     lv,
		AddSource: false,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// shorten time to RFC3339 without monotonic junk
			if a.Key == slog.TimeKey {
				if t := a.Value.Time(); !t.IsZero() {
					a.Value = slog.StringValue(t.UTC().Format(time.RFC3339Nano))
				}
			}
			return a
		},
	})
	return slog.New(h)
}

// SetDefault wires this logger as the process-wide default.
func SetDefault(l *slog.Logger) { slog.SetDefault(l) }

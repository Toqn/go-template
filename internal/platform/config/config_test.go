package config

import (
	"testing"
	"time"
)

// clearEnv resets every variable Load reads so the host environment cannot
// leak into a test. t.Setenv restores the original values afterwards.
func clearEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{"APP_ENV", "HTTP_ADDR", "LOG_LEVEL", "SHUTDOWN_TIMEOUT"} {
		t.Setenv(k, "")
	}
}

func TestLoadDefaults(t *testing.T) {
	clearEnv(t)
	c, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := Config{Stage: Dev, HTTPAddr: ":8080", LogLevel: "info", ShutdownTimeout: 5 * time.Second}
	if c != want {
		t.Fatalf("got %+v, want %+v", c, want)
	}
}

func TestLoadOverrides(t *testing.T) {
	clearEnv(t)
	t.Setenv("APP_ENV", "prod")
	t.Setenv("HTTP_ADDR", "127.0.0.1:9000")
	t.Setenv("LOG_LEVEL", "DEBUG")
	t.Setenv("SHUTDOWN_TIMEOUT", "30s")
	c, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := Config{Stage: Prod, HTTPAddr: "127.0.0.1:9000", LogLevel: "debug", ShutdownTimeout: 30 * time.Second}
	if c != want {
		t.Fatalf("got %+v, want %+v", c, want)
	}
}

func TestLoadInvalid(t *testing.T) {
	tests := map[string]struct{ key, val string }{
		"stage":             {"APP_ENV", "staging"},
		"log level":         {"LOG_LEVEL", "verbose"},
		"shutdown garbage":  {"SHUTDOWN_TIMEOUT", "soon"},
		"shutdown negative": {"SHUTDOWN_TIMEOUT", "-1s"},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			clearEnv(t)
			t.Setenv(tc.key, tc.val)
			if _, err := Load(); err == nil {
				t.Fatalf("expected error for %s=%q", tc.key, tc.val)
			}
		})
	}
}

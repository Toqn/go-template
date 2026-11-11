// Package config loads and validates process configuration from the environment.
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

// Stage is the deployment environment the process runs in.
type Stage string

// Supported stages.
const (
	Prod Stage = "prod"
	Test Stage = "test"
	Dev  Stage = "dev"
)

// Config is the process configuration.
type Config struct {
	Stage           Stage
	HTTPAddr        string        // e.g., ":8080"
	LogLevel        string        // "debug"|"info"|"warn"|"error"
	ShutdownTimeout time.Duration // graceful shutdown deadline

	// Add adapter-specific settings below
	// e.g., for a database, message broker, external API, etc.
}

// Load reads the configuration from environment variables, applying defaults,
// and validates it.
func Load() (Config, error) {
	shutdown, err := time.ParseDuration(getenv("SHUTDOWN_TIMEOUT", "5s"))
	if err != nil {
		return Config{}, fmt.Errorf("SHUTDOWN_TIMEOUT: %w", err)
	}
	c := Config{
		Stage:           Stage(getenv("APP_ENV", "dev")),
		HTTPAddr:        getenv("HTTP_ADDR", ":8080"),
		LogLevel:        strings.ToLower(getenv("LOG_LEVEL", "info")),
		ShutdownTimeout: shutdown,
	}
	return c, c.Validate()
}

// Validate reports the first invalid setting, if any.
func (c Config) Validate() error {
	if c.Stage != Prod && c.Stage != Test && c.Stage != Dev {
		return errors.New("APP_ENV must be dev|test|prod")
	}
	switch c.LogLevel {
	case "debug", "info", "warn", "error":
	default:
		return errors.New("LOG_LEVEL must be debug|info|warn|error")
	}
	if c.HTTPAddr == "" {
		return errors.New("HTTP_ADDR required")
	}
	if c.ShutdownTimeout <= 0 {
		return errors.New("SHUTDOWN_TIMEOUT must be positive")
	}
	// force adapter settings when those adapters are wired.
	return nil
}

func getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

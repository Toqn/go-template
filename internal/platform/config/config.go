package config

import (
	"errors"
	"os"
	"time"
)

type Stage string

const (
	Prod Stage = "prod"
	Test Stage = "test"
	Dev  Stage = "dev"
)

type Config struct {
	Stage           Stage
	HTTPAddr        string        // e.g., ":8080"
	LogLevel        string        // "debug"|"info"|"warn"|"error"
	ShutdownTimeout time.Duration // graceful shutdown deadline

	// Add adapter-specific settings below
	// e.g., for a database, message broker, external API, etc.
}

func Load() (Config, error) {
	c := Config{
		Stage:    Stage(getenv("APP_ENV", "dev")),
		HTTPAddr: getenv("HTTP_ADDR", ":8080"),
		LogLevel: getenv("LOG_LEVEL", "info"),
	}
	return c, c.Validate()
}

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

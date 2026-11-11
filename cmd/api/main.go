// Command api is the HTTP entrypoint and composition root of the service.
package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	v1 "github.com/toqn/go-template/api/v1"
	"github.com/toqn/go-template/internal/platform/config"
	"github.com/toqn/go-template/internal/platform/logx"
	"github.com/toqn/go-template/internal/platform/trace"
)

// version is set at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	if err := run(); err != nil {
		slog.Error("fatal", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	logger := logx.New(cfg.LogLevel)
	logx.SetDefault(logger)

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           requestLogger(routes()),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		logger.Info("starting server", "addr", srv.Addr, "stage", cfg.Stage, "version", version)
		errCh <- srv.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		// ListenAndServe only returns early on failure (e.g. port in use).
		return fmt.Errorf("listen: %w", err)
	case <-ctx.Done():
	}
	stop() // a second signal now kills the process immediately

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown: %w", err)
	}
	logger.Info("server stopped")
	return nil
}

func routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthz)
	return mux
}

func healthz(w http.ResponseWriter, r *http.Request) {
	_, end := trace.Start(r.Context(), "healthz")
	defer end(nil)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v1.Health{Status: "ok"})
}

const maxReqIDLen = 128

// validReqID reports whether a client-supplied request id is safe to reuse.
func validReqID(id string) bool {
	if id == "" || len(id) > maxReqIDLen {
		return false
	}
	for _, c := range id {
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9':
		case c == '-', c == '_', c == '.', c == ':':
		default:
			return false
		}
	}
	return true
}

// requestLogger assigns each request an id (echoed in X-Request-ID) and logs
// method, path, status and duration once the request completes. Health checks
// are logged at debug level so probes don't flood the logs.
func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqID := r.Header.Get("X-Request-ID")
		if !validReqID(reqID) {
			reqID = rand.Text()
		}
		w.Header().Set("X-Request-ID", reqID)

		ctx := logx.WithRequestID(r.Context(), reqID)
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r.WithContext(ctx))

		level := slog.LevelInfo
		if r.URL.Path == "/healthz" {
			level = slog.LevelDebug
		}
		logx.FromContext(ctx).Log(ctx, level, "request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", float64(time.Since(start).Microseconds())/1000,
		)
	})
}

// statusRecorder captures the response status code for logging.
type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (s *statusRecorder) WriteHeader(code int) {
	if !s.wroteHeader {
		s.status = code
		s.wroteHeader = true
	}
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusRecorder) Write(b []byte) (int, error) {
	s.wroteHeader = true
	return s.ResponseWriter.Write(b)
}

// Unwrap lets http.ResponseController reach the underlying writer.
func (s *statusRecorder) Unwrap() http.ResponseWriter { return s.ResponseWriter }

package main

import (
	"context"
	"log"
	"log/slog"

	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/toqn/phito-backend/internal/platform/config"
	plog "github.com/toqn/phito-backend/internal/platform/log"
	"github.com/toqn/phito-backend/internal/platform/trace"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	logger := plog.New(cfg.LogLevel)
	plog.SetDefault(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		ctx, end := trace.Start(r.Context(), "healthz")
		defer end(nil)
		plog.FromContext(ctx).Info("healthz check")
		w.WriteHeader(http.StatusOK)
	})

	muxWithMiddleware := requestLogger(mux)

	srv := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      muxWithMiddleware,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		logger.Info("starting server", "addr", srv.Addr, "stage", cfg.Stage)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("could not start server", "error", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	<-ctx.Done()

	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("could not gracefully shutdown server", "error", err)
	} else {
		logger.Info("server stopped")
	}
}

func newReqID() string {
	// TODO keep trivial; swap later for ULIDs
	return time.Now().UTC().Format("20060102T150405.000000000Z07:00")
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = newReqID()
		}
		ctx := plog.WithRequestID(r.Context(), reqID)
		slog.Info("request_start", "method", r.Method, "path", r.URL.Path, "req_id", reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

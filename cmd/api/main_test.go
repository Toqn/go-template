package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	v1 "github.com/toqn/go-template/api/v1"
	"github.com/toqn/go-template/internal/platform/logx"
)

func TestHealthz(t *testing.T) {
	rec := httptest.NewRecorder()
	requestLogger(routes()).ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q", ct)
	}
	var body v1.Health
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body.Status != "ok" {
		t.Fatalf("status field = %q, want ok", body.Status)
	}
}

func TestHealthzRejectsPost(t *testing.T) {
	rec := httptest.NewRecorder()
	routes().ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/healthz", nil))
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestRequestLoggerRequestID(t *testing.T) {
	tests := map[string]struct {
		incoming string
		reuse    bool
	}{
		"missing":   {"", false},
		"valid":     {"abc-123_x.y:z", true},
		"too long":  {strings.Repeat("a", maxReqIDLen+1), false},
		"bad chars": {"abc\ndef", false},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var seen string
			h := requestLogger(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				seen = logx.RequestID(r.Context())
			}))
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tc.incoming != "" {
				req.Header.Set("X-Request-ID", tc.incoming)
			}
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)

			got := rec.Header().Get("X-Request-ID")
			if got == "" || got != seen {
				t.Fatalf("response id %q, context id %q", got, seen)
			}
			if tc.reuse != (got == tc.incoming) {
				t.Fatalf("incoming %q, got %q, want reuse=%v", tc.incoming, got, tc.reuse)
			}
		})
	}
}

func TestStatusRecorder(t *testing.T) {
	rec := &statusRecorder{ResponseWriter: httptest.NewRecorder(), status: http.StatusOK}
	rec.WriteHeader(http.StatusTeapot)
	rec.WriteHeader(http.StatusInternalServerError) // superfluous; net/http ignores it
	if rec.status != http.StatusTeapot {
		t.Fatalf("status = %d, want %d", rec.status, http.StatusTeapot)
	}
}

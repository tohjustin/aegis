package server

import (
	"encoding/json"
	"net/http"
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/tohjustin/aegis/internal/observability"
	"github.com/tohjustin/aegis/internal/version"
)

func newAdminHandler(metrics *observability.Metrics, ready *atomic.Bool) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /livez", func(w http.ResponseWriter, r *http.Request) {
		writePlain(w, http.StatusOK, "ok\n")
	})
	mux.HandleFunc("GET /readyz", func(w http.ResponseWriter, r *http.Request) {
		if !ready.Load() {
			writePlain(w, http.StatusServiceUnavailable, "not ready\n")
			return
		}
		writePlain(w, http.StatusOK, "ok\n")
	})
	mux.Handle("GET /metrics", promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{}))
	mux.HandleFunc("GET /version", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, version.Current())
	})
	return mux
}

func writePlain(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	if status >= 500 {
		msg = "internal error"
	}
	writeJSON(w, status, map[string]string{"error": msg})
}

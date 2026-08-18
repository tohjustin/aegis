package observability

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strconv"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusRecorder) Write(b []byte) (int, error) {
	n, err := s.ResponseWriter.Write(b)
	s.bytes += n
	return n, err
}

func (s *statusRecorder) Unwrap() http.ResponseWriter {
	return s.ResponseWriter
}

// Instrument logs access records and records HTTP metrics. Route labels use the
// ServeMux pattern, or "unmatched" when no pattern was selected.
func Instrument(logger *slog.Logger, metrics *Metrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			start := time.Now()
			next.ServeHTTP(rec, r)

			route := r.Pattern
			if route == "" {
				route = "unmatched"
			}
			duration := time.Since(start)
			code := strconv.Itoa(rec.status)
			metrics.Requests.WithLabelValues(r.Method, route, code).Inc()
			metrics.Duration.WithLabelValues(r.Method, route).Observe(duration.Seconds())

			logger.LogAttrs(r.Context(), slog.LevelInfo, "request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("route", route),
				slog.Int("status", rec.status),
				slog.Int("bytes", rec.bytes),
				slog.Duration("duration", duration),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
			)
		})
	}
}

// Recover is inner middleware so Instrument still records a 500 after a panic.
func Recover(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				v := recover()
				if v == nil {
					return
				}
				if v == http.ErrAbortHandler {
					panic(v)
				}
				logger.LogAttrs(r.Context(), slog.LevelError, "panic recovered",
					slog.Any("panic", v),
					slog.String("stack", string(debug.Stack())),
				)
				writeJSONError(w, http.StatusInternalServerError, "internal error")
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

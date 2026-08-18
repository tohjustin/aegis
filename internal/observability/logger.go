// Package observability configures structured logging and Prometheus metrics.
package observability

import (
	"io"
	"log/slog"
	"strings"
)

// NewLogger constructs a slog logger and installs it as the process default.
func NewLogger(level, format string, w io.Writer) (*slog.Logger, error) {
	var lvl slog.Level
	if err := lvl.UnmarshalText([]byte(strings.ToLower(level))); err != nil {
		return nil, err
	}
	opts := &slog.HandlerOptions{Level: lvl}

	var handler slog.Handler
	switch strings.ToLower(format) {
	case "text":
		handler = slog.NewTextHandler(w, opts)
	default:
		handler = slog.NewJSONHandler(w, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger, nil
}

package server

import (
	"log/slog"
	"net/http"

	"github.com/tohjustin/aegis/internal/config"
	"github.com/tohjustin/aegis/pkg/badge"
)

func (a *app) writeBadge(w http.ResponseWriter, r *http.Request, params *badge.Params) {
	generated, err := badge.Create(params)
	if err != nil {
		a.logger.LogAttrs(r.Context(), slog.LevelError, "failed to create badge",
			slog.String("url", r.URL.RequestURI()),
			slog.Any("err", err),
		)
		if err := a.writeErrorBadge(w, "internal server error"); err != nil {
			a.logger.LogAttrs(r.Context(), slog.LevelError, "failed to create error badge",
				slog.String("url", r.URL.RequestURI()),
				slog.Any("err", err),
			)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	setBadgeHeaders(w, a.cfg)
	if _, err := w.Write([]byte(generated)); err != nil {
		a.logger.LogAttrs(r.Context(), slog.LevelError, "failed to write HTTP response",
			slog.String("url", r.URL.RequestURI()),
			slog.Any("err", err),
		)
	}
}

func (a *app) writeErrorBadge(w http.ResponseWriter, status string) error {
	generated, err := badge.Create(&badge.Params{
		Subject: "aegis",
		Status:  status,
	})
	if err != nil {
		return err
	}
	setBadgeHeaders(w, a.cfg)
	_, err = w.Write([]byte(generated))
	return err
}

func setBadgeHeaders(w http.ResponseWriter, cfg *config.Config) {
	if !cfg.ExcludeCacheControlHeaders {
		w.Header().Set("Cache-Control", "public, max-age=3600, s-maxage=3600")
	}
	w.Header().Set("Content-Type", "image/svg+xml;utf-8")
}

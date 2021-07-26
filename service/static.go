package service

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/tohjustin/aegis/pkg/badge"
	"github.com/tohjustin/aegis/service/config"
)

type staticService struct {
	name   string
	config *config.Config
	logger *zap.Logger
}

// NewStaticService returns a HTTP handler for the static badge service
func NewStaticService(configuration *config.Config,
	logger *zap.Logger) (BadgeService, error) {
	if configuration == nil {
		return nil, fmt.Errorf("missing config dependency")
	}
	if logger == nil {
		return nil, fmt.Errorf("missing logger dependency")
	}

	return &staticService{
		name:   "static",
		config: configuration,
		logger: logger,
	}, nil
}

func (service *staticService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	generatedBadge, err := badge.Create(&badge.Params{
		Style:   badge.Style(r.URL.Query().Get("style")),
		Subject: r.URL.Query().Get("subject"),
		Status:  r.URL.Query().Get("status"),
		Color:   r.URL.Query().Get("color"),
		Icon:    r.URL.Query().Get("icon"),
	})
	if err != nil {
		service.logger.Error("Failed to create badge",
			zap.String("url", r.URL.RequestURI()),
			zap.String("service", service.name),
			zap.Error(err))
		if err := internalServerError(w, service.config); err != nil {
			service.logger.Error("Failed to create error badge",
				zap.String("url", r.URL.RequestURI()),
				zap.String("service", service.name),
				zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if !service.config.ExcludeCacheControlHeaders {
		// cache response in browser for 1 hour (3600), CDN for 1 hour (3600)
		w.Header().Set("Cache-Control", "public, max-age=3600, s-maxage=3600")
	}
	w.Header().Set("Content-Type", "image/svg+xml;utf-8")
	_, err = w.Write([]byte(generatedBadge))
	service.logger.Error("Failed to write HTTP response",
		zap.String("url", r.URL.RequestURI()),
		zap.String("service", service.name),
		zap.Error(err))
}

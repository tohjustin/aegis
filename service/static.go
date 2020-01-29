package service

import (
	"fmt"
	"net/http"

	"go.uber.org/zap"

	"github.com/tohjustin/badger/pkg/badge"
	"github.com/tohjustin/badger/service/config"
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
			zap.String("service", service.name),
			zap.Error(err))
		internalServerError(w)
		return
	}

	// cache response in browser for 1 day (86400), CDN for 1 year (31536000)
	w.Header().Set("Cache-Control", "public, immutable, max-age=86400, s-maxage=31536000")
	w.Header().Set("Content-Type", "image/svg+xml;utf-8")
	w.Write([]byte(generatedBadge))
}

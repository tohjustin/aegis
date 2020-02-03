package service

import (
	"net/http"

	"github.com/tohjustin/aegis/pkg/badge"
	"github.com/tohjustin/aegis/service/config"
)

func generateErrorBadge(w http.ResponseWriter,
	configuration *config.Config, status string) error {
	generatedBadge, err := badge.Create(&badge.Params{
		Subject: "aegis",
		Status:  status,
	})
	if err != nil {
		return err
	}

	if !configuration.ExcludeCacheControlHeaders {
		// cache response in browser for 1 hour (3600), CDN for 1 hour (3600)
		w.Header().Set("Cache-Control", "public, max-age=3600, s-maxage=3600")
	}
	w.Header().Set("Content-Type", "image/svg+xml;utf-8")
	w.Write([]byte(generatedBadge))
	return nil
}

// badRequest handles HTTP requests that are malformed
func badRequest(w http.ResponseWriter,
	configuration *config.Config) error {
	return generateErrorBadge(w, configuration, "bad request")
}

// internalServerError handles HTTP requests that results in internal server error
func internalServerError(w http.ResponseWriter,
	configuration *config.Config) error {
	return generateErrorBadge(w, configuration, "internal server error")
}

// notFound handles HTTP requests for methods that don't exist
func notFound(w http.ResponseWriter,
	configuration *config.Config) error {
	return generateErrorBadge(w, configuration, "not found")
}

// serviceNotFound handles HTTP requests for services that don't exist
func serviceNotFound(w http.ResponseWriter,
	configuration *config.Config) error {
	return generateErrorBadge(w, configuration, "service not found")
}

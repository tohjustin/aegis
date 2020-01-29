package service

import (
	"net/http"

	"github.com/tohjustin/badger/pkg/badge"
	"github.com/tohjustin/badger/service/config"
)

func generateErrorBadge(w http.ResponseWriter,
	configuration *config.Config, status string) error {
	generatedBadge, err := badge.Create(&badge.Params{
		Subject: "badger",
		Status:  status,
	})
	if err != nil {
		return err
	}

	if !configuration.ExcludeCacheControlHeaders {
		// cache response in browser for 1 day (86400), CDN for 1 year (31536000)
		w.Header().Set("Cache-Control", "public, immutable, max-age=86400, s-maxage=31536000")
	}
	w.Header().Set("Content-Type", "image/svg+xml;utf-8")
	w.Write([]byte(generatedBadge))
	return nil
}

// badRequest handles HTTP requests that are malformed
func badRequest(w http.ResponseWriter,
	configuration *config.Config) error {
	return generateErrorBadge(w, configuration, "400 Bad Request")
}

// internalServerError handles HTTP requests that results in internal server error
func internalServerError(w http.ResponseWriter,
	configuration *config.Config) error {
	return generateErrorBadge(w, configuration, "500 Internal Server Error")
}

// notFound handles HTTP requests for services or methods that don't exist
func notFound(w http.ResponseWriter,
	configuration *config.Config) error {
	return generateErrorBadge(w, configuration, "404 Not Found")
}

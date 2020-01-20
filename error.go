package main

import (
	"fmt"
	"net/http"

	"github.com/tohjustin/badger/pkg/badge"
)

func generateErrorBadge(w http.ResponseWriter, status string) {
	generatedBadge, err := badge.Create(&badge.Params{
		Subject: "badger",
		Status:  status,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// cache response in browser for 1 day (86400), CDN for 1 year (31536000)
	w.Header().Set("Cache-Control", "public, immutable, max-age=86400, s-maxage=31536000")
	w.Header().Set("Content-Type", "image/svg+xml;utf-8")
	w.Write([]byte(generatedBadge))
}

// badRequest handles HTTP requests that are malformed
func badRequest(w http.ResponseWriter) {
	generateErrorBadge(w, "400 Bad Request")
}

// internalServerError handles HTTP requests that results in internal server error
func internalServerError(w http.ResponseWriter) {
	generateErrorBadge(w, "500 Internal Server Error")
}

// notFound handles HTTP requests for services or methods that don't exist
func notFound(w http.ResponseWriter) {
	generateErrorBadge(w, "404 Not Found")
}

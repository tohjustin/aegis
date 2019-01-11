package main

import (
	"fmt"
	"net/http"

	"github.com/tohjustin/badger/pkg/badge"
)

func badgeServiceWriteResponse(w http.ResponseWriter, response string) {
	// cache response in browser for 1 day (86400), CDN for 1 year (31536000)
	w.Header().Set("Cache-Control", "public, immutable, max-age=86400, s-maxage=31536000")
	w.Header().Set("Content-Type", "image/svg+xml;utf-8")
	w.Write([]byte(response))
}

func badgeServiceErrorHandler(w http.ResponseWriter, r *http.Request) {
	style := r.URL.Query().Get("style")

	createOptions := badge.Options{Style: badge.Style(style)}
	generatedBadge, err := badge.Create("badger", "400 Bad Request", &createOptions)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	badgeServiceWriteResponse(w, generatedBadge)
}

func badgeServiceHandler(w http.ResponseWriter, r *http.Request) {
	color := r.URL.Query().Get("color")
	icon := r.URL.Query().Get("icon")
	status := r.URL.Query().Get("status")
	style := r.URL.Query().Get("style")
	subject := r.URL.Query().Get("subject")

	createOptions := badge.Options{Color: color, Icon: icon, Style: badge.Style(style)}
	generatedBadge, err := badge.Create(subject, status, &createOptions)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	badgeServiceWriteResponse(w, generatedBadge)
}

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
	generatedBadge, err := badge.Create(&badge.Params{
		Style:   badge.Style(r.URL.Query().Get("style")),
		Subject: "badger",
		Status:  "400 Bad Request",
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	badgeServiceWriteResponse(w, generatedBadge)
}

func badgeServiceHandler(w http.ResponseWriter, r *http.Request) {
	generatedBadge, err := badge.Create(&badge.Params{
		Style:   badge.Style(r.URL.Query().Get("style")),
		Subject: r.URL.Query().Get("subject"),
		Status:  r.URL.Query().Get("status"),
		Color:   r.URL.Query().Get("color"),
		Icon:    r.URL.Query().Get("icon"),
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	badgeServiceWriteResponse(w, generatedBadge)
}

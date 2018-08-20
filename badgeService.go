package main

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/tohjustin/badger/pkg/badge"
)

func badgeServiceHandler(w http.ResponseWriter, r *http.Request) {
	routeVariables := mux.Vars(r)
	subject, _ := url.PathUnescape(routeVariables["subject"])
	status, _ := url.PathUnescape(routeVariables["status"])
	color := routeVariables["color"]
	icon := r.URL.Query().Get("icon")
	style := r.URL.Query().Get("style")

	generatedBadge, err := badge.GenerateSVG(style, subject, status, color, icon)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml;utf-8")
	w.Write([]byte(generatedBadge))
}

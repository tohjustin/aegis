package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/hashicorp/golang-lru"
	"github.com/tohjustin/badger/pkg/badge"
)

const badgeServiceCacheSize = 5000
const maxStatusLength = 40
const maxSubjectLength = 40

var (
	badgeServiceCache *lru.Cache
)

func badgeServiceInit() error {
	cache, err := lru.New(badgeServiceCacheSize)
	if err != nil {
		return err
	}

	badgeServiceCache = cache
	return nil
}

func createCharLimitExceededErrorMsg(paramName string, maxLength int) string {
	return fmt.Sprintf("You have exceeded the maximum character limit:\n"+
		" - Received: \"%s\"\n"+
		" - Expected: \"/badge/<SUBJECT>/<STATUS>/<COLOR>\","+
		" where SUBJECT is not more than %d characters long", paramName, maxLength)
}

func badgeServiceHandler(w http.ResponseWriter, r *http.Request) {
	routeVariables := mux.Vars(r)
	subject, _ := url.PathUnescape(routeVariables["subject"])
	status, _ := url.PathUnescape(routeVariables["status"])
	color := routeVariables["color"]
	style := r.URL.Query().Get("style")

	cacheKey := subject + "/" + status + "/" + color + "?style=" + style
	svgBadge, ok := badgeServiceCache.Get(cacheKey)
	if !ok {
		if len(subject) > maxSubjectLength {
			errorMsg := createCharLimitExceededErrorMsg(subject, maxSubjectLength)
			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}

		if len(status) > maxStatusLength {
			errorMsg := createCharLimitExceededErrorMsg(status, maxStatusLength)
			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}

		generatedBadge, err := badge.GenerateSVG(style, subject, status, color)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		svgBadge = generatedBadge
		badgeServiceCache.Add(cacheKey, svgBadge)
	}

	w.Header().Set("Content-Type", "image/svg+xml;utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, svgBadge)
}

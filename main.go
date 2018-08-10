package main

import (
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/hashicorp/golang-lru"
	"github.com/tohjustin/badger/pkg/badge"
	"github.com/urfave/negroni"
)

const maxStatusLength = 40
const maxSubjectLength = 40
const cacheSize = 5000
const defaultPort = "8080"

var (
	badgeServiceCache *lru.Cache
)

func badgeHandler(w http.ResponseWriter, r *http.Request) {
	routeVariables := mux.Vars(r)
	subject, _ := routeVariables["subject"]
	status, _ := routeVariables["status"]
	color := routeVariables["color"]
	style := r.URL.Query().Get("style")

	cacheKey := subject + "/" + status + "/" + color + "?style=" + style
	svgBadge, ok := badgeServiceCache.Get(cacheKey)
	if !ok {
		if len(subject) > maxSubjectLength {
			errorMsg := fmt.Sprintf("Max character length exceeded:\n"+
				" - Received: \"%s\"\n"+
				" - Expected: \"/badge/<SUBJECT>/<STATUS>/<COLOR>\","+
				" where SUBJECT is not more than %d characters long", subject, maxStatusLength)
			http.Error(w, errorMsg, http.StatusBadRequest)
			return
		}

		if len(status) > maxStatusLength {
			errorMsg := fmt.Sprintf("Max character length exceeded:\n"+
				" - Received: \"%s\"\n"+
				" - Expected: \"/badge/<SUBJECT>/<STATUS>/<COLOR>\","+
				" where STATUS is not more than %d characters long", status, maxStatusLength)
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

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	n := negroni.New()

	// Setup recovery middleware
	n.Use(negroni.NewRecovery())

	// Setup logging middleware
	logger := negroni.NewLogger()
	papertrailHost := os.Getenv("PAPERTRAIL_HOST")
	if papertrailHost != "" {
		w, err := syslog.Dial("udp", papertrailHost, 0, "badger-server")
		if err != nil {
			log.Fatal("failed to dial syslog")
		}
		logger.ALogger = log.New(w, "[negroni] ", 0)
	}
	n.Use(logger)

	// Setup LRU cache
	badgeServiceCache, _ = lru.New(cacheSize)

	// Setup routes
	router := mux.NewRouter()
	router.HandleFunc(`/badge/{subject}/{status}/{color}`, badgeHandler).Methods("GET")
	n.UseHandler(router)

	http.ListenAndServe(":"+port, n)
}

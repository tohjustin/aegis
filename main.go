package main

import (
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/tohjustin/badger/pkg/badge"
	"github.com/urfave/negroni"
)

func mapSubexpNames(m, n []string) map[string]string {
	m, n = m[1:], n[1:]
	r := make(map[string]string, len(m))
	for i := range n {
		r[n[i]] = m[i]
	}

	return r
}

func badgeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	badgeParams := vars["badgeParams"]
	badgeStyle := r.URL.Query().Get("style")

	badgeParamsPattern := regexp.MustCompile(`^(?P<subject>.+)-(?P<status>.+)-(?P<color>.+)\.svg$`)
	matched := badgeParamsPattern.FindStringSubmatch(badgeParams)
	if matched == nil {
		errorMsg := fmt.Sprintf("Invalid URL format:\n"+
			" - Received: \"%s\"\n"+
			" - Expected: \"<SUBJECT>-<STATUS>-<COLOR>.svg\"", badgeParams)
		http.Error(w, errorMsg, http.StatusBadRequest)
		return
	}

	result := mapSubexpNames(matched, badgeParamsPattern.SubexpNames())
	svgBadge, err := badge.GenerateSVG(badgeStyle, result["subject"], result["status"], result["color"])
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml;utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, svgBadge)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	n := negroni.New()

	// Setup recovery middleware
	n.Use(negroni.NewRecovery())

	// Setup logging middleware
	logger := negroni.NewLogger()
	if papertrailHost := os.Getenv("PAPERTRAIL_HOST"); papertrailHost != "" {
		w, err := syslog.Dial("udp", papertrailHost, 0, "badger-server")
		if err != nil {
			log.Fatal("failed to dial syslog")
		}
		logger.ALogger = log.New(w, "[negroni] ", 0)
	}
	n.Use(logger)

	// Setup routes
	router := mux.NewRouter()
	router.HandleFunc(`/{badgeParams}`, badgeHandler).Methods("GET")
	n.UseHandler(router)

	http.ListenAndServe(":"+port, n)
}

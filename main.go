package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/mux"
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
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request")
		return
	}

	// TODO: Validate <color>
	result := mapSubexpNames(matched, badgeParamsPattern.SubexpNames())
	svgBadge, err := generateBadge(badgeStyle, result["subject"], result["status"], result["color"])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Internal Server Error")
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

	r := mux.NewRouter()
	r.HandleFunc(`/badge/{badgeParams}`, badgeHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+port, r))
}

package main

import (
	"fmt"
	"log"
	"log/syslog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

const defaultPort = "8080"

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

	// initialize handlers
	badgeServiceInit()

	// setup router
	mux := mux.NewRouter()
	mux.UseEncodedPath()
	mux.HandleFunc(`/badge/{subject}/{status}`, badgeServiceHandler).Methods("GET")
	mux.HandleFunc(`/badge/{subject}/{status}/{color}`, badgeServiceHandler).Methods("GET")

	n.UseHandler(mux)

	http.ListenAndServe(":"+port, n)
}

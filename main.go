package main

import (
	"log"
	"log/syslog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

const defaultPort = "8080"

func main() {
	logEndpoint := os.Getenv("PAPERTRAIL_HOST")
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// initialize handlers
	badgeServiceInit()

	// setup router
	mux := mux.NewRouter()
	mux.UseEncodedPath()
	mux.HandleFunc(`/badge/{subject}/{status}`, badgeServiceHandler).Methods("GET")
	mux.HandleFunc(`/badge/{subject}/{status}/{color}`, badgeServiceHandler).Methods("GET")

	// setup middlewares
	n := negroni.New()
	logger := negroni.NewLogger()
	if logEndpoint != "" {
		w, err := syslog.Dial("udp", logEndpoint, 0, "badger-server")
		if err != nil {
			log.Fatal("Failed to dial syslog at: ", logEndpoint)
			return
		}

		logger.ALogger = log.New(w, "[negroni] ", 0)
	}
	n.Use(logger)
	n.Use(newRecoveryMiddleware())
	n.UseHandler(mux)

	http.ListenAndServe(":"+port, n)
}

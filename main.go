package main

import (
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

	// setup router
	mux := mux.NewRouter()
	mux.UseEncodedPath()
	mux.HandleFunc(`/badge/{subject}/{status}`, badgeServiceHandler).Methods("GET")
	mux.HandleFunc(`/badge/{subject}/{status}/{color}`, badgeServiceHandler).Methods("GET")

	// setup middlewares
	n := negroni.New()
	n.Use(newLoggerMiddleware(logEndpoint))
	n.Use(newRecoveryMiddleware())
	n.UseHandler(mux)

	http.ListenAndServe(":"+port, n)
}

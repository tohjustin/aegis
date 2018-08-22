package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

const defaultPort = "8080"

func newRouter() http.Handler {
	mux := mux.NewRouter()
	mux.UseEncodedPath()
	mux.HandleFunc(`/badge/{subject}/{status}/{color}`, badgeServiceHandler).Methods("GET")
	mux.HandleFunc(`/badge/{subject}/{status}`, badgeServiceHandler).Methods("GET")
	mux.PathPrefix("/").HandlerFunc(badgeServiceErrorHandler).Methods("GET")

	return mux
}

func main() {
	logEndpoint := os.Getenv("PAPERTRAIL_HOST")
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	n := negroni.New()
	n.Use(newLoggerMiddleware(logEndpoint))
	n.Use(newRecoveryMiddleware())
	n.UseHandler(newRouter())

	http.ListenAndServe(":"+port, n)
}

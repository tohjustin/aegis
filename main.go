package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

const defaultPort = "8080"

func newRouter() http.Handler {
	githubService := NewGithubService()

	mux := mux.NewRouter()
	mux.UseEncodedPath()
	mux.HandleFunc(`/static`, badgeServiceHandler).Methods("GET")
	mux.HandleFunc(`/github/{owner}/{repo}/{requestType}`, githubService.Handler).Methods("GET")
	mux.PathPrefix("/").HandlerFunc(badgeServiceErrorHandler).Methods("GET")

	return mux
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	n := negroni.New()
	n.Use(newLoggerMiddleware())
	n.Use(newRecoveryMiddleware())
	n.UseHandler(newRouter())

	srv := http.Server{
		Addr:         ":" + port,
		Handler:      n,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Printf("HTTP service listening on port %s...", port)

	// gracefully shutdowns server
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(nil); err != nil {
			log.Printf("HTTP service Shutdown: %v", err)
		}

		log.Printf("HTTP service shutdown successfully...")
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("HTTP service ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

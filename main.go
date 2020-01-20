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

// BadgeService represents a badge service
type BadgeService interface {
	http.Handler
}

// GitRepositoryService represents a badge service for git respository providers
type GitRepositoryService interface {
	BadgeService
	getForkCount(owner string, repo string) (int, error)
	getIssueCount(owner string, repo string, issueState string) (int, error)
	getPullRequestCount(owner string, repo string, pullRequestState string) (int, error)
	getStarCount(owner string, repo string) (int, error)
}

func newRouter() http.Handler {
	staticService := newStaticServiceHandler()
	bitbucketService := newBitbucketServiceHandler()
	githubService := newGithubServiceHandler()
	gitlabService := newGitlabServiceHandler()

	mux := mux.NewRouter()
	mux.UseEncodedPath()
	mux.Handle(`/static`, staticService).Methods("GET")
	mux.Handle(`/bitbucket/{owner}/{repo}/{requestType}`, bitbucketService).Methods("GET")
	mux.Handle(`/github/{owner}/{repo}/{requestType}`, githubService).Methods("GET")
	mux.Handle(`/gitlab/{owner}/{repo}/{requestType}`, gitlabService).Methods("GET")

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

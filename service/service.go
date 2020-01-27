package service

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

const (
	// DefaultPort default port server listens on
	DefaultPort = 8080
)

// BadgeService represents a badge service
type BadgeService interface {
	http.Handler
}

// GitProviderService represents a badge service for git providers
type GitProviderService interface {
	BadgeService
	getForkCount(owner string, repo string) (int, error)
	getIssueCount(owner string, repo string, issueState string) (int, error)
	getPullRequestCount(owner string, repo string, pullRequestState string) (int, error)
	getStarCount(owner string, repo string) (int, error)
}

// Info contains build information about the application
type Info struct {
	ExecutableName string
	LongName       string
	Version        string
	GitHash        string
}

// Config represents a server's configuration
type Config struct {
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Application represents a badge generation application
type Application struct {
	info   Info
	config *Config

	staticService    *BadgeService
	bitbucketService *GitProviderService
	githubService    *GitProviderService
	gitlabService    *GitProviderService
}

// handler setup routes & returns a HTTP handler for the application server
func (app *Application) handler() http.Handler {
	mux := mux.NewRouter()

	mux.UseEncodedPath()
	mux.Handle(`/static`, *app.staticService).Methods("GET")
	mux.Handle(`/bitbucket/{method}/{owner}/{repo}`, *app.bitbucketService).Methods("GET")
	mux.Handle(`/github/{method}/{owner}/{repo}`, *app.githubService).Methods("GET")
	mux.Handle(`/gitlab/{method}/{owner}/{repo}`, *app.gitlabService).Methods("GET")

	return mux
}

// Start starts the application
func (app *Application) Start() {
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.Port),
		ReadTimeout:  app.config.ReadTimeout,
		WriteTimeout: app.config.WriteTimeout,
		Handler:      app.handler(),
	}

	// gracefully shutdowns server
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := httpServer.Shutdown(nil); err != nil {
			log.Printf("HTTP service Shutdown: %v", err)
		}

		log.Printf("HTTP service shutdown successfully...")
		close(idleConnsClosed)
	}()

	log.Printf("HTTP service listening on port %d...", app.config.Port)
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("HTTP service ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

// newConfig returns a new set of server configuration
func newConfig(port int) *Config {
	return &Config{
		Port:         port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

// New creates and returns a new instance of Application.
func New(appInfo Info) (*Application, error) {
	githubAccessToken := os.Getenv("GITHUB_ACCESS_TOKEN")

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = DefaultPort
	}
	config := newConfig(port)

	staticService := NewStaticService()
	bitbucketService := NewBitbucketService()
	githubService, err := NewGithubService(githubAccessToken)
	if err != nil {
		return nil, err
	}
	gitlabService := NewGitlabService()

	app := &Application{
		info:             appInfo,
		config:           config,
		staticService:    &staticService,
		bitbucketService: &bitbucketService,
		githubService:    &githubService,
		gitlabService:    &gitlabService,
	}

	return app, nil
}

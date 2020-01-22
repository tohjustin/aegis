package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

const (
	// DefaultPort default port server listens on
	DefaultPort = "8080"
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

// Config represents a server's configuration
type Config struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Server represents a server instance
type Server struct {
	config *Config

	staticService    *BadgeService
	bitbucketService *GitProviderService
	githubService    *GitProviderService
	gitlabService    *GitProviderService
}

// Handler returns a handler for the server
func (s *Server) Handler() http.Handler {
	mux := mux.NewRouter()

	mux.UseEncodedPath()
	mux.Handle(`/static`, *s.staticService).Methods("GET")
	mux.Handle(`/bitbucket/{method}/{owner}/{repo}`, *s.bitbucketService).Methods("GET")
	mux.Handle(`/github/{method}/{owner}/{repo}`, *s.githubService).Methods("GET")
	mux.Handle(`/gitlab/{method}/{owner}/{repo}`, *s.gitlabService).Methods("GET")

	return mux
}

// Run starts the server
func (s *Server) Run() {
	httpServer := &http.Server{
		Addr:         ":" + s.config.Port,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		Handler:      s.Handler(),
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

	log.Printf("HTTP service listening on port %s...", s.config.Port)
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("HTTP service ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

// NewConfig returns a new set of server configuration
func NewConfig(port string) *Config {
	if port == "" {
		port = DefaultPort
	}

	return &Config{
		Port:         port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

// NewServer returns a new server instance
func NewServer(
	config *Config,
	staticService *BadgeService,
	bitbucketService *GitProviderService,
	githubService *GitProviderService,
	gitlabService *GitProviderService,
) *Server {
	return &Server{
		config:           config,
		staticService:    staticService,
		bitbucketService: bitbucketService,
		githubService:    githubService,
		gitlabService:    gitlabService,
	}
}

func main() {
	port := os.Getenv("PORT")
	githubAccessToken := os.Getenv("GITHUB_ACCESS_TOKEN")

	config := NewConfig(port)
	staticService := NewStaticService()
	bitbucketService := NewBitbucketService()
	githubService, err := NewGithubService(githubAccessToken)
	if err != nil {
		panic(err)
	}
	gitlabService := NewGitlabService()

	server := NewServer(
		config,
		&staticService,
		&bitbucketService,
		&githubService,
		&gitlabService,
	)

	server.Run()
}

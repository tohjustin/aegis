package service

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/tohjustin/badger/service/config"
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
	ShortName      string
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
	info    Info
	rootCmd *cobra.Command

	staticService    *BadgeService
	bitbucketService *GitProviderService
	githubService    *GitProviderService
	gitlabService    *GitProviderService
}

func (app *Application) execute() {
	// Setup dependencies
	staticService := NewStaticService()
	bitbucketService := NewBitbucketService()
	githubService, err := NewGithubService(os.Getenv("GITHUB_ACCESS_TOKEN"))
	if err != nil {
		log.Fatalf("Unable to setup GitHub service: %v", err)
	}
	gitlabService := NewGitlabService()
	app.staticService = &staticService
	app.bitbucketService = &bitbucketService
	app.githubService = &githubService
	app.gitlabService = &gitlabService

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port()),
		ReadTimeout:  config.ReadTimeout(),
		WriteTimeout: config.WriteTimeout(),
		Handler:      app.handler(),
	}

	// gracefully shutdowns server
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := httpServer.Shutdown(nil); err != nil {
			log.Printf("Server Shutdown: %v\n", err)
		}

		log.Printf("Server shutdown successfully...\n")
		close(idleConnsClosed)
	}()

	// Start HTTP server
	log.Printf("Server listening on port %d...\n", config.Port())
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("Server ListenAndServe: %v\n", err)
	}

	<-idleConnsClosed
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
func (app *Application) Start() error {
	return app.rootCmd.Execute()
}

// New creates and returns a new instance of Application.
func New(appInfo Info) (*Application, error) {
	app := &Application{
		info: appInfo,
	}

	// Setup commands
	rootCmd := &cobra.Command{
		Use:   appInfo.ExecutableName,
		Short: appInfo.ShortName,
		Long:  appInfo.LongName,
		Run: func(cmd *cobra.Command, args []string) {
			app.execute()
		},
	}
	versionCmd := &cobra.Command{
		Use:  "version",
		Long: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Badger v%s (%s)\n", appInfo.Version, appInfo.GitHash)
		},
	}
	rootCmd.AddCommand(versionCmd)

	// Setup Flags
	flagSet := new(flag.FlagSet)
	config.Flags(flagSet)
	rootCmd.Flags().AddGoFlagSet(flagSet)

	app.rootCmd = rootCmd

	return app, nil
}

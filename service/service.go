package service

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/tohjustin/aegis/service/config"
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

// Application represents a badge generation application
type Application struct {
	info    Info
	config  *config.Config
	logger  *zap.Logger
	rootCmd *cobra.Command

	staticService    *BadgeService
	bitbucketService *GitProviderService
	githubService    *GitProviderService
	gitlabService    *GitProviderService
}

func (app *Application) init() {
	logger, err := newLogger()
	if err != nil {
		log.Fatalf("Failed to get logger: %v", err)
	}
	app.logger = logger

	config, err := config.New()
	if err != nil {
		log.Fatalf("Failed to get config: %v", err)
	}
	app.config = config
}

func (app *Application) execute() {
	app.logger.Info("Starting "+app.info.LongName+"...",
		zap.String("Version", app.info.Version),
		zap.String("GitHash", app.info.GitHash),
		zap.Int("NumCPU", runtime.NumCPU()))

	// Setup dependencies
	app.logger.Info("Initializing services...")
	staticService, err := NewStaticService(app.config, app.logger)
	if err != nil {
		log.Fatalf("Failed to get static service: %v", err)
	}
	bitbucketService, err := NewBitbucketService(app.config, app.logger)
	if err != nil {
		log.Fatalf("Failed to get Bitbucket service: %v", err)
	}
	githubService, err := NewGithubService(app.config, app.logger)
	if err != nil {
		log.Fatalf("Failed to get GitHub service: %v", err)
	}
	gitlabService, err := NewGitlabService(app.config, app.logger)
	if err != nil {
		log.Fatalf("Failed to get GitLab service: %v", err)
	}
	app.staticService = &staticService
	app.bitbucketService = &bitbucketService
	app.githubService = &githubService
	app.gitlabService = &gitlabService

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
		s := <-sigint
		app.logger.Info("Received signal from OS", zap.String("signal", s.String()))

		app.logger.Info("Starting shutdown...")
		if err := httpServer.Shutdown(context.Background()); err != nil {
			app.logger.Error("Encountered error during shutdown", zap.Error(err))
		}

		app.logger.Info("Shutdown complete.")
		close(idleConnsClosed)
	}()

	// Start HTTP server
	app.logger.Info("HTTP server listening...", zap.Uint("Port", app.config.Port))
	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		app.logger.Error("HTTP server encountered an error", zap.Error(err))
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

	if url := app.config.RootRedirectURL; url != "" {
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, url, http.StatusFound)
		}).Methods("GET")
	}
	// return service-not-found badge for all unmatched routes
	mux.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := serviceNotFound(w, app.config)
		if err != nil {
			app.logger.Error("Failed to create error badge",
				zap.String("url", r.URL.RequestURI()),
				zap.Error(err))
		}
	}).Methods("GET")

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
		Use:  appInfo.ExecutableName,
		Long: appInfo.LongName,
		Run: func(cmd *cobra.Command, args []string) {
			app.init()
			app.execute()
		},
	}
	versionCmd := &cobra.Command{
		Use:  "version",
		Long: "Print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s v%s (%s)\n", appInfo.ShortName, appInfo.Version, appInfo.GitHash)
		},
	}
	rootCmd.AddCommand(versionCmd)

	// Setup Flags
	flagSet := new(flag.FlagSet)
	addFlagsFns := []func(*flag.FlagSet){
		config.Flags,
		loggerFlags,
	}
	for _, addFlags := range addFlagsFns {
		addFlags(flagSet)
	}
	rootCmd.Flags().AddGoFlagSet(flagSet)

	app.rootCmd = rootCmd

	return app, nil
}

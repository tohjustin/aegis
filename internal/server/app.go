package server

import (
	"log/slog"
	"net/http"

	"github.com/tohjustin/aegis/internal/config"
	"github.com/tohjustin/aegis/internal/gitprovider"
)

type app struct {
	cfg       *config.Config
	logger    *slog.Logger
	mux       *http.ServeMux
	github    *gitprovider.GitHub
	gitlab    *gitprovider.GitLab
	bitbucket *gitprovider.Bitbucket
}

func newAppHandler(cfg *config.Config, logger *slog.Logger, github *gitprovider.GitHub, gitlab *gitprovider.GitLab, bitbucket *gitprovider.Bitbucket) http.Handler {
	a := &app{
		cfg:       cfg,
		logger:    logger,
		mux:       http.NewServeMux(),
		github:    github,
		gitlab:    gitlab,
		bitbucket: bitbucket,
	}

	a.mux.HandleFunc("GET /static", a.handleStatic)
	a.mux.HandleFunc("GET /github/{method}/{owner}/{repo}", a.handleGit("github", githubMethods(github)))
	a.mux.HandleFunc("GET /gitlab/{method}/{owner}/{repo}", a.handleGit("gitlab", gitlabMethods(gitlab)))
	a.mux.HandleFunc("GET /bitbucket/{method}/{owner}/{repo}", a.handleGit("bitbucket", bitbucketMethods(bitbucket)))
	if cfg.RootRedirectURL != "" {
		a.mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, cfg.RootRedirectURL, http.StatusFound)
		})
	}
	return a
}

func (a *app) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, pattern := a.mux.Handler(r)
	if pattern == "" {
		if r.Method == http.MethodGet {
			if err := a.writeErrorBadge(w, "service not found"); err != nil {
				a.logger.LogAttrs(r.Context(), slog.LevelError, "failed to create error badge",
					slog.String("url", r.URL.RequestURI()),
					slog.Any("err", err),
				)
			}
			return
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	a.mux.ServeHTTP(w, r)
}

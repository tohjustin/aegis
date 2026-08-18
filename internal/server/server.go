// Package server wires application and admin HTTP listeners.
package server

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"sync/atomic"

	"github.com/tohjustin/aegis/internal/config"
	"github.com/tohjustin/aegis/internal/gitprovider"
	"github.com/tohjustin/aegis/internal/observability"
)

const maxBodyBytes = 1 << 20

// Server owns the app and admin HTTP servers plus the ready flag.
type Server struct {
	app    *http.Server
	admin  *http.Server
	logger *slog.Logger
	ready  *atomic.Bool
}

// New constructs app and admin servers without binding listeners.
func New(cfg *config.Config, logger *slog.Logger, metrics *observability.Metrics) (*Server, error) {
	github, err := gitprovider.NewGitHub(cfg.GithubAccessToken)
	if err != nil {
		return nil, err
	}

	ready := &atomic.Bool{}
	appHandler := newAppHandler(cfg, logger, github, gitprovider.NewGitLab(nil), gitprovider.NewBitbucket(nil))
	appHandler = maxBytes(maxBodyBytes)(appHandler)
	appHandler = observability.Recover(logger)(appHandler)
	appHandler = observability.Instrument(logger, metrics)(appHandler)

	adminHandler := observability.Recover(logger)(newAdminHandler(metrics, ready))

	errorLog := slog.NewLogLogger(logger.Handler(), slog.LevelError)
	return &Server{
		app: &http.Server{
			ReadTimeout:       cfg.ReadTimeout,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       cfg.IdleTimeout,
			ErrorLog:          errorLog,
			Handler:           appHandler,
		},
		admin: &http.Server{
			ReadTimeout:       cfg.ReadTimeout,
			ReadHeaderTimeout: cfg.ReadHeaderTimeout,
			WriteTimeout:      cfg.WriteTimeout,
			IdleTimeout:       cfg.IdleTimeout,
			ErrorLog:          errorLog,
			Handler:           adminHandler,
		},
		logger: logger,
		ready:  ready,
	}, nil
}

// ServeApp serves application traffic on ln.
func (s *Server) ServeApp(ln net.Listener) error {
	return serve(s.app, ln)
}

// ServeAdmin serves admin traffic on ln.
func (s *Server) ServeAdmin(ln net.Listener) error {
	return serve(s.admin, ln)
}

// SetReady toggles the readiness flag used by /readyz.
func (s *Server) SetReady(v bool) {
	s.ready.Store(v)
}

// Shutdown stops the app server first, then the admin server.
func (s *Server) Shutdown(ctx context.Context) error {
	s.SetReady(false)
	return errors.Join(s.app.Shutdown(ctx), s.admin.Shutdown(ctx))
}

func serve(srv *http.Server, ln net.Listener) error {
	err := srv.Serve(ln)
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func maxBytes(n int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, n)
			next.ServeHTTP(w, r)
		})
	}
}

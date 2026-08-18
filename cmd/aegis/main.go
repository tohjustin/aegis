package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/tohjustin/aegis/internal/config"
	"github.com/tohjustin/aegis/internal/observability"
	"github.com/tohjustin/aegis/internal/server"
	"github.com/tohjustin/aegis/internal/version"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	if wantsVersion(os.Args[1:]) {
		fmt.Println(version.String())
		return nil
	}

	cfg, err := config.Load(os.Args[1:], os.Stderr)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return nil
		}
		return err
	}

	logger, err := observability.NewLogger(cfg.LogLevel, cfg.LogFormat, os.Stdout)
	if err != nil {
		return err
	}

	info := version.Current()
	metrics, err := observability.NewMetrics(cfg.HistogramBuckets, info)
	if err != nil {
		return err
	}

	srv, err := server.New(cfg, logger, metrics)
	if err != nil {
		return err
	}

	appLn, err := net.Listen("tcp", cfg.AppAddr)
	if err != nil {
		return fmt.Errorf("listen app: %w", err)
	}
	adminLn, err := net.Listen("tcp", cfg.AdminAddr)
	if err != nil {
		_ = appLn.Close()
		return fmt.Errorf("listen admin: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	logger.LogAttrs(ctx, slog.LevelInfo, "starting",
		slog.String("version", info.Version),
		slog.String("git_commit", info.GitCommit),
		slog.String("git_tag", info.GitTag),
		slog.String("platform", info.Platform),
		slog.String("app_addr", appLn.Addr().String()),
		slog.String("admin_addr", adminLn.Addr().String()),
	)

	errc := make(chan error, 2)
	go func() { errc <- srv.ServeApp(appLn) }()
	go func() { errc <- srv.ServeAdmin(adminLn) }()
	srv.SetReady(true)

	var runErr error
	select {
	case <-ctx.Done():
		logger.LogAttrs(context.Background(), slog.LevelInfo, "shutdown signal received")
	case runErr = <-errc:
		if runErr != nil {
			logger.LogAttrs(context.Background(), slog.LevelError, "server exited", slog.Any("err", runErr))
		}
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	return errors.Join(runErr, srv.Shutdown(shutdownCtx))
}

func wantsVersion(args []string) bool {
	for _, arg := range args {
		if arg == "--" {
			return false
		}
		if arg == "-v" || arg == "--version" {
			return true
		}
	}
	return false
}

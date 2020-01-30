package config

import (
	"flag"
	"fmt"
	"os"
	"time"
)

const (
	portCfg                       = "port"
	readTimeoutCfg                = "read-timeout"
	writeTimeoutCfg               = "write-timeout"
	excludeCacheControlHeadersCfg = "exclude-cache-control-headers"
	githubAccessTokenCfg          = "github-access-token"
)

var (
	port                       *uint
	readTimeout                *uint
	writeTimeout               *uint
	excludeCacheControlHeaders *bool
	githubAccessToken          *string
)

// Config contains all application configuration
type Config struct {
	Port                       uint
	ReadTimeout                time.Duration
	WriteTimeout               time.Duration
	ExcludeCacheControlHeaders bool
	GithubAccessToken          string
}

// Flags adds flags related to the application to the given flagset.
func Flags(flags *flag.FlagSet) {
	// server configs
	port = flags.Uint(portCfg, 8080, "Port exposing badget service.")
	readTimeout = flags.Uint(readTimeoutCfg, 2000, "Maximum duration in milliseconds for reading the entire request, including the body.")
	writeTimeout = flags.Uint(writeTimeoutCfg, 2000, "Maximum duration in milliseconds before timing out writes of the response.")
	excludeCacheControlHeaders = flags.Bool(excludeCacheControlHeadersCfg, false, "Flag to exclude HTTP Cache-Control headers from responses.")

	// service configs
	githubAccessToken = flags.String(githubAccessTokenCfg, os.Getenv("GITHUB_ACCESS_TOKEN"), "GitHub Access Token for GitHub badge service.")
}

// New returns an instance of all application configuration
func New() (*Config, error) {
	if port == nil || readTimeout == nil || writeTimeout == nil ||
		excludeCacheControlHeaders == nil || githubAccessToken == nil {
		return nil, fmt.Errorf("configuration flags are not set")
	}

	return &Config{
		Port:                       *port,
		ReadTimeout:                time.Duration(*readTimeout) * time.Millisecond,
		WriteTimeout:               time.Duration(*writeTimeout) * time.Millisecond,
		ExcludeCacheControlHeaders: *excludeCacheControlHeaders,
		GithubAccessToken:          *githubAccessToken,
	}, nil
}

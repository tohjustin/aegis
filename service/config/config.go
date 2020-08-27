package config

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"time"
)

const (
	portCfg                       = "port"
	readTimeoutCfg                = "read-timeout"
	writeTimeoutCfg               = "write-timeout"
	excludeCacheControlHeadersCfg = "exclude-cache-control-headers"
	rootRedirectURLCfg            = "root-redirect-url"
	githubAccessTokenCfg          = "github-access-token"
)

var (
	port                       *uint
	readTimeout                *uint
	writeTimeout               *uint
	excludeCacheControlHeaders *bool
	rootRedirectURL            *string
	githubAccessToken          *string
)

// Config contains all application configuration
type Config struct {
	Port                       uint
	ReadTimeout                time.Duration
	WriteTimeout               time.Duration
	ExcludeCacheControlHeaders bool
	RootRedirectURL            string
	GithubAccessToken          string
}

// Flags adds flags related to the application to the given flagset.
func Flags(flags *flag.FlagSet) {
	// server configs
	port = flags.Uint(portCfg, 8080, "Port exposing badge service.")
	readTimeout = flags.Uint(readTimeoutCfg, 2000, "Maximum duration in milliseconds for reading the entire request, including the body.")
	writeTimeout = flags.Uint(writeTimeoutCfg, 2000, "Maximum duration in milliseconds before timing out writes of the response.")
	excludeCacheControlHeaders = flags.Bool(excludeCacheControlHeadersCfg, false, "Flag to exclude HTTP Cache-Control headers from responses.")
	rootRedirectURL = flags.String(rootRedirectURLCfg, os.Getenv("ROOT_REDIRECT_URL"), "URL to redirect for all root path requests.")

	// service configs
	githubAccessToken = flags.String(githubAccessTokenCfg, os.Getenv("GITHUB_ACCESS_TOKEN"), "GitHub Access Token for GitHub badge service.")
}

// New returns an instance of all application configuration
func New() (*Config, error) {
	if port == nil || readTimeout == nil || writeTimeout == nil ||
		excludeCacheControlHeaders == nil || githubAccessToken == nil {
		return nil, fmt.Errorf("configuration flags are not set")
	}

	if *rootRedirectURL != "" {
		if _, err := url.ParseRequestURI(*rootRedirectURL); err != nil {
			return nil, fmt.Errorf("Config.RootRedirectURL URL is invalid: %s", *rootRedirectURL)
		}
	}

	return &Config{
		Port:                       *port,
		ReadTimeout:                time.Duration(*readTimeout) * time.Millisecond,
		WriteTimeout:               time.Duration(*writeTimeout) * time.Millisecond,
		ExcludeCacheControlHeaders: *excludeCacheControlHeaders,
		RootRedirectURL:            *rootRedirectURL,
		GithubAccessToken:          *githubAccessToken,
	}, nil
}

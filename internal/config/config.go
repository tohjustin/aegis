// Package config loads process configuration from defaults, environment, then flags.
package config

import (
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultAppAddr           = ":8080"
	defaultAdminAddr         = ":9090"
	defaultLogLevel          = "info"
	defaultLogFormat         = "json"
	defaultReadTimeout       = 15 * time.Second
	defaultReadHeaderTimeout = 5 * time.Second
	defaultWriteTimeout      = 15 * time.Second
	defaultIdleTimeout       = 60 * time.Second
	defaultShutdownTimeout   = 15 * time.Second
	defaultHistogramBuckets  = "0.005,0.01,0.025,0.05,0.1,0.25,0.5,1,2.5,5,10"
)

// Config is the fully validated process configuration.
type Config struct {
	AppAddr                    string
	AdminAddr                  string
	LogLevel                   string
	LogFormat                  string
	ReadTimeout                time.Duration
	ReadHeaderTimeout          time.Duration
	WriteTimeout               time.Duration
	IdleTimeout                time.Duration
	ShutdownTimeout            time.Duration
	HistogramBuckets           []float64
	ExcludeCacheControlHeaders bool
	RootRedirectURL            string
	GithubAccessToken          string
}

// Load reads defaults, then environment variables, then CLI flags (flags win).
// Errors are intended for stderr; logging is not used.
func Load(args []string, errOut io.Writer) (*Config, error) {
	fs := flag.NewFlagSet("aegis", flag.ContinueOnError)
	fs.SetOutput(errOut)

	appAddrDefault, err := envString("APP_ADDR", defaultAppAddr)
	if err != nil {
		return nil, err
	}
	adminAddrDefault, err := envString("ADMIN_ADDR", defaultAdminAddr)
	if err != nil {
		return nil, err
	}
	logLevelDefault, err := envString("LOG_LEVEL", defaultLogLevel)
	if err != nil {
		return nil, err
	}
	logFormatDefault, err := envString("LOG_FORMAT", defaultLogFormat)
	if err != nil {
		return nil, err
	}
	readTimeoutDefault, err := envDuration("READ_TIMEOUT", defaultReadTimeout)
	if err != nil {
		return nil, err
	}
	readHeaderTimeoutDefault, err := envDuration("READ_HEADER_TIMEOUT", defaultReadHeaderTimeout)
	if err != nil {
		return nil, err
	}
	writeTimeoutDefault, err := envDuration("WRITE_TIMEOUT", defaultWriteTimeout)
	if err != nil {
		return nil, err
	}
	idleTimeoutDefault, err := envDuration("IDLE_TIMEOUT", defaultIdleTimeout)
	if err != nil {
		return nil, err
	}
	shutdownTimeoutDefault, err := envDuration("SHUTDOWN_TIMEOUT", defaultShutdownTimeout)
	if err != nil {
		return nil, err
	}
	bucketsDefault, err := envString("HISTOGRAM_BUCKETS", defaultHistogramBuckets)
	if err != nil {
		return nil, err
	}
	excludeDefault, err := envBool("EXCLUDE_CACHE_CONTROL_HEADERS", false)
	if err != nil {
		return nil, err
	}
	rootRedirectDefault, err := envString("ROOT_REDIRECT_URL", "")
	if err != nil {
		return nil, err
	}
	githubTokenDefault, err := envString("GITHUB_ACCESS_TOKEN", "")
	if err != nil {
		return nil, err
	}

	appAddr := fs.String("app-addr", appAddrDefault, "Address for application traffic.")
	adminAddr := fs.String("admin-addr", adminAddrDefault, "Address for admin traffic (probes, metrics, version).")
	logLevel := fs.String("log-level", logLevelDefault, "Log level (debug, info, warn, error).")
	logFormat := fs.String("log-format", logFormatDefault, "Log format (json, text).")
	readTimeout := fs.Duration("read-timeout", readTimeoutDefault, "Maximum duration for reading the entire request.")
	readHeaderTimeout := fs.Duration("read-header-timeout", readHeaderTimeoutDefault, "Maximum duration for reading request headers.")
	writeTimeout := fs.Duration("write-timeout", writeTimeoutDefault, "Maximum duration before timing out writes of the response.")
	idleTimeout := fs.Duration("idle-timeout", idleTimeoutDefault, "Maximum idle time for keep-alive connections.")
	shutdownTimeout := fs.Duration("shutdown-timeout", shutdownTimeoutDefault, "Graceful shutdown timeout.")
	histogramBuckets := fs.String("histogram-buckets", bucketsDefault, "Comma-separated Prometheus histogram buckets in seconds.")
	excludeCacheControlHeaders := fs.Bool("exclude-cache-control-headers", excludeDefault, "Exclude Cache-Control headers from badge responses.")
	rootRedirectURL := fs.String("root-redirect-url", rootRedirectDefault, "URL to redirect for GET / requests.")
	githubAccessToken := fs.String("github-access-token", githubTokenDefault, "GitHub access token for the GitHub badge service.")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if fs.NArg() > 0 {
		return nil, fmt.Errorf("unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}

	buckets, err := parseBuckets(*histogramBuckets)
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		AppAddr:                    *appAddr,
		AdminAddr:                  *adminAddr,
		LogLevel:                   strings.ToLower(*logLevel),
		LogFormat:                  strings.ToLower(*logFormat),
		ReadTimeout:                *readTimeout,
		ReadHeaderTimeout:          *readHeaderTimeout,
		WriteTimeout:               *writeTimeout,
		IdleTimeout:                *idleTimeout,
		ShutdownTimeout:            *shutdownTimeout,
		HistogramBuckets:           buckets,
		ExcludeCacheControlHeaders: *excludeCacheControlHeaders,
		RootRedirectURL:            *rootRedirectURL,
		GithubAccessToken:          *githubAccessToken,
	}
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	if strings.TrimSpace(c.AppAddr) == "" {
		return fmt.Errorf("app-addr must be non-empty")
	}
	if strings.TrimSpace(c.AdminAddr) == "" {
		return fmt.Errorf("admin-addr must be non-empty")
	}
	if c.AppAddr == c.AdminAddr {
		return fmt.Errorf("app-addr and admin-addr must differ")
	}
	switch strings.ToLower(c.LogLevel) {
	case "debug", "info", "warn", "error":
	default:
		return fmt.Errorf("log-level must be one of debug, info, warn, error")
	}
	switch strings.ToLower(c.LogFormat) {
	case "json", "text":
	default:
		return fmt.Errorf("log-format must be one of json, text")
	}
	for _, d := range []struct {
		name string
		val  time.Duration
	}{
		{"read-timeout", c.ReadTimeout},
		{"read-header-timeout", c.ReadHeaderTimeout},
		{"write-timeout", c.WriteTimeout},
		{"idle-timeout", c.IdleTimeout},
		{"shutdown-timeout", c.ShutdownTimeout},
	} {
		if d.val <= 0 {
			return fmt.Errorf("%s must be positive", d.name)
		}
	}
	if err := validateBuckets(c.HistogramBuckets); err != nil {
		return err
	}
	if c.RootRedirectURL != "" {
		if _, err := url.ParseRequestURI(c.RootRedirectURL); err != nil {
			return fmt.Errorf("root-redirect-url is invalid: %s", c.RootRedirectURL)
		}
	}
	if strings.TrimSpace(c.GithubAccessToken) == "" {
		return fmt.Errorf("github-access-token is required")
	}
	return nil
}

func envString(key, fallback string) (string, error) {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback, nil
	}
	return v, nil
}

func envDuration(key string, fallback time.Duration) (time.Duration, error) {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %w", key, err)
	}
	return d, nil
}

func envBool(key string, fallback bool) (bool, error) {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback, nil
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return false, fmt.Errorf("invalid %s: %w", key, err)
	}
	return b, nil
}

func parseBuckets(s string) ([]float64, error) {
	parts := strings.Split(s, ",")
	out := make([]float64, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		f, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid histogram bucket %q: %w", part, err)
		}
		out = append(out, f)
	}
	return out, nil
}

func validateBuckets(buckets []float64) error {
	if len(buckets) == 0 {
		return fmt.Errorf("histogram buckets must be non-empty")
	}
	for i, v := range buckets {
		if v <= 0 {
			return fmt.Errorf("histogram bucket %d must be positive", i)
		}
		if i > 0 && v <= buckets[i-1] {
			return fmt.Errorf("histogram buckets must be strictly ascending")
		}
	}
	return nil
}

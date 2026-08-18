package config

import (
	"io"
	"testing"
	"time"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("GITHUB_ACCESS_TOKEN", "token")
	for _, key := range []string{
		"APP_ADDR", "ADMIN_ADDR", "LOG_LEVEL", "LOG_FORMAT",
		"READ_TIMEOUT", "READ_HEADER_TIMEOUT", "WRITE_TIMEOUT",
		"IDLE_TIMEOUT", "SHUTDOWN_TIMEOUT", "HISTOGRAM_BUCKETS",
		"EXCLUDE_CACHE_CONTROL_HEADERS", "ROOT_REDIRECT_URL",
	} {
		t.Setenv(key, "")
	}

	cfg, err := Load(nil, io.Discard)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.AppAddr != defaultAppAddr {
		t.Errorf("AppAddr = %q, want %q", cfg.AppAddr, defaultAppAddr)
	}
	if cfg.AdminAddr != defaultAdminAddr {
		t.Errorf("AdminAddr = %q, want %q", cfg.AdminAddr, defaultAdminAddr)
	}
	if cfg.LogLevel != defaultLogLevel {
		t.Errorf("LogLevel = %q, want %q", cfg.LogLevel, defaultLogLevel)
	}
	if cfg.LogFormat != defaultLogFormat {
		t.Errorf("LogFormat = %q, want %q", cfg.LogFormat, defaultLogFormat)
	}
	if cfg.ReadTimeout != defaultReadTimeout {
		t.Errorf("ReadTimeout = %s, want %s", cfg.ReadTimeout, defaultReadTimeout)
	}
	if cfg.GithubAccessToken != "token" {
		t.Errorf("GithubAccessToken = %q, want token", cfg.GithubAccessToken)
	}
}

func TestLoadEnvThenFlagsWin(t *testing.T) {
	t.Setenv("GITHUB_ACCESS_TOKEN", "env-token")
	t.Setenv("APP_ADDR", ":7000")
	t.Setenv("ADMIN_ADDR", ":7001")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("READ_TIMEOUT", "3s")

	cfg, err := Load([]string{
		"-app-addr", ":8000",
		"-log-level", "warn",
		"-github-access-token", "flag-token",
	}, io.Discard)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg.AppAddr != ":8000" {
		t.Errorf("AppAddr = %q, want :8000 (flag)", cfg.AppAddr)
	}
	if cfg.AdminAddr != ":7001" {
		t.Errorf("AdminAddr = %q, want :7001 (env)", cfg.AdminAddr)
	}
	if cfg.LogLevel != "warn" {
		t.Errorf("LogLevel = %q, want warn (flag)", cfg.LogLevel)
	}
	if cfg.ReadTimeout != 3*time.Second {
		t.Errorf("ReadTimeout = %s, want 3s (env)", cfg.ReadTimeout)
	}
	if cfg.GithubAccessToken != "flag-token" {
		t.Errorf("GithubAccessToken = %q, want flag-token", cfg.GithubAccessToken)
	}
}

func TestLoadValidation(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		args []string
	}{
		{
			name: "missing github token",
			env:  map[string]string{"GITHUB_ACCESS_TOKEN": ""},
		},
		{
			name: "same addrs",
			env:  map[string]string{"GITHUB_ACCESS_TOKEN": "token"},
			args: []string{"-app-addr", ":8080", "-admin-addr", ":8080"},
		},
		{
			name: "bad log level",
			env:  map[string]string{"GITHUB_ACCESS_TOKEN": "token"},
			args: []string{"-log-level", "verbose"},
		},
		{
			name: "non-positive timeout",
			env:  map[string]string{"GITHUB_ACCESS_TOKEN": "token"},
			args: []string{"-write-timeout", "0s"},
		},
		{
			name: "bad buckets",
			env:  map[string]string{"GITHUB_ACCESS_TOKEN": "token"},
			args: []string{"-histogram-buckets", "1,0.5"},
		},
		{
			name: "bad redirect url",
			env:  map[string]string{"GITHUB_ACCESS_TOKEN": "token"},
			args: []string{"-root-redirect-url", "not a url"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}
			if _, err := Load(tt.args, io.Discard); err == nil {
				t.Fatal("Load expected error")
			}
		})
	}
}

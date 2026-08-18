// Package version holds build metadata injected via ldflags at link time.
package version

import (
	"fmt"
	"runtime"
)

var (
	// Version is the application semantic version.
	Version = "dev"
	// BuildDate is the UTC build timestamp.
	BuildDate = "unknown"
	// GitCommit is the git SHA at build time.
	GitCommit = "unknown"
	// GitTag is the git tag at build time, if any.
	GitTag = "unknown"
	// GitTreeState is "clean" or "dirty".
	GitTreeState = "unknown"
)

// Info is a snapshot of build metadata for logs, metrics, and /version.
type Info struct {
	Version      string `json:"version"`
	BuildDate    string `json:"build_date"`
	GitCommit    string `json:"git_commit"`
	GitTag       string `json:"git_tag"`
	GitTreeState string `json:"git_tree_state"`
	GoVersion    string `json:"go_version"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

// Current returns build metadata for the running binary.
func Current() Info {
	return Info{
		Version:      Version,
		BuildDate:    BuildDate,
		GitCommit:    GitCommit,
		GitTag:       GitTag,
		GitTreeState: GitTreeState,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String returns a short version line for -v/--version.
func String() string {
	return fmt.Sprintf("aegis %s (%s)", Version, GitCommit)
}

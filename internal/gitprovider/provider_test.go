package gitprovider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGitLabCounts(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/issues"):
			w.Header().Set("X-Total", "7")
		case strings.Contains(r.URL.Path, "/merge_requests"):
			w.Header().Set("X-Total", "4")
		default:
			_ = json.NewEncoder(w).Encode(gitlabProject{StarCount: 12, ForksCount: 3})
		}
	}))
	t.Cleanup(srv.Close)

	client := NewGitLab(srv.Client())
	client.baseURL = srv.URL

	ctx := context.Background()
	forks, err := client.ForkCount(ctx, "owner", "repo")
	if err != nil {
		t.Fatalf("ForkCount: %v", err)
	}
	if forks != 3 {
		t.Errorf("forks = %d, want 3", forks)
	}
	stars, err := client.StarCount(ctx, "owner", "repo")
	if err != nil {
		t.Fatalf("StarCount: %v", err)
	}
	if stars != 12 {
		t.Errorf("stars = %d, want 12", stars)
	}
	issues, err := client.IssueCount(ctx, "owner", "repo", "opened")
	if err != nil {
		t.Fatalf("IssueCount: %v", err)
	}
	if issues != 7 {
		t.Errorf("issues = %d, want 7", issues)
	}
	mrs, err := client.PullRequestCount(ctx, "owner", "repo", "merged")
	if err != nil {
		t.Fatalf("PullRequestCount: %v", err)
	}
	if mrs != 4 {
		t.Errorf("mrs = %d, want 4", mrs)
	}
}

func TestBitbucketCounts(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /repositories/{owner}/{repo}/forks", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(bitbucketSize{Size: 5})
	})
	mux.HandleFunc("GET /repositories/{owner}/{repo}/issues", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(bitbucketSize{Size: 9})
	})
	mux.HandleFunc("GET /repositories/{owner}/{repo}/pullrequests", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(bitbucketSize{Size: 2})
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	client := NewBitbucket(srv.Client())
	client.baseURL = srv.URL

	ctx := context.Background()
	forks, err := client.ForkCount(ctx, "owner", "repo")
	if err != nil {
		t.Fatalf("ForkCount: %v", err)
	}
	if forks != 5 {
		t.Errorf("forks = %d, want 5", forks)
	}
	issues, err := client.IssueCount(ctx, "owner", "repo", "open")
	if err != nil {
		t.Fatalf("IssueCount: %v", err)
	}
	if issues != 9 {
		t.Errorf("issues = %d, want 9", issues)
	}
	prs, err := client.PullRequestCount(ctx, "owner", "repo", "open")
	if err != nil {
		t.Fatalf("PullRequestCount: %v", err)
	}
	if prs != 2 {
		t.Errorf("prs = %d, want 2", prs)
	}
	stars, err := client.StarCount(ctx, "owner", "repo")
	if err != nil {
		t.Fatalf("StarCount: %v", err)
	}
	if stars != -2 {
		t.Errorf("stars = %d, want -2", stars)
	}
}

func TestGitLabUnexpectedStatus(t *testing.T) {
	t.Parallel()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(srv.Close)

	client := NewGitLab(srv.Client())
	client.baseURL = srv.URL
	if _, err := client.ForkCount(context.Background(), "o", "r"); err == nil {
		t.Fatal("expected error")
	}
}

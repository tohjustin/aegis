package gitprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

const gitlabAPI = "https://gitlab.com/api/v4"

// GitLab fetches repository metrics from the GitLab REST API.
type GitLab struct {
	baseURL string
	client  *http.Client
}

type gitlabProject struct {
	StarCount  int `json:"star_count"`
	ForksCount int `json:"forks_count"`
}

// NewGitLab returns a GitLab client. If httpClient is nil, a client with a 10s timeout is used.
func NewGitLab(httpClient *http.Client) *GitLab {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &GitLab{baseURL: gitlabAPI, client: httpClient}
}

// ForkCount returns the number of forks for a project.
func (g *GitLab) ForkCount(ctx context.Context, owner, repo string) (int, error) {
	project, err := g.project(ctx, owner, repo)
	if err != nil {
		return 0, err
	}
	return project.ForksCount, nil
}

// IssueCount returns the number of issues for a project, optionally filtered by state.
func (g *GitLab) IssueCount(ctx context.Context, owner, repo, issueState string) (int, error) {
	url := fmt.Sprintf("%s/projects/%s%%2F%s/issues", g.baseURL, owner, repo)
	switch issueState {
	case "opened", "closed":
		url = fmt.Sprintf("%s?state=%s", url, issueState)
	}
	return g.xTotal(ctx, url)
}

// PullRequestCount returns the number of merge requests for a project, optionally filtered by state.
func (g *GitLab) PullRequestCount(ctx context.Context, owner, repo, pullRequestState string) (int, error) {
	url := fmt.Sprintf("%s/projects/%s%%2F%s/merge_requests", g.baseURL, owner, repo)
	switch pullRequestState {
	case "opened", "closed", "locked", "merged":
		url = fmt.Sprintf("%s?state=%s", url, pullRequestState)
	}
	return g.xTotal(ctx, url)
}

// StarCount returns the number of stars for a project.
func (g *GitLab) StarCount(ctx context.Context, owner, repo string) (int, error) {
	project, err := g.project(ctx, owner, repo)
	if err != nil {
		return 0, err
	}
	return project.StarCount, nil
}

func (g *GitLab) project(ctx context.Context, owner, repo string) (*gitlabProject, error) {
	url := fmt.Sprintf("%s/projects/%s%%2F%s", g.baseURL, owner, repo)
	resp, err := g.get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var project gitlabProject
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, fmt.Errorf("gitlab decode project: %w", err)
	}
	return &project, nil
}

func (g *GitLab) xTotal(ctx context.Context, url string) (int, error) {
	resp, err := g.get(ctx, url)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	total, err := strconv.Atoi(resp.Header.Get("X-Total"))
	if err != nil {
		return 0, fmt.Errorf("gitlab X-Total: %w", err)
	}
	return total, nil
}

func (g *GitLab) get(ctx context.Context, url string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("gitlab request: %w", err)
	}
	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gitlab fetch: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		return nil, fmt.Errorf("gitlab: unexpected status %d", resp.StatusCode)
	}
	return resp, nil
}

package gitprovider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const bitbucketAPI = "https://api.bitbucket.org/2.0"

// Bitbucket fetches repository metrics from the Bitbucket Cloud REST API.
type Bitbucket struct {
	baseURL string
	client  *http.Client
}

type bitbucketSize struct {
	Size int `json:"size"`
}

// NewBitbucket returns a Bitbucket client. If httpClient is nil, a client with a 10s timeout is used.
func NewBitbucket(httpClient *http.Client) *Bitbucket {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &Bitbucket{baseURL: bitbucketAPI, client: httpClient}
}

// ForkCount returns the number of forks for a repository.
func (b *Bitbucket) ForkCount(ctx context.Context, owner, repo string) (int, error) {
	url := fmt.Sprintf("%s/repositories/%s/%s/forks?fields=size", b.baseURL, owner, repo)
	return b.size(ctx, url)
}

// IssueCount returns the number of issues for a repository, optionally filtered by state.
func (b *Bitbucket) IssueCount(ctx context.Context, owner, repo, issueState string) (int, error) {
	url := fmt.Sprintf("%s/repositories/%s/%s/issues", b.baseURL, owner, repo)
	switch issueState {
	case "new", "open", "resolved", "invalid", "duplicate", "wontfix", "closed":
		url = fmt.Sprintf("%s?fields=size&q=(state+=+\"%s\")", url, issueState)
	case "on-hold":
		url = fmt.Sprintf("%s?fields=size&q=(state+=+\"on%%20hold\")", url)
	}
	return b.size(ctx, url)
}

// PullRequestCount returns the number of pull requests for a repository, optionally filtered by state.
func (b *Bitbucket) PullRequestCount(ctx context.Context, owner, repo, pullRequestState string) (int, error) {
	url := fmt.Sprintf("%s/repositories/%s/%s/pullrequests", b.baseURL, owner, repo)
	switch pullRequestState {
	case "merged", "superseded", "open", "declined":
		url = fmt.Sprintf("%s?fields=size&q=(state+=+\"%s\")", url, pullRequestState)
	}
	return b.size(ctx, url)
}

// StarCount is unimplemented for Bitbucket Cloud and returns -2.
func (b *Bitbucket) StarCount(context.Context, string, string) (int, error) {
	return -2, nil
}

func (b *Bitbucket) size(ctx context.Context, url string) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return 0, fmt.Errorf("bitbucket request: %w", err)
	}
	resp, err := b.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("bitbucket fetch: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("bitbucket: unexpected status %d", resp.StatusCode)
	}

	var body bitbucketSize
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return 0, fmt.Errorf("bitbucket decode: %w", err)
	}
	return body.Size, nil
}

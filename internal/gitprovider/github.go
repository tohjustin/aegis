package gitprovider

import (
	"context"
	"fmt"

	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// GitHub fetches repository metrics from the GitHub GraphQL API.
type GitHub struct {
	client *githubv4.Client
}

// NewGitHub returns a GitHub client authenticated with the given access token.
func NewGitHub(accessToken string) (*GitHub, error) {
	if accessToken == "" {
		return nil, fmt.Errorf("github access token is required")
	}
	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	httpClient := oauth2.NewClient(context.Background(), src)
	return &GitHub{client: githubv4.NewClient(httpClient)}, nil
}

// ForkCount returns the number of forks for a repository.
func (g *GitHub) ForkCount(ctx context.Context, owner, repo string) (int, error) {
	var query struct {
		Repository struct {
			Forks struct {
				TotalCount int
			}
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}
	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"repo":  githubv4.String(repo),
	}
	if err := g.client.Query(ctx, &query, variables); err != nil {
		return 0, fmt.Errorf("github forks: %w", err)
	}
	return query.Repository.Forks.TotalCount, nil
}

// IssueCount returns the number of issues for a repository, optionally filtered by state.
func (g *GitHub) IssueCount(ctx context.Context, owner, repo, issueState string) (int, error) {
	var issueStates []githubv4.IssueState
	switch issueState {
	case "open":
		issueStates = []githubv4.IssueState{githubv4.IssueStateOpen}
	case "closed":
		issueStates = []githubv4.IssueState{githubv4.IssueStateClosed}
	default:
		issueStates = []githubv4.IssueState{
			githubv4.IssueStateOpen,
			githubv4.IssueStateClosed,
		}
	}
	var query struct {
		Repository struct {
			Issues struct {
				TotalCount int
			} `graphql:"issues(states: $states)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}
	variables := map[string]interface{}{
		"owner":  githubv4.String(owner),
		"repo":   githubv4.String(repo),
		"states": issueStates,
	}
	if err := g.client.Query(ctx, &query, variables); err != nil {
		return 0, fmt.Errorf("github issues: %w", err)
	}
	return query.Repository.Issues.TotalCount, nil
}

// PullRequestCount returns the number of pull requests for a repository, optionally filtered by state.
func (g *GitHub) PullRequestCount(ctx context.Context, owner, repo, pullRequestState string) (int, error) {
	var pullRequestStates []githubv4.PullRequestState
	switch pullRequestState {
	case "open":
		pullRequestStates = []githubv4.PullRequestState{githubv4.PullRequestStateOpen}
	case "closed":
		pullRequestStates = []githubv4.PullRequestState{githubv4.PullRequestStateClosed}
	case "merged":
		pullRequestStates = []githubv4.PullRequestState{githubv4.PullRequestStateMerged}
	default:
		pullRequestStates = []githubv4.PullRequestState{
			githubv4.PullRequestStateOpen,
			githubv4.PullRequestStateClosed,
			githubv4.PullRequestStateMerged,
		}
	}
	var query struct {
		Repository struct {
			PullRequests struct {
				TotalCount int
			} `graphql:"pullRequests(states: $states)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}
	variables := map[string]interface{}{
		"owner":  githubv4.String(owner),
		"repo":   githubv4.String(repo),
		"states": pullRequestStates,
	}
	if err := g.client.Query(ctx, &query, variables); err != nil {
		return 0, fmt.Errorf("github pull requests: %w", err)
	}
	return query.Repository.PullRequests.TotalCount, nil
}

// StarCount returns the number of stars for a repository.
func (g *GitHub) StarCount(ctx context.Context, owner, repo string) (int, error) {
	var query struct {
		Repository struct {
			Stargazers struct {
				TotalCount int
			}
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}
	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"repo":  githubv4.String(repo),
	}
	if err := g.client.Query(ctx, &query, variables); err != nil {
		return 0, fmt.Errorf("github stars: %w", err)
	}
	return query.Repository.Stargazers.TotalCount, nil
}

package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shurcooL/githubv4"
	"go.uber.org/zap"
	"golang.org/x/oauth2"

	"github.com/tohjustin/aegis/pkg/badge"
	"github.com/tohjustin/aegis/service/config"
)

type githubService struct {
	name   string
	client *githubv4.Client
	config *config.Config
	logger *zap.Logger
}

// NewGithubService returns a HTTP handler for the Github badge service
func NewGithubService(configuration *config.Config,
	logger *zap.Logger) (GitProviderService, error) {
	if configuration == nil {
		return nil, fmt.Errorf("missing config dependency")
	}
	if logger == nil {
		return nil, fmt.Errorf("missing logger dependency")
	}

	accessToken := configuration.GithubAccessToken
	if accessToken == "" {
		return nil, fmt.Errorf("missing GitHub access token")
	}

	// Create new Github GraphQL client
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)

	return &githubService{
		name:   "github",
		client: githubv4.NewClient(httpClient),
		config: configuration,
		logger: logger,
	}, nil
}

func (service *githubService) getForkCount(owner string, repo string) (int, error) {
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

	err := service.client.Query(context.Background(), &query, variables)
	return query.Repository.Forks.TotalCount, err
}

func (service *githubService) getIssueCount(owner string, repo string, issueState string) (int, error) {
	var issueStates []githubv4.IssueState
	var query struct {
		Repository struct {
			Issues struct {
				TotalCount int
			} `graphql:"issues(states: $states)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}
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
	variables := map[string]interface{}{
		"owner":  githubv4.String(owner),
		"repo":   githubv4.String(repo),
		"states": issueStates,
	}

	err := service.client.Query(context.Background(), &query, variables)
	return query.Repository.Issues.TotalCount, err
}

func (service *githubService) getPullRequestCount(owner string, repo string, pullRequestState string) (int, error) {
	var pullRequestStates []githubv4.PullRequestState
	var query struct {
		Repository struct {
			PullRequests struct {
				TotalCount int
			} `graphql:"pullRequests(states: $states)"`
		} `graphql:"repository(owner: $owner, name: $repo)"`
	}
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
	variables := map[string]interface{}{
		"owner":  githubv4.String(owner),
		"repo":   githubv4.String(repo),
		"states": pullRequestStates,
	}

	err := service.client.Query(context.Background(), &query, variables)
	return query.Repository.PullRequests.TotalCount, err
}

func (service *githubService) getStarCount(owner string, repo string) (int, error) {
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

	err := service.client.Query(context.Background(), &query, variables)
	return query.Repository.Stargazers.TotalCount, err
}

func (service *githubService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	routeVariables := mux.Vars(r)
	owner := routeVariables["owner"]
	repo := routeVariables["repo"]
	method := routeVariables["method"]

	// Fetch data
	var color, status, subject string
	var value int
	var err error
	switch method {
	case "forks":
		subject = "forks"
		value, err = service.getForkCount(owner, repo)
	case "issues":
		state := r.URL.Query().Get("state")
		switch state {
		case "":
			subject = "issues"
		case "open":
			subject = "open issues"
		case "closed":
			subject = "closed issues"
		default:
			service.logger.Info("Unsupported state",
				zap.String("url", r.URL.RequestURI()),
				zap.String("service", service.name),
				zap.String("method", method),
				zap.String("state", state))
			if err := badRequest(w, service.config); err != nil {
				service.logger.Error("Failed to create error badge",
					zap.String("url", r.URL.RequestURI()),
					zap.String("service", service.name),
					zap.String("method", method),
					zap.Error(err))
			}
			return
		}
		value, err = service.getIssueCount(owner, repo, state)
	case "pull-requests":
		state := r.URL.Query().Get("state")
		switch state {
		case "":
			subject = "PRs"
		case "open":
			subject = "open PRs"
		case "closed":
			subject = "closed PRs"
		case "merged":
			subject = "merged PRs"
		default:
			service.logger.Info("Unsupported state",
				zap.String("url", r.URL.RequestURI()),
				zap.String("service", service.name),
				zap.String("method", method),
				zap.String("state", state))
			if err := badRequest(w, service.config); err != nil {
				service.logger.Error("Failed to create error badge",
					zap.String("url", r.URL.RequestURI()),
					zap.String("service", service.name),
					zap.String("method", method),
					zap.Error(err))
			}
			return
		}
		value, err = service.getPullRequestCount(owner, repo, state)
	case "stars":
		subject = "stars"
		value, err = service.getStarCount(owner, repo)
	default:
		service.logger.Info("Unsupported method",
			zap.String("url", r.URL.RequestURI()),
			zap.String("service", service.name),
			zap.String("method", method))
		if err := notFound(w, service.config); err != nil {
			service.logger.Error("Failed to create error badge",
				zap.String("url", r.URL.RequestURI()),
				zap.String("service", service.name),
				zap.String("method", method),
				zap.Error(err))
		}
		return
	}
	if err != nil {
		service.logger.Error("Failed to fetch data",
			zap.String("url", r.URL.RequestURI()),
			zap.String("service", service.name),
			zap.String("method", method),
			zap.Error(err))
		if err := internalServerError(w, service.config); err != nil {
			service.logger.Error("Failed to create error badge",
				zap.String("url", r.URL.RequestURI()),
				zap.String("service", service.name),
				zap.String("method", method),
				zap.Error(err))
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	status = formatIntegerWithMetricPrefix(value)

	// Overwrite any badge texts
	if queryColor := r.URL.Query().Get("color"); queryColor != "" {
		color = queryColor
	}
	if queryStatus := r.URL.Query().Get("status"); queryStatus != "" {
		status = queryStatus
	}
	if querySubject := r.URL.Query().Get("subject"); querySubject != "" {
		subject = querySubject
	}

	// Generate badge
	generatedBadge, err := badge.Create(&badge.Params{
		Style:   badge.Style(r.URL.Query().Get("style")),
		Subject: subject,
		Status:  status,
		Color:   color,
		Icon:    r.URL.Query().Get("icon"),
	})
	if err != nil {
		service.logger.Error("Failed to create badge",
			zap.String("url", r.URL.RequestURI()),
			zap.String("service", service.name),
			zap.String("method", method),
			zap.Error(err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !service.config.ExcludeCacheControlHeaders {
		// cache response in browser for 1 hour (3600), CDN for 1 hour (3600)
		w.Header().Set("Cache-Control", "public, max-age=3600, s-maxage=3600")
	}
	w.Header().Set("Content-Type", "image/svg+xml;utf-8")
	_, err = w.Write([]byte(generatedBadge))
	service.logger.Error("Failed to write HTTP response",
		zap.String("url", r.URL.RequestURI()),
		zap.String("service", service.name),
		zap.String("method", method),
		zap.Error(err))
}

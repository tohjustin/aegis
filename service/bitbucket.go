package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/tohjustin/aegis/pkg/badge"
	"github.com/tohjustin/aegis/service/config"
)

type bitbucketService struct {
	name   string
	config *config.Config
	logger *zap.Logger
}

type bitbucketFilteredResponse struct {
	Size int `json:"size"`
}

// NewBitbucketService returns a HTTP handler for the Bitbucket badge service
func NewBitbucketService(configuration *config.Config,
	logger *zap.Logger) (GitProviderService, error) {
	if configuration == nil {
		return nil, fmt.Errorf("missing config dependency")
	}
	if logger == nil {
		return nil, fmt.Errorf("missing logger dependency")
	}

	return &bitbucketService{
		name:   "bitbucket",
		config: configuration,
		logger: logger,
	}, nil
}

func (service *bitbucketService) fetch(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func (service *bitbucketService) getForkCount(owner string, repo string) (int, error) {
	url := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/forks?&fields=size", owner, repo)
	resp, err := service.fetch(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var forks bitbucketFilteredResponse
	if err := json.NewDecoder(resp.Body).Decode(&forks); err != nil {
		return -1, err
	}

	return forks.Size, nil
}

func (service *bitbucketService) getIssueCount(owner string, repo string, issueState string) (int, error) {
	url := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/issues", owner, repo)
	switch issueState {
	case "new":
		url = fmt.Sprintf("%s?&fields=size&q=(state+=+\"%s\")", url, issueState)
	case "open":
		url = fmt.Sprintf("%s?&fields=size&q=(state+=+\"%s\")", url, issueState)
	case "resolved":
		url = fmt.Sprintf("%s?&fields=size&q=(state+=+\"%s\")", url, issueState)
	case "on-hold":
		url = fmt.Sprintf("%s?&fields=size&q=(state+=+\"on%%20hold\")", url)
	case "invalid":
		url = fmt.Sprintf("%s?&fields=size&q=(state+=+\"%s\")", url, issueState)
	case "duplicate":
		url = fmt.Sprintf("%s?&fields=size&q=(state+=+\"%s\")", url, issueState)
	case "wontfix":
		url = fmt.Sprintf("%s?&fields=size&q=(state+=+\"%s\")", url, issueState)
	case "closed":
		url = fmt.Sprintf("%s?&fields=size&q=(state+=+\"%s\")", url, issueState)
	}
	resp, err := service.fetch(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var issues bitbucketFilteredResponse
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		return 0, err
	}

	return issues.Size, nil
}

func (service *bitbucketService) getPullRequestCount(owner string, repo string, pullRequestState string) (int, error) {
	url := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/pullrequests", owner, repo)
	switch pullRequestState {
	case "merged":
		url = fmt.Sprintf("%s?&fields=size&q=(state+=+\"%s\")", url, pullRequestState)
	case "superseded":
		url = fmt.Sprintf("%s?&fields=size&q=(state+=+\"%s\")", url, pullRequestState)
	case "open":
		url = fmt.Sprintf("%s?&fields=size&q=(state+=+\"%s\")", url, pullRequestState)
	case "declined":
		url = fmt.Sprintf("%s?&fields=size&q=(state+=+\"%s\")", url, pullRequestState)
	}
	resp, err := service.fetch(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	var pullRequests bitbucketFilteredResponse
	if err := json.NewDecoder(resp.Body).Decode(&pullRequests); err != nil {
		return 0, err
	}

	return pullRequests.Size, nil
}

func (service *bitbucketService) getStarCount(owner string, repo string) (int, error) {
	return -2, nil
}

func (service *bitbucketService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		case "new":
			subject = "new issues"
		case "open":
			subject = "open issues"
		case "resolved":
			subject = "resolved issues"
		case "on-hold":
			subject = "on-hold issues"
		case "invalid":
			subject = "invalid issues"
		case "duplicate":
			subject = "duplicate issues"
		case "wontfix":
			subject = "wontfix issues"
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
		case "merged":
			subject = "merged PRs"
		case "superseded":
			subject = "superseded PRs"
		case "open":
			subject = "open PRs"
		case "declined":
			subject = "declined PRs"
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

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tohjustin/badger/pkg/badge"
)

type BitbucketFilteredResponse struct {
	Size int `json:"size"`
}

func NewBitbucketService() RepositoryService {
	return &bitbucketService{}
}

type bitbucketService struct{}

func (service *bitbucketService) getForkCount(owner string, repo string) (int, error) {
	url := fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s/%s/forks?&fields=size", owner, repo)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return 0, err
	}
	defer resp.Body.Close()

	var forks BitbucketFilteredResponse
	if err := json.NewDecoder(resp.Body).Decode(&forks); err != nil {
		log.Println(err)
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
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return 0, err
	}
	defer resp.Body.Close()

	var issues BitbucketFilteredResponse
	if err := json.NewDecoder(resp.Body).Decode(&issues); err != nil {
		log.Println(err)
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
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return 0, err
	}
	defer resp.Body.Close()

	var pullRequests BitbucketFilteredResponse
	if err := json.NewDecoder(resp.Body).Decode(&pullRequests); err != nil {
		log.Println(err)
		return 0, err
	}

	return pullRequests.Size, nil
}

func (service *bitbucketService) getStargazerCount(owner string, repo string) (int, error) {
	return -2, nil
}

func (service *bitbucketService) Handler(w http.ResponseWriter, r *http.Request) {
	routeVariables := mux.Vars(r)
	owner := routeVariables["owner"]
	repo := routeVariables["repo"]
	requestType := routeVariables["requestType"]

	// Fetch data
	var color, status, subject string
	var value int
	var err error
	switch requestType {
	case "forks":
		subject = "forks"
		value, err = service.getForkCount(owner, repo)
	case "issues":
		state := r.URL.Query().Get("state")
		switch state {
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
			subject = "issues"
		}
		value, err = service.getIssueCount(owner, repo, state)
	case "pull-requests":
		state := r.URL.Query().Get("state")
		switch state {
		case "merged":
			subject = "merged PRs"
		case "superseded":
			subject = "superseded PRs"
		case "open":
			subject = "open PRs"
		case "declined":
			subject = "declined PRs"
		default:
			subject = "PRs"
		}
		value, err = service.getPullRequestCount(owner, repo, state)
	case "stars":
		subject = "stars"
		value, err = service.getStargazerCount(owner, repo)
	}

	// Compute status
	if err != nil {
		status = err.Error()
	} else {
		status = strconv.Itoa(value)
	}

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
	icon := r.URL.Query().Get("icon")
	style := r.URL.Query().Get("style")

	// Generate badge
	createOptions := badge.Options{Color: color, Icon: icon, Style: badge.Style(style)}
	generatedBadge, err := badge.Create(subject, status, &createOptions)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	// cache response in browser for 1 hour (3600), CDN for 1 hour (3600)
	w.Header().Set("Cache-Control", "public, max-age=3600, s-maxage=3600")
	w.Header().Set("Content-Type", "image/svg+xml;utf-8")
	w.Write([]byte(generatedBadge))
}

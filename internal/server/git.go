package server

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/tohjustin/aegis/internal/gitprovider"
	"github.com/tohjustin/aegis/pkg/badge"
)

type countFunc func(ctx context.Context, owner, repo, state string) (int, error)

type gitMethod struct {
	subjects map[string]string
	fetch    countFunc
}

func githubMethods(g *gitprovider.GitHub) map[string]gitMethod {
	return map[string]gitMethod{
		"forks": {
			subjects: map[string]string{"": "forks"},
			fetch: func(ctx context.Context, owner, repo, _ string) (int, error) {
				return g.ForkCount(ctx, owner, repo)
			},
		},
		"issues": {
			subjects: map[string]string{
				"":       "issues",
				"open":   "open issues",
				"closed": "closed issues",
			},
			fetch: g.IssueCount,
		},
		"pull-requests": {
			subjects: map[string]string{
				"":       "PRs",
				"open":   "open PRs",
				"closed": "closed PRs",
				"merged": "merged PRs",
			},
			fetch: g.PullRequestCount,
		},
		"stars": {
			subjects: map[string]string{"": "stars"},
			fetch: func(ctx context.Context, owner, repo, _ string) (int, error) {
				return g.StarCount(ctx, owner, repo)
			},
		},
	}
}

func gitlabMethods(g *gitprovider.GitLab) map[string]gitMethod {
	return map[string]gitMethod{
		"forks": {
			subjects: map[string]string{"": "forks"},
			fetch: func(ctx context.Context, owner, repo, _ string) (int, error) {
				return g.ForkCount(ctx, owner, repo)
			},
		},
		"issues": {
			subjects: map[string]string{
				"":       "issues",
				"opened": "opened issues",
				"closed": "closed issues",
			},
			fetch: g.IssueCount,
		},
		"merge-requests": {
			subjects: map[string]string{
				"":       "MRs",
				"opened": "opened MRs",
				"closed": "closed MRs",
				"locked": "locked MRs",
				"merged": "merged MRs",
			},
			fetch: g.PullRequestCount,
		},
		"stars": {
			subjects: map[string]string{"": "stars"},
			fetch: func(ctx context.Context, owner, repo, _ string) (int, error) {
				return g.StarCount(ctx, owner, repo)
			},
		},
	}
}

func bitbucketMethods(b *gitprovider.Bitbucket) map[string]gitMethod {
	return map[string]gitMethod{
		"forks": {
			subjects: map[string]string{"": "forks"},
			fetch: func(ctx context.Context, owner, repo, _ string) (int, error) {
				return b.ForkCount(ctx, owner, repo)
			},
		},
		"issues": {
			subjects: map[string]string{
				"":          "issues",
				"new":       "new issues",
				"open":      "open issues",
				"resolved":  "resolved issues",
				"on-hold":   "on-hold issues",
				"invalid":   "invalid issues",
				"duplicate": "duplicate issues",
				"wontfix":   "wontfix issues",
				"closed":    "closed issues",
			},
			fetch: b.IssueCount,
		},
		"pull-requests": {
			subjects: map[string]string{
				"":           "PRs",
				"merged":     "merged PRs",
				"superseded": "superseded PRs",
				"open":       "open PRs",
				"declined":   "declined PRs",
			},
			fetch: b.PullRequestCount,
		},
		"stars": {
			subjects: map[string]string{"": "stars"},
			fetch: func(ctx context.Context, owner, repo, _ string) (int, error) {
				return b.StarCount(ctx, owner, repo)
			},
		},
	}
}

func (a *app) handleGit(service string, methods map[string]gitMethod) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		owner := r.PathValue("owner")
		repo := r.PathValue("repo")
		method := r.PathValue("method")

		spec, ok := methods[method]
		if !ok {
			a.logger.LogAttrs(r.Context(), slog.LevelInfo, "unsupported method",
				slog.String("url", r.URL.RequestURI()),
				slog.String("service", service),
				slog.String("method", method),
			)
			if err := a.writeErrorBadge(w, "not found"); err != nil {
				a.logger.LogAttrs(r.Context(), slog.LevelError, "failed to create error badge",
					slog.String("url", r.URL.RequestURI()),
					slog.String("service", service),
					slog.String("method", method),
					slog.Any("err", err),
				)
			}
			return
		}

		state := r.URL.Query().Get("state")
		subject, ok := spec.subjects[state]
		if !ok {
			a.logger.LogAttrs(r.Context(), slog.LevelInfo, "unsupported state",
				slog.String("url", r.URL.RequestURI()),
				slog.String("service", service),
				slog.String("method", method),
				slog.String("state", state),
			)
			if err := a.writeErrorBadge(w, "bad request"); err != nil {
				a.logger.LogAttrs(r.Context(), slog.LevelError, "failed to create error badge",
					slog.String("url", r.URL.RequestURI()),
					slog.String("service", service),
					slog.String("method", method),
					slog.Any("err", err),
				)
			}
			return
		}

		value, err := spec.fetch(r.Context(), owner, repo, state)
		if err != nil {
			a.logger.LogAttrs(r.Context(), slog.LevelError, "failed to fetch data",
				slog.String("url", r.URL.RequestURI()),
				slog.String("service", service),
				slog.String("method", method),
				slog.Any("err", err),
			)
			if err := a.writeErrorBadge(w, "internal server error"); err != nil {
				a.logger.LogAttrs(r.Context(), slog.LevelError, "failed to create error badge",
					slog.String("url", r.URL.RequestURI()),
					slog.String("service", service),
					slog.String("method", method),
					slog.Any("err", err),
				)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}

		status := gitprovider.FormatCount(value)
		color := r.URL.Query().Get("color")
		if queryStatus := r.URL.Query().Get("status"); queryStatus != "" {
			status = queryStatus
		}
		if querySubject := r.URL.Query().Get("subject"); querySubject != "" {
			subject = querySubject
		}

		a.writeBadge(w, r, &badge.Params{
			Style:   badge.Style(r.URL.Query().Get("style")),
			Subject: subject,
			Status:  status,
			Color:   color,
			Icon:    r.URL.Query().Get("icon"),
		})
	}
}

package server

import (
	"net/http"

	"github.com/tohjustin/aegis/pkg/badge"
)

func (a *app) handleStatic(w http.ResponseWriter, r *http.Request) {
	a.writeBadge(w, r, &badge.Params{
		Style:   badge.Style(r.URL.Query().Get("style")),
		Subject: r.URL.Query().Get("subject"),
		Status:  r.URL.Query().Get("status"),
		Color:   r.URL.Query().Get("color"),
		Icon:    r.URL.Query().Get("icon"),
	})
}

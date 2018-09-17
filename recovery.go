package main

import (
	"fmt"
	"net/http"

	"github.com/urfave/negroni"
)

type textPanicFormatter struct{}

func (t *textPanicFormatter) FormatPanicError(
	rw http.ResponseWriter,
	r *http.Request,
	infos *negroni.PanicInformation,
) {
	if rw.Header().Get("Content-Type") == "" {
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
	}
	fmt.Fprintf(rw, "500 internal server error")
}

func newRecoveryMiddleware() (recovery *negroni.Recovery) {
	recovery = negroni.NewRecovery()
	recovery.Formatter = &textPanicFormatter{}
	recovery.PrintStack = false
	return
}

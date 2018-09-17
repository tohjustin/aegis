package main

import (
	"github.com/urfave/negroni"
)

func newLoggerMiddleware() (logger *negroni.Logger) {
	logger = negroni.NewLogger()
	return
}

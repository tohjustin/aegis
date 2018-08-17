package main

import (
	"log"
	"log/syslog"

	"github.com/urfave/negroni"
)

func newLoggerMiddleware(logEndpoint string) (logger *negroni.Logger) {
	logger = negroni.NewLogger()
	if logEndpoint != "" {
		w, err := syslog.Dial("udp", logEndpoint, 0, "badger-server")
		if err != nil {
			log.Fatal("Failed to dial syslog at: ", logEndpoint)
			return
		}

		logger.ALogger = log.New(w, "[negroni] ", 0)
	}

	return
}

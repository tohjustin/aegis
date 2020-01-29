package main

import (
	"log"

	"github.com/tohjustin/aegis/internal/version"
	"github.com/tohjustin/aegis/service"
)

func main() {
	handleErr := func(err error) {
		if err != nil {
			log.Fatalf("Failed to run the service: %v", err)
		}
	}

	info := service.Info{
		ExecutableName: "aegis",
		ShortName:      "Aegis",
		LongName:       "Aegis badge generation service",
		Version:        version.Version,
		GitHash:        version.GitHash,
	}

	svc, err := service.New(info)
	handleErr(err)

	err = svc.Start()
	handleErr(err)
}

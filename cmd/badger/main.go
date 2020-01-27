package main

import (
	"log"

	"github.com/tohjustin/badger/service"
)

var (
	// Version represents the application semantic version, variable will be replaced at link time after `make` has been run.
	Version = "latest"
	// GitHash represents the application Git SHA-1 hash, variable will be replaced at link time after `make` has been run.
	GitHash = "<UNKNOWN>"
)

func main() {
	info := service.Info{
		ExecutableName: "badger",
		LongName:       "Badger Badge Generation Service",
		Version:        Version,
		GitHash:        GitHash,
	}

	svc, err := service.New(info)
	if err != nil {
		log.Fatalf("Failed to run the service: %v", err)
	}

	svc.Start()
}

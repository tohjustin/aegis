package main

import (
	"log"

	"github.com/tohjustin/badger/internal/version"
	"github.com/tohjustin/badger/service"
)

func main() {
	info := service.Info{
		ExecutableName: "badger",
		LongName:       "Badger Badge Generation Service",
		Version:        version.Version,
		GitHash:        version.GitHash,
	}

	svc, err := service.New(info)
	if err != nil {
		log.Fatalf("Failed to run the service: %v", err)
	}

	svc.Start()
}

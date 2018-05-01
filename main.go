package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	svgBadge, err := generateClassicBadge("Hello World", "100%", "#1896DE")
	if err != nil {
		fmt.Fprint(w, err)
	} else {
		w.Header().Set("Content-Type", "image/svg+xml;utf-8")
		fmt.Fprint(w, svgBadge)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

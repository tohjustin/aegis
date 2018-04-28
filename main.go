package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
)

const badgeTemplate = `
<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="116" height="20">
	<linearGradient id="smooth" x2="0" y2="100%">
		<stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
		<stop offset="1" stop-opacity=".1"/>
	</linearGradient>
	<clipPath id="round">
		<rect width="116" height="20" rx="3" fill="#fff"/>
	</clipPath>
	<g clip-path="url(#round)">
		<rect width="73" height="20" fill="#555"/>
		<rect x="73" width="80" height="20" fill="{{.Color}}"/>
		<rect width="180" height="20" fill="url(#smooth)"/>>
	</g>
	<g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="110">
		<text x="375" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="630">{{.Subject}}</text>
		<text x="375" y="140" transform="scale(.1)" textLength="630">{{.Subject}}</text>
		<text x="935" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="330">{{.Status}}</text>
		<text x="935" y="140" transform="scale(.1)" textLength="330">{{.Status}}</text>
	</g>
</svg>
`

type badge struct {
	Subject string
	Status  string
	Color   string
}

func generateBadge(subject string, status string, color string) string {
	helloWorldBadge := badge{
		Subject: subject,
		Status:  status,
		Color:   color,
	}

	var svgBuffer bytes.Buffer
	t := template.New("")
	t.Parse(badgeTemplate)
	t.Execute(&svgBuffer, helloWorldBadge)

	return svgBuffer.String()
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml;utf-8")
	fmt.Fprint(w, generateBadge("Hello World", "100%", "#1896DE"))
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

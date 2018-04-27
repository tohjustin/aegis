package main

import (
	"fmt"
	"log"
	"net/http"
)

const badge = `
<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="116" height="20">
	<linearGradient id="b" x2="0" y2="100%">
		<stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
		<stop offset="1" stop-opacity=".1"/>
	</linearGradient>
	<clipPath id="a">
		<rect width="116" height="20" rx="3" fill="#fff"/>
	</clipPath>
	<g clip-path="url(#a)">
		<path fill="#555" d="M0 0h73v20H0z"/>
		<path fill="#1896DE" d="M73 0h43v20H73z"/>
		<path fill="url(#b)" d="M0 0h116v20H0z"/>
	</g>
	<g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="110">
		<text x="375" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="630">Hello World</text>
		<text x="375" y="140" transform="scale(.1)" textLength="630">Hello World</text>
		<text x="935" y="150" fill="#010101" fill-opacity=".3" transform="scale(.1)" textLength="330">100%</text>
		<text x="935" y="140" transform="scale(.1)" textLength="330">100%</text>
	</g>
</svg>
`

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml;utf-8")
	fmt.Fprint(w, badge)
}

func main() {
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

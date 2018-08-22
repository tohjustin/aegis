package main

import (
	"net/http"
	"testing"

	"github.com/tohjustin/badger/pkg/badge"
)

func TestBadgeServiceHandler(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/badge/testSubject/testStatus/ff0000",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody:   badge.GenerateSVGUnsafe("classic", "testSubject", "testStatus", "ff0000", ""),
	})
}

func TestBadgeServiceHandlerWithCSSColorName(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/badge/testSubject/testStatus/red",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody:   badge.GenerateSVGUnsafe("classic", "testSubject", "testStatus", "red", ""),
	})
}

func TestBadgeServiceHandlerWithNoColor(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/badge/testSubject/testStatus",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody:   badge.GenerateSVGUnsafe("classic", "testSubject", "testStatus", "", ""),
	})
}

func TestBadgeServiceHandlerWithBadColor(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/badge/testSubject/testStatus/badColor",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody:   badge.GenerateSVGUnsafe("classic", "testSubject", "testStatus", "", ""),
	})
}

func TestBadgeServiceHandlerWithIconQuery(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/badge/testSubject/testStatus/ff0000?icon=brands/docker",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody:   badge.GenerateSVGUnsafe("classic", "testSubject", "testStatus", "ff0000", "brands/docker"),
	})
}

func TestBadgeServiceHandlerWithBadIconQuery(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/badge/testSubject/testStatus/ff0000?icon=badIcon",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody:   badge.GenerateSVGUnsafe("classic", "testSubject", "testStatus", "ff0000", ""),
	})
}

func TestBadgeServiceHandlerWithStyleQuery(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/badge/testSubject/testStatus/ff0000?style=semaphore",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody:   badge.GenerateSVGUnsafe("semaphore", "testSubject", "testStatus", "ff0000", ""),
	})
}

func TestBadgeServiceHandlerWithBadStyleQuery(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/badge/testSubject/testStatus/ff0000?style=badStyle",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody:   badge.GenerateSVGUnsafe("classic", "testSubject", "testStatus", "ff0000", ""),
	})
}

func TestBadgeServiceHandlerWithBadHTTPMethods(t *testing.T) {
	t.Parallel()

	// service only accepts "GET" requests
	badHTTPMethods := []string{
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}

	for _, badHTTPMethod := range badHTTPMethods {
		runHTTPTest(t, httpTestCase{
			requestMethod:   badHTTPMethod,
			requestPath:     "/badge/testSubject/testStatus/ff0000",
			expectedHeaders: nil,
			expectedStatus:  405,
			expectedBody:    "",
		})
	}
}

func TestBadgeServiceErrorHandler(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody:   badge.GenerateSVGUnsafe("classic", "badger", "400 Bad Request", "", ""),
	})
}

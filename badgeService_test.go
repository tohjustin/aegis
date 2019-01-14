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
		requestPath:   "/static?subject=testSubject&status=testStatus&color=ff0000",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody: badge.CreateUnsafe("testSubject", "testStatus",
			&badge.Options{Color: "ff0000"}),
	})
}

func TestBadgeServiceHandlerWithCSSColorName(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/static?subject=testSubject&status=testStatus&color=red",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody: badge.CreateUnsafe("testSubject", "testStatus",
			&badge.Options{Color: "red"}),
	})
}

func TestBadgeServiceHandlerWithNoColor(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/static?subject=testSubject&status=testStatus",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody:   badge.CreateUnsafe("testSubject", "testStatus", nil),
	})
}

func TestBadgeServiceHandlerWithBadColor(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/static?subject=testSubject&status=testStatus&color=badColor",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody:   badge.CreateUnsafe("testSubject", "testStatus", nil),
	})
}

func TestBadgeServiceHandlerWithIconQuery(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/static?subject=testSubject&status=testStatus&icon=brands/docker",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody: badge.CreateUnsafe("testSubject", "testStatus",
			&badge.Options{Icon: "brands/docker"}),
	})
}

func TestBadgeServiceHandlerWithBadIconQuery(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/static?subject=testSubject&status=testStatus&icon=badIcon",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody:   badge.CreateUnsafe("testSubject", "testStatus", nil),
	})
}

func TestBadgeServiceHandlerWithStyleQuery(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/static?subject=testSubject&status=testStatus&style=semaphore",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody: badge.CreateUnsafe("testSubject", "testStatus",
			&badge.Options{Style: "semaphore"}),
	})
}

func TestBadgeServiceHandlerWithBadStyleQuery(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/static?subject=testSubject&status=testStatus&style=badStyle",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody:   badge.CreateUnsafe("testSubject", "testStatus", nil),
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
			requestPath:     "/static?subject=testSubject&status=testStatus&color=ff0000",
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
		expectedBody:   badge.CreateUnsafe("badger", "400 Bad Request", nil),
	})
}

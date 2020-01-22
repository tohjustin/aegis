package main

import (
	"net/http"
	"testing"

	"github.com/tohjustin/badger/pkg/badge"
)

func createBadge(params *badge.Params) string {
	generatedBadge, err := badge.Create(params)
	if err != nil {
		panic(err)
	}

	return generatedBadge
}

func TestStaticBadgeService(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/static?subject=testSubject&status=testStatus&color=ff0000",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody: createBadge(&badge.Params{
			Subject: "testSubject",
			Status:  "testStatus",
			Color:   "ff0000",
		}),
	})
}

func TestStaticBadgeServiceWithCSSColorName(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/static?subject=testSubject&status=testStatus&color=red",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody: createBadge(&badge.Params{
			Subject: "testSubject",
			Status:  "testStatus",
			Color:   "red",
		}),
	})
}

func TestStaticBadgeServiceWithNoColor(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/static?subject=testSubject&status=testStatus",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody: createBadge(&badge.Params{
			Subject: "testSubject",
			Status:  "testStatus",
			Color:   badge.DefaultColor,
		}),
	})
}

func TestStaticBadgeServiceWithBadColor(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/static?subject=testSubject&status=testStatus&color=badColor",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody: createBadge(&badge.Params{
			Subject: "testSubject",
			Status:  "testStatus",
			Color:   badge.DefaultColor,
		}),
	})
}

func TestStaticBadgeServiceWithIconQuery(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/static?subject=testSubject&status=testStatus&icon=brands/docker",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody: createBadge(&badge.Params{
			Subject: "testSubject",
			Status:  "testStatus",
			Icon:    "brands/docker",
		}),
	})
}

func TestStaticBadgeServiceWithBadIconQuery(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/static?subject=testSubject&status=testStatus&icon=badIcon",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody: createBadge(&badge.Params{
			Subject: "testSubject",
			Status:  "testStatus",
		}),
	})
}

func TestStaticBadgeServiceWithStyleQuery(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/static?subject=testSubject&status=testStatus&style=semaphoreci",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody: createBadge(&badge.Params{
			Style:   badge.SemaphoreCIStyle,
			Subject: "testSubject",
			Status:  "testStatus",
		}),
	})
}

func TestStaticBadgeServiceWithBadStyleQuery(t *testing.T) {
	t.Parallel()

	runHTTPTest(t, httpTestCase{
		requestMethod: "GET",
		requestPath:   "/static?subject=testSubject&status=testStatus&style=badStyle",
		expectedHeaders: map[string]string{
			"Cache-Control": "public, immutable, max-age=86400, s-maxage=31536000",
			"Content-Type":  "image/svg+xml;utf-8",
		},
		expectedStatus: 200,
		expectedBody: createBadge(&badge.Params{
			Style:   badge.DefaultStyle,
			Subject: "testSubject",
			Status:  "testStatus",
		}),
	})
}

func TestStaticBadgeServiceWithBadHTTPMethods(t *testing.T) {
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

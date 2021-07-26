package service

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tohjustin/aegis/service/config"
	"go.uber.org/zap/zaptest"
)

type httpTestCase struct {
	requestMethod   string
	requestPath     string
	expectedHeaders map[string]string
	expectedStatus  int
	expectedBody    string
}

func runHTTPTest(t *testing.T, testCase httpTestCase) {
	req, err := http.NewRequest(testCase.requestMethod, testCase.requestPath, nil)
	if err != nil {
		t.Fatal(err)
	}

	// TODO: Create proper mock dependencies & service generators
	mockLogger := zaptest.NewLogger(t)
	mockConfig := &config.Config{}
	mockStaticService, err := NewStaticService(mockConfig, mockLogger)
	if err != nil {
		t.Fatalf(err.Error())
	}
	mockGitProviderService, err := NewGitlabService(mockConfig, mockLogger)
	if err != nil {
		t.Fatalf(err.Error())
	}

	testServer := &Application{
		info:             Info{},
		logger:           mockLogger,
		config:           mockConfig,
		rootCmd:          nil,
		staticService:    &mockStaticService,
		bitbucketService: &mockGitProviderService,
		githubService:    &mockGitProviderService,
		gitlabService:    &mockGitProviderService,
	}
	res := httptest.NewRecorder()
	testServer.handler().ServeHTTP(res, req)

	// Check response header
	for fieldName, expectedFieldValue := range testCase.expectedHeaders {
		if fieldValue := res.Result().Header.Get(fieldName); fieldValue != expectedFieldValue {
			t.Errorf("handler returned wrong %v header: got %v want %v",
				fieldName, fieldValue, expectedFieldValue)
		}
	}

	// Check response status code
	if status := res.Code; status != testCase.expectedStatus {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, testCase.expectedStatus)
	}

	// Check response body
	if body := res.Body.String(); body != testCase.expectedBody {
		t.Errorf("handler returned unexpected body: got %v want %v",
			body, testCase.expectedBody)
	}
}

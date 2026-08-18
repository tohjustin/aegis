package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSONError(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	writeJSONError(rec, http.StatusBadRequest, "bad request")
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", rec.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["error"] != "bad request" {
		t.Errorf("error = %q", body["error"])
	}

	rec = httptest.NewRecorder()
	writeJSONError(rec, http.StatusInternalServerError, "secret")
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["error"] != "internal error" {
		t.Errorf("5xx leaked %q", body["error"])
	}
}

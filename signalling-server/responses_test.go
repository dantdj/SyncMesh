package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	recorder := httptest.NewRecorder()
	headers := http.Header{}
	headers.Set("X-Test", "value")

	payload := envelope{
		"status": "ok",
		"count":  2,
	}

	if err := writeJSON(recorder, http.StatusCreated, payload, headers); err != nil {
		t.Fatalf("writeJSON returned error: %v", err)
	}

	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", recorder.Code)
	}

	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("expected Content-Type application/json, got %q", got)
	}

	if got := recorder.Header().Get("X-Test"); got != "value" {
		t.Fatalf("expected X-Test header to be set, got %q", got)
	}

	var decoded map[string]any
	if err := json.NewDecoder(recorder.Body).Decode(&decoded); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if decoded["status"] != "ok" {
		t.Fatalf("expected status ok, got %v", decoded["status"])
	}
}

func TestErrorResponses(t *testing.T) {
	recorder := httptest.NewRecorder()
	errorResponse(recorder, http.StatusBadRequest, "bad request")

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", recorder.Code)
	}

	var decoded map[string]any
	if err := json.NewDecoder(recorder.Body).Decode(&decoded); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if decoded["error"] != "bad request" {
		t.Fatalf("expected error message, got %v", decoded["error"])
	}
}

func TestNotFoundResponse(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	recorder := httptest.NewRecorder()

	notFoundResponse(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", recorder.Code)
	}
}

func TestMethodNotAllowedResponse(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/ping", nil)
	recorder := httptest.NewRecorder()

	methodNotAllowedResponse(recorder, req)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", recorder.Code)
	}
}

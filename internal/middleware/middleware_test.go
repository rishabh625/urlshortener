package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Simple test handler that just responds with status 200
func testHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestLatencyMiddleware(t *testing.T) {
	// Create a buffer to capture log output
	var logBuffer bytes.Buffer
	log.SetOutput(&logBuffer) // Capture log output

	// Create a new handler wrapped with the LatencyMiddleware
	handler := LatencyMiddleware(http.HandlerFunc(testHandler))

	// Create a test request
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Call the handler
	handler.ServeHTTP(rr, req)

	// Check the status code to ensure the handler is working
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status 200, got %v", status)
	}

	// Check if the log output contains the expected latency message
	expectedLogOutput := "Method: GET, Path: /test, Latency:"
	if containsLatencyLog(logBuffer.String(), expectedLogOutput) {
		t.Errorf("expected log output to contain '%s', but got %s", expectedLogOutput, logBuffer.String())
	}
}

func containsLatencyLog(logOutput, expectedLog string) bool {
	return strings.Contains(logOutput, expectedLog)
}

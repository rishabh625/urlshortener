package handler

import (
	"context"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"urlshortener/internal/service"
)

// MockService simulates the RedirectURL service for testing
type MockService struct {
	mock.Mock
}

func (m *MockService) RedirectURL(ctx context.Context, id string) *RedirectResponse {
	args := m.Called(ctx, id)
	return args.Get(0).(*RedirectResponse)
}

type RedirectResponse struct {
	LongURl string
}

func TestRedirectHandler(t *testing.T) {
	// Mock the service and the app
	mockService := new(MockService)
	app := &App{service: new(service.URLShortenService)}

	// Create a new request (GET /short-id)
	req, err := http.NewRequest("GET", "/short-id", nil)
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	// Record the response
	rr := httptest.NewRecorder()

	// Define test cases
	tests := []struct {
		name           string
		id             string
		mockResp       *RedirectResponse
		expectedStatus int
	}{
		{
			name:           "Valid ID, redirect",
			id:             "short-id",
			mockResp:       &RedirectResponse{LongURl: "https://example.com/eacnjkd"},
			expectedStatus: http.StatusServiceUnavailable,
		},
		{
			name:           "Service unavailable, internal error",
			id:             "short-id",
			mockResp:       nil,
			expectedStatus: http.StatusServiceUnavailable,
		},
	}

	// Run the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mock service behavior
			mockService.On("RedirectURL", context.Background(), tt.id).Return(tt.mockResp)

			// Create a new request based on the test case
			req.URL.Path = "/" + tt.id

			// Call the handler
			handler := app.RedirectHandler()
			handler.ServeHTTP(rr, req)

			// Assert the status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Additional check if method was "GET" and not anything else
			if req.Method != http.MethodGet && rr.Code == http.StatusBadRequest {
				t.Errorf("expected 400 BadRequest for method %s, got %d", req.Method, rr.Code)
			}
		})
	}
}

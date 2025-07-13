package handler_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pierre-teodoresco/calculator-api-go/internal/handler"
)

func TestHealthHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Valid GET request",
			method:         "GET",
			expectedStatus: http.StatusOK,
			expectedBody:   "API is running fine",
		},
		{
			name:           "POST request (should still work)",
			method:         "POST",
			expectedStatus: http.StatusOK,
			expectedBody:   "API is running fine",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/health", nil)
			w := httptest.NewRecorder()

			handler.HealthHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d but got %d", tt.expectedStatus, w.Code)
			}

			body, err := io.ReadAll(w.Body)
			if err != nil {
				t.Errorf("Failed to read response body: %v", err)
			}

			if string(body) != tt.expectedBody {
				t.Errorf("Expected body '%s' but got '%s'", tt.expectedBody, string(body))
			}
		})
	}
}

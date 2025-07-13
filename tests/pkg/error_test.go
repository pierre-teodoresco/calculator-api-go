package pkg_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pierre-teodoresco/calculator-api-go/pkg"
)

func TestSetJSONHeader(t *testing.T) {
	w := httptest.NewRecorder()
	pkg.SetJSONHeader(w)

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json but got %s", contentType)
	}
}

func TestSendJSONError(t *testing.T) {
	tests := []struct {
		name           string
		error          *pkg.Error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Bad Request Error",
			error:          &pkg.Error{Message: "Invalid input", Status: http.StatusBadRequest},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Invalid input"}`,
		},
		{
			name:           "Internal Server Error",
			error:          &pkg.Error{Message: "Something went wrong", Status: http.StatusInternalServerError},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Something went wrong"}`,
		},
		{
			name:           "Unsupported Media Type",
			error:          &pkg.Error{Message: "Content-Type must be application/json", Status: http.StatusUnsupportedMediaType},
			expectedStatus: http.StatusUnsupportedMediaType,
			expectedBody:   `{"error":"Content-Type must be application/json"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			pkg.SendJSONError(w, tt.error)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d but got %d", tt.expectedStatus, w.Code)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json but got %s", contentType)
			}

			var response map[string]interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			if err != nil {
				t.Errorf("Failed to decode response: %v", err)
			}

			if response["error"] != tt.error.Message {
				t.Errorf("Expected error message '%s' but got '%s'", tt.error.Message, response["error"])
			}
		})
	}
}

func TestSendJSON(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		content        interface{}
		expectedStatus int
	}{
		{
			name:           "Success response with struct",
			statusCode:     http.StatusOK,
			content:        struct{ Result int }{Result: 42},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Created response with map",
			statusCode:     http.StatusCreated,
			content:        map[string]string{"message": "Created successfully"},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "Response with string",
			statusCode:     http.StatusAccepted,
			content:        "Processing",
			expectedStatus: http.StatusAccepted,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			pkg.SendJSON(w, tt.statusCode, tt.content)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d but got %d", tt.expectedStatus, w.Code)
			}

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json but got %s", contentType)
			}

			var response interface{}
			err := json.NewDecoder(w.Body).Decode(&response)
			if err != nil {
				t.Errorf("Failed to decode response: %v", err)
			}

			// Basic check that response is not nil
			if response == nil {
				t.Errorf("Expected non-nil response")
			}
		})
	}
}

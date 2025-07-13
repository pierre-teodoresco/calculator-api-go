package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pierre-teodoresco/calculator-api-go/internal/handler"
	"github.com/pierre-teodoresco/calculator-api-go/pkg"
)

func TestParseRequest(t *testing.T) {
	tests := []struct {
		name           string
		contentType    string
		body           string
		expectedArgs   handler.Args
		expectedErr    bool
		expectedStatus int
	}{
		{
			name:         "Valid JSON request",
			contentType:  "application/json",
			body:         `{"a": 5, "b": 3}`,
			expectedArgs: handler.Args{A: 5, B: 3},
			expectedErr:  false,
		},
		{
			name:         "Valid JSON without content-type",
			contentType:  "",
			body:         `{"a": 10, "b": 2}`,
			expectedArgs: handler.Args{A: 10, B: 2},
			expectedErr:  false,
		},
		{
			name:           "Invalid content-type",
			contentType:    "text/plain",
			body:           `{"a": 5, "b": 3}`,
			expectedErr:    true,
			expectedStatus: http.StatusUnsupportedMediaType,
		},
		{
			name:           "Malformed JSON",
			contentType:    "application/json",
			body:           `{"a": 5, "b":}`,
			expectedErr:    true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Invalid field type",
			contentType:    "application/json",
			body:           `{"a": "string", "b": 3}`,
			expectedErr:    true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Unknown field",
			contentType:    "application/json",
			body:           `{"a": 5, "b": 3, "c": 1}`,
			expectedErr:    true,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:         "Missing fields",
			contentType:  "application/json",
			body:         `{"a": 5}`,
			expectedArgs: handler.Args{A: 5, B: 0},
			expectedErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", strings.NewReader(tt.body))
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}

			args, err := handler.ParseRequest(req)

			if tt.expectedErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if err.Status != tt.expectedStatus {
					t.Errorf("Expected status %d but got %d", tt.expectedStatus, err.Status)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if args != tt.expectedArgs {
					t.Errorf("Expected args %v but got %v", tt.expectedArgs, args)
				}
			}
		})
	}
}

func TestAddHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedResult int
		expectError    bool
	}{
		{
			name:           "Valid addition",
			requestBody:    `{"a": 5, "b": 3}`,
			expectedStatus: http.StatusOK,
			expectedResult: 8,
			expectError:    false,
		},
		{
			name:           "Addition with negative numbers",
			requestBody:    `{"a": -5, "b": 3}`,
			expectedStatus: http.StatusOK,
			expectedResult: -2,
			expectError:    false,
		},
		{
			name:           "Addition with zero",
			requestBody:    `{"a": 0, "b": 10}`,
			expectedStatus: http.StatusOK,
			expectedResult: 10,
			expectError:    false,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"a": 5, "b":}`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/add", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.AddHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d but got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectError {
				var result handler.Result
				err := json.NewDecoder(w.Body).Decode(&result)
				if err != nil {
					t.Errorf("Failed to decode response: %v", err)
				}
				if result.Value != tt.expectedResult {
					t.Errorf("Expected result %d but got %d", tt.expectedResult, result.Value)
				}
			} else {
				var errorResponse pkg.Error
				err := json.NewDecoder(w.Body).Decode(&errorResponse)
				if err != nil {
					t.Errorf("Failed to decode error response: %v", err)
				}
				if errorResponse.Message == "" {
					t.Errorf("Expected error message but got empty string")
				}
			}
		})
	}
}

func TestMultiplyHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedResult int
		expectError    bool
	}{
		{
			name:           "Valid multiplication",
			requestBody:    `{"a": 5, "b": 3}`,
			expectedStatus: http.StatusOK,
			expectedResult: 15,
			expectError:    false,
		},
		{
			name:           "Multiplication with zero",
			requestBody:    `{"a": 5, "b": 0}`,
			expectedStatus: http.StatusOK,
			expectedResult: 0,
			expectError:    false,
		},
		{
			name:           "Multiplication with negative numbers",
			requestBody:    `{"a": -5, "b": 3}`,
			expectedStatus: http.StatusOK,
			expectedResult: -15,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/multiply", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.MultiplyHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d but got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectError {
				var result handler.Result
				err := json.NewDecoder(w.Body).Decode(&result)
				if err != nil {
					t.Errorf("Failed to decode response: %v", err)
				}
				if result.Value != tt.expectedResult {
					t.Errorf("Expected result %d but got %d", tt.expectedResult, result.Value)
				}
			}
		})
	}
}

func TestSubtractHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedResult int
		expectError    bool
	}{
		{
			name:           "Valid subtraction",
			requestBody:    `{"a": 5, "b": 3}`,
			expectedStatus: http.StatusOK,
			expectedResult: 2,
			expectError:    false,
		},
		{
			name:           "Subtraction resulting in negative",
			requestBody:    `{"a": 3, "b": 5}`,
			expectedStatus: http.StatusOK,
			expectedResult: -2,
			expectError:    false,
		},
		{
			name:           "Subtraction with zero",
			requestBody:    `{"a": 5, "b": 0}`,
			expectedStatus: http.StatusOK,
			expectedResult: 5,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/subtract", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.SubtractHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d but got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectError {
				var result handler.Result
				err := json.NewDecoder(w.Body).Decode(&result)
				if err != nil {
					t.Errorf("Failed to decode response: %v", err)
				}
				if result.Value != tt.expectedResult {
					t.Errorf("Expected result %d but got %d", tt.expectedResult, result.Value)
				}
			}
		})
	}
}

func TestDivideHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedResult int
		expectError    bool
	}{
		{
			name:           "Valid division",
			requestBody:    `{"a": 10, "b": 2}`,
			expectedStatus: http.StatusOK,
			expectedResult: 5,
			expectError:    false,
		},
		{
			name:           "Division by zero",
			requestBody:    `{"a": 10, "b": 0}`,
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:           "Division with negative numbers",
			requestBody:    `{"a": -10, "b": 2}`,
			expectedStatus: http.StatusOK,
			expectedResult: -5,
			expectError:    false,
		},
		{
			name:           "Division with remainder (integer division)",
			requestBody:    `{"a": 10, "b": 3}`,
			expectedStatus: http.StatusOK,
			expectedResult: 3,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/divide", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.DivideHandler(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d but got %d", tt.expectedStatus, w.Code)
			}

			if !tt.expectError {
				var result handler.Result
				err := json.NewDecoder(w.Body).Decode(&result)
				if err != nil {
					t.Errorf("Failed to decode response: %v", err)
				}
				if result.Value != tt.expectedResult {
					t.Errorf("Expected result %d but got %d", tt.expectedResult, result.Value)
				}
			} else {
				var errorResponse pkg.Error
				err := json.NewDecoder(w.Body).Decode(&errorResponse)
				if err != nil {
					t.Errorf("Failed to decode error response: %v", err)
				}
				if errorResponse.Message == "" {
					t.Errorf("Expected error message but got empty string")
				}
			}
		})
	}
}

func TestHandlersContentType(t *testing.T) {
	handlers := []struct {
		name    string
		handler func(http.ResponseWriter, *http.Request)
	}{
		{"AddHandler", handler.AddHandler},
		{"MultiplyHandler", handler.MultiplyHandler},
		{"SubtractHandler", handler.SubtractHandler},
		{"DivideHandler", handler.DivideHandler},
	}

	for _, h := range handlers {
		t.Run(h.name+"_ContentType", func(t *testing.T) {
			req := httptest.NewRequest("POST", "/test", strings.NewReader(`{"a": 5, "b": 3}`))
			w := httptest.NewRecorder()

			h.handler(w, req)

			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json but got %s", contentType)
			}
		})
	}
}

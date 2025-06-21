package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestErrorResponse_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      ErrorResponse
		expected string
	}{
		{
			name: "Error with internal error",
			err: ErrorResponse{
				Message: "Task not found",
				Code:    http.StatusNotFound,
				Err:     errors.New("internal error details"),
			},
			expected: "internal error details",
		},
		{
			name: "Error without internal error",
			err: ErrorResponse{
				Message: "Invalid request",
				Code:    http.StatusBadRequest,
			},
			expected: "Invalid request",
		},
		{
			name: "Error with empty message and no internal error",
			err: ErrorResponse{
				Code: http.StatusInternalServerError,
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestWriteErrorResponse(t *testing.T) {
	tests := []struct {
		name           string
		err            ErrorResponse
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Standard error response",
			err: ErrorResponse{
				Message: "Task not found",
				Code:    http.StatusNotFound,
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"Not Found","message":"Task not found","code":404}`,
		},
		{
			name: "Bad request error",
			err: ErrorResponse{
				Message: "Invalid input",
				Code:    http.StatusBadRequest,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"Bad Request","message":"Invalid input","code":400}`,
		},
		{
			name: "Internal server error",
			err: ErrorResponse{
				Message: "Something went wrong",
				Code:    http.StatusInternalServerError,
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"Internal Server Error","message":"Something went wrong","code":500}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			writeErrorResponse(w, tt.err)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check content type
			expectedContentType := "application/json"
			if contentType := w.Header().Get("Content-Type"); contentType != expectedContentType {
				t.Errorf("Expected Content-Type %s, got %s", expectedContentType, contentType)
			}

			// Check response body (trim whitespace for comparison)
			body := strings.TrimSpace(w.Body.String())
			if body != tt.expectedBody {
				t.Errorf("Expected body %s, got %s", tt.expectedBody, body)
			}
		})
	}
}

func TestWriteJSONResponse(t *testing.T) {
	tests := []struct {
		name           string
		data           interface{}
		statusCode     int
		expectedStatus int
		expectedBody   string
		expectError    bool
	}{
		{
			name:           "Simple string response",
			data:           "hello",
			statusCode:     http.StatusOK,
			expectedStatus: http.StatusOK,
			expectedBody:   `"hello"`,
			expectError:    false,
		},
		{
			name: "JSON object response",
			data: map[string]string{
				"message": "success",
			},
			statusCode:     http.StatusCreated,
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"message":"success"}`,
			expectError:    false,
		},
		{
			name:           "Empty response",
			data:           nil,
			statusCode:     http.StatusNoContent,
			expectedStatus: http.StatusNoContent,
			expectedBody:   `null`,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := writeJSONResponse(w, tt.data, tt.statusCode)

			// Check error expectation
			if tt.expectError && err == nil {
				t.Error("Expected an error, but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected no error, but got: %v", err)
			}

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check content type
			expectedContentType := "application/json"
			if contentType := w.Header().Get("Content-Type"); contentType != expectedContentType {
				t.Errorf("Expected Content-Type %s, got %s", expectedContentType, contentType)
			}

			// Check response body (trim whitespace for comparison)
			body := strings.TrimSpace(w.Body.String())
			if body != tt.expectedBody {
				t.Errorf("Expected body %s, got %s", tt.expectedBody, body)
			}
		})
	}
}

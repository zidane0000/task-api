package handlers

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse represents a structured API error response
type ErrorResponse struct {
	Status  string `json:"error"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
	Err     error  `json:"-"` // Internal error (not exposed to client)
}

func (e ErrorResponse) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

var (
	ErrInvalidJSON      = ErrorResponse{Message: "Invalid JSON in request body", Code: http.StatusBadRequest}
	ErrTaskNotFound     = ErrorResponse{Message: "Task not found", Code: http.StatusNotFound}
	ErrInvalidTaskID    = ErrorResponse{Message: "Invalid task ID", Code: http.StatusBadRequest}
	ErrMethodNotAllowed = ErrorResponse{Message: "Method not allowed", Code: http.StatusMethodNotAllowed}
	ErrInternalServer   = ErrorResponse{Message: "Internal server error", Code: http.StatusInternalServerError}
)

// writeErrorResponse writes a structured error response to the client
func writeErrorResponse(w http.ResponseWriter, err ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)

	response := ErrorResponse{
		Status:  http.StatusText(err.Code),
		Message: err.Message,
		Code:    err.Code,
	}

	// If JSON encoding fails, fall back to http.Error
	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		http.Error(w, err.Message, err.Code)
	}
}

// writeJSONResponse writes a successful JSON response
func writeJSONResponse(w http.ResponseWriter, data interface{}, statusCode int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

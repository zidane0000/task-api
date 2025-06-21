package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"task-api/internal/models"
	"task-api/internal/storage"
	"testing"

	"github.com/go-chi/chi/v5"
)

// Mock implementation right in the test file
type mockTaskStorage struct{}

func (m *mockTaskStorage) Create(task *models.Task) (*models.Task, error) {
	return nil, errors.New("storage create failed")
}

func (m *mockTaskStorage) GetAll() ([]*models.Task, error) {
	return nil, errors.New("storage getall failed")
}

func (m *mockTaskStorage) GetByID(id int) (*models.Task, error) {
	return nil, errors.New("storage getbyid failed")
}

func (m *mockTaskStorage) Update(task *models.Task) error {
	return errors.New("storage update failed")
}

func (m *mockTaskStorage) Delete(id int) error {
	return errors.New("storage delete failed")
}

// setupTestHandler creates a handler with in-memory storage for testing
func setupTestHandler() *TaskHandler {
	testStorage := storage.NewInMemoryStorage()
	return NewTaskHandler(testStorage)
}

// setupTestHandlerWithMock creates a handler with mock storage for testing
func setupTestHandlerWithMock() *TaskHandler {
	testStorage := &mockTaskStorage{}
	return NewTaskHandler(testStorage)
}

// TestTaskHandler_GetAllTasks_EmptyStorage tests retrieving tasks from empty storage
func TestTaskHandler_GetAllTasks_EmptyStorage(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	handler.GetAllTasks(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var tasks []*models.Task
	if err := json.NewDecoder(w.Body).Decode(&tasks); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tasks))
	}
}

// TestTaskHandler_GetAllTasks tests retrieving all tasks
func TestTaskHandler_GetAllTasks(t *testing.T) {
	handler := setupTestHandler()

	// Add some test tasks
	task1, _ := models.NewTask("Task 1", 0)
	_, err := handler.storage.Create(task1)
	if err != nil {
		t.Fatalf("Failed to create task1: %v", err)
	}
	task2, _ := models.NewTask("Task 2", 1)
	_, err = handler.storage.Create(task2)
	if err != nil {
		t.Fatalf("Failed to create task2: %v", err)
	}

	// Test GET /tasks
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	handler.GetAllTasks(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var tasks []*models.Task
	if err := json.NewDecoder(w.Body).Decode(&tasks); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
}

// TestTaskHandler_CreateTask tests creating a new task
func TestTaskHandler_CreateTask(t *testing.T) {
	handler := setupTestHandler()

	tests := []struct {
		name           string
		requestBody    interface{}
		expectedStatus int
	}{
		{
			name:           "valid task",
			requestBody:    map[string]interface{}{"name": "New Task", "status": 0},
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "empty name",
			requestBody:    map[string]interface{}{"name": "", "status": 0},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid status",
			requestBody:    map[string]interface{}{"name": "Test Task", "status": 2},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handler.CreateTask(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusCreated {
				var createdTask models.Task
				if err := json.NewDecoder(w.Body).Decode(&createdTask); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if createdTask.ID == 0 {
					t.Error("Expected task to have an ID")
				}
			}
		})
	}
}

// TestTaskHandler_CreateTask_InvalidJSON tests creating task with invalid JSON
func TestTaskHandler_CreateTask_InvalidJSON(t *testing.T) {
	handler := setupTestHandler()

	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBufferString("{invalid json"))
	w := httptest.NewRecorder()

	handler.CreateTask(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}
}

// TestTaskHandler_UpdateTask tests updating an existing task
func TestTaskHandler_UpdateTask(t *testing.T) {
	handler := setupTestHandler()

	// Create a task first
	task, _ := models.NewTask("Original Task", 0)
	createdTask, _ := handler.storage.Create(task)

	tests := []struct {
		name           string
		taskID         string
		requestBody    interface{}
		expectedStatus int
	}{
		{
			name:           "valid update",
			taskID:         "1",
			requestBody:    map[string]interface{}{"name": "Updated Task", "status": 1},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "non-existent task",
			taskID:         "999",
			requestBody:    map[string]interface{}{"name": "Updated Task", "status": 1},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid task data",
			taskID:         "1",
			requestBody:    map[string]interface{}{"name": "", "status": 0},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid task ID",
			taskID:         "abc",
			requestBody:    map[string]interface{}{"name": "Updated Task", "status": 1},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/tasks/"+tt.taskID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Add chi URL parameter to route context
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.taskID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			handler.UpdateTask(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedStatus == http.StatusOK {
				var updatedTask models.Task
				if err := json.NewDecoder(w.Body).Decode(&updatedTask); err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}

				if updatedTask.ID != createdTask.ID {
					t.Error("Task ID should remain the same after update")
				}
			}
		})
	}
}

// TestTaskHandler_DeleteTask tests deleting a task
func TestTaskHandler_DeleteTask(t *testing.T) {
	handler := setupTestHandler()

	// Create a task first
	task, _ := models.NewTask("Task to Delete", 0)
	createdTask, _ := handler.storage.Create(task)

	tests := []struct {
		name           string
		taskID         string
		expectedStatus int
	}{
		{
			name:           "valid deletion",
			taskID:         "1",
			expectedStatus: http.StatusNoContent,
		},
		{
			name:           "non-existent task",
			taskID:         "999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "invalid task ID",
			taskID:         "abc",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/tasks/"+tt.taskID, nil)

			// Add chi URL parameter to route context
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.taskID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			handler.DeleteTask(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Verify task is deleted for successful deletion
			if tt.expectedStatus == http.StatusNoContent {
				_, err := handler.storage.GetByID(createdTask.ID)
				if err == nil {
					t.Error("Task should be deleted but still exists")
				}
			}
		})
	}
}

// TestTaskHandler_GetAllTasks_StorageError tests GetAllTasks when storage fails
func TestTaskHandler_GetAllTasks_StorageError(t *testing.T) {
	handler := setupTestHandlerWithMock()

	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	w := httptest.NewRecorder()

	handler.GetAllTasks(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var errorResponse map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&errorResponse); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	if errorResponse["message"] != "Internal server error" {
		t.Errorf("Expected 'Internal server error', got %v", errorResponse["message"])
	}
}

// TestTaskHandler_CreateTask_StorageError tests CreateTask when storage fails
func TestTaskHandler_CreateTask_StorageError(t *testing.T) {
	handler := setupTestHandlerWithMock()

	requestBody := map[string]interface{}{"name": "Valid Task", "status": 0}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateTask(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	var errorResponse map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&errorResponse); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	if errorResponse["message"] != "Internal server error" {
		t.Errorf("Expected 'Internal server error', got %v", errorResponse["message"])
	}
}

// TestTaskHandler_UpdateTask_StorageErrors tests UpdateTask with various storage failures
func TestTaskHandler_UpdateTask_StorageErrors(t *testing.T) {
	handler := setupTestHandlerWithMock()

	tests := []struct {
		name           string
		taskID         string
		requestBody    interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "storage GetByID fails",
			taskID:         "1",
			requestBody:    map[string]interface{}{"name": "Updated Task", "status": 1},
			expectedStatus: http.StatusNotFound,
			expectedError:  "Task not found",
		},
		{
			name:           "storage Update fails after successful GetByID",
			taskID:         "1",
			requestBody:    map[string]interface{}{"name": "Updated Task", "status": 1},
			expectedStatus: http.StatusNotFound, // GetByID fails first
			expectedError:  "Task not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPut, "/tasks/"+tt.taskID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Add chi URL parameter to route context
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.taskID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			handler.UpdateTask(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var errorResponse map[string]interface{}
			if err := json.NewDecoder(w.Body).Decode(&errorResponse); err != nil {
				t.Fatalf("Failed to decode error response: %v", err)
			}

			if errorResponse["message"] != tt.expectedError {
				t.Errorf("Expected error '%s', got %v", tt.expectedError, errorResponse["message"])
			}
		})
	}
}

// TestTaskHandler_DeleteTask_StorageErrors tests DeleteTask with various storage failures
func TestTaskHandler_DeleteTask_StorageErrors(t *testing.T) {
	handler := setupTestHandlerWithMock()

	tests := []struct {
		name           string
		taskID         string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "storage GetByID fails",
			taskID:         "1",
			expectedStatus: http.StatusNotFound,
			expectedError:  "Task not found",
		},
		{
			name:           "storage Delete fails after successful GetByID",
			taskID:         "1",
			expectedStatus: http.StatusNotFound, // GetByID fails first
			expectedError:  "Task not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/tasks/"+tt.taskID, nil)

			// Add chi URL parameter to route context
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.taskID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			w := httptest.NewRecorder()

			handler.DeleteTask(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var errorResponse map[string]interface{}
			if err := json.NewDecoder(w.Body).Decode(&errorResponse); err != nil {
				t.Fatalf("Failed to decode error response: %v", err)
			}

			if errorResponse["message"] != tt.expectedError {
				t.Errorf("Expected error '%s', got %v", tt.expectedError, errorResponse["message"])
			}
		})
	}
}

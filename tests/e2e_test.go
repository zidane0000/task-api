package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"task-api/internal/models"
	"testing"
	"time"
)

const (
	baseURL = "http://localhost:8080"
	timeout = 30 * time.Second
)

// TestTaskAPI_E2E tests the complete task API workflow against a real running server
func TestTaskAPI_E2E(t *testing.T) {
	// Wait for server to be ready
	if !waitForServer(t) {
		t.Fatal("Server is not ready")
	}

	// Test 1: Health check
	t.Run("GET /health", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/health")
		if err != nil {
			t.Fatalf("Failed to GET /health: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}
	})

	// Test 2: Get initial tasks count
	var initialTaskCount int
	t.Run("GET /tasks - initial state", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/tasks")
		if err != nil {
			t.Fatalf("Failed to GET /tasks: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var tasks []*models.Task
		if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		initialTaskCount = len(tasks)
		t.Logf("Initial task count: %d", initialTaskCount)
	})

	// Test 3: Create first task
	var task1 *models.Task
	t.Run("POST /tasks - create first task", func(t *testing.T) {
		taskData := map[string]interface{}{
			"name":   "Complete project documentation",
			"status": 0,
		}
		body, _ := json.Marshal(taskData)

		resp, err := http.Post(baseURL+"/tasks", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to POST /tasks: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}

		if err := json.NewDecoder(resp.Body).Decode(&task1); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if task1.ID == 0 {
			t.Error("Expected task to have an ID")
		}
		if task1.Name != "Complete project documentation" {
			t.Errorf("Expected name 'Complete project documentation', got '%s'", task1.Name)
		}
		if task1.Status != 0 {
			t.Errorf("Expected status 0, got %d", task1.Status)
		}
	})

	// Test 4: Create second task
	var task2 *models.Task
	t.Run("POST /tasks - create second task", func(t *testing.T) {
		taskData := map[string]interface{}{
			"name":   "Review code changes",
			"status": 1,
		}
		body, _ := json.Marshal(taskData)

		resp, err := http.Post(baseURL+"/tasks", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to POST /tasks: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", resp.StatusCode)
		}

		if err := json.NewDecoder(resp.Body).Decode(&task2); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if task2.ID == 0 {
			t.Error("Expected task to have an ID")
		}
		if task2.ID == task1.ID {
			t.Error("Expected tasks to have unique IDs")
		}
	})

	// Test 5: Get all tasks (should have initial count + 2 tasks now)
	t.Run("GET /tasks - with new tasks", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/tasks")
		if err != nil {
			t.Fatalf("Failed to GET /tasks: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var tasks []*models.Task
		if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		expectedCount := initialTaskCount + 2
		if len(tasks) != expectedCount {
			t.Errorf("Expected %d tasks, got %d", expectedCount, len(tasks))
		}
	})

	// Test 6: Update first task
	t.Run("PUT /tasks/{id} - update task", func(t *testing.T) {
		updateData := map[string]interface{}{
			"name":   "Complete project documentation - UPDATED",
			"status": 1,
		}
		body, _ := json.Marshal(updateData)

		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/tasks/%d", baseURL, task1.ID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to PUT /tasks/%d: %v", task1.ID, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var updatedTask models.Task
		if err := json.NewDecoder(resp.Body).Decode(&updatedTask); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		if updatedTask.Name != "Complete project documentation - UPDATED" {
			t.Errorf("Expected updated name, got '%s'", updatedTask.Name)
		}
		if updatedTask.Status != 1 {
			t.Errorf("Expected status 1, got %d", updatedTask.Status)
		}
	})

	// Test 7: Delete second task
	t.Run("DELETE /tasks/{id} - delete task", func(t *testing.T) {
		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/tasks/%d", baseURL, task2.ID), nil)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to DELETE /tasks/%d: %v", task2.ID, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent {
			t.Errorf("Expected status 204, got %d", resp.StatusCode)
		}
	})

	// Test 8: Verify task was deleted
	t.Run("GET /tasks - after deletion", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/tasks")
		if err != nil {
			t.Fatalf("Failed to GET /tasks: %v", err)
		}
		defer resp.Body.Close()

		var tasks []*models.Task
		if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		expectedCount := initialTaskCount + 1 // initial + 1 remaining task
		if len(tasks) != expectedCount {
			t.Errorf("Expected %d tasks after deletion, got %d", expectedCount, len(tasks))
		}
	})
}

// TestTaskAPI_ErrorCases_E2E tests various error scenarios against a real server
func TestTaskAPI_ErrorCases_E2E(t *testing.T) {
	// Wait for server to be ready
	if !waitForServer(t) {
		t.Fatal("Server is not ready")
	}

	// Test error cases
	t.Run("POST /tasks - invalid JSON", func(t *testing.T) {
		resp, err := http.Post(baseURL+"/tasks", "application/json", bytes.NewBufferString("{invalid json"))
		if err != nil {
			t.Fatalf("Failed to POST /tasks: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("POST /tasks - invalid task data", func(t *testing.T) {
		taskData := map[string]interface{}{
			"name":   "",
			"status": 0,
		}
		body, _ := json.Marshal(taskData)

		resp, err := http.Post(baseURL+"/tasks", "application/json", bytes.NewBuffer(body))
		if err != nil {
			t.Fatalf("Failed to POST /tasks: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", resp.StatusCode)
		}
	})

	t.Run("PUT /tasks/999 - non-existent task", func(t *testing.T) {
		updateData := map[string]interface{}{
			"name":   "Updated Task",
			"status": 1,
		}
		body, _ := json.Marshal(updateData)

		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodPut, baseURL+"/tasks/999", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to PUT /tasks/999: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})

	t.Run("DELETE /tasks/999 - non-existent task", func(t *testing.T) {
		client := &http.Client{}
		req, _ := http.NewRequest(http.MethodDelete, baseURL+"/tasks/999", nil)

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to DELETE /tasks/999: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", resp.StatusCode)
		}
	})
}

// waitForServer waits for the server to be ready by checking the health endpoint
func waitForServer(t *testing.T) bool {
	client := &http.Client{Timeout: 5 * time.Second}
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		resp, err := client.Get(baseURL + "/health")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return true
			}
		}
		time.Sleep(1 * time.Second)
	}

	t.Logf("Server at %s did not become ready within %v", baseURL, timeout)
	return false
}

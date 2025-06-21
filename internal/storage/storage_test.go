package storage

import (
	"os"
	"task-api/internal/models"
	"testing"
)

// TestTaskStorage_Create tests task creation functionality
func TestTaskStorage_Create(t *testing.T) {
	storage := NewTaskStorage()

	tests := []struct {
		name     string
		taskName string
		status   int
		wantErr  bool
	}{
		{
			name:     "create valid incomplete task",
			taskName: "Test task",
			status:   0,
			wantErr:  false,
		},
		{
			name:     "create valid completed task",
			taskName: "Completed task",
			status:   1,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create task using models.NewTask for validation
			task, err := models.NewTask(tt.taskName, tt.status)
			if err != nil {
				t.Fatalf("Failed to create task model: %v", err)
			}

			// Store the task
			result, err := storage.Create(task)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify task was assigned an ID
			if result.ID == 0 {
				t.Error("Expected task to be assigned an ID, got 0")
			}

			// Verify task content is preserved
			if result.Name != tt.taskName {
				t.Errorf("Expected name %q, got %q", tt.taskName, result.Name)
			}

			if result.Status != tt.status {
				t.Errorf("Expected status %d, got %d", tt.status, result.Status)
			}
		})
	}
}

// TestTaskStorage_Create_UniqueIDs tests that created tasks get unique IDs
func TestTaskStorage_Create_UniqueIDs(t *testing.T) {
	storage := NewTaskStorage()

	task1, _ := models.NewTask("Task 1", 0)
	task2, _ := models.NewTask("Task 2", 1)

	result1, err := storage.Create(task1)
	if err != nil {
		t.Fatalf("Failed to create first task: %v", err)
	}

	result2, err := storage.Create(task2)
	if err != nil {
		t.Fatalf("Failed to create second task: %v", err)
	}

	if result1.ID == result2.ID {
		t.Errorf("Expected unique IDs, but both tasks got ID %d", result1.ID)
	}

	if result1.ID == 0 || result2.ID == 0 {
		t.Error("Expected non-zero IDs for both tasks")
	}
}

// TestTaskStorage_GetAll tests retrieving all tasks
func TestTaskStorage_GetAll(t *testing.T) {
	storage := NewTaskStorage()

	// Test empty storage
	tasks, err := storage.GetAll()
	if err != nil {
		t.Errorf("Unexpected error getting all tasks from empty storage: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks in empty storage, got %d", len(tasks))
	}
	// Add some tasks
	task1, _ := models.NewTask("Task 1", 0)
	task2, _ := models.NewTask("Task 2", 1)

	_, err = storage.Create(task1)
	if err != nil {
		t.Fatalf("Failed to create task1: %v", err)
	}
	_, err = storage.Create(task2)
	if err != nil {
		t.Fatalf("Failed to create task2: %v", err)
	}

	// Test with tasks
	tasks, err = storage.GetAll()
	if err != nil {
		t.Errorf("Unexpected error getting all tasks: %v", err)
	}
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
}

// TestTaskStorage_GetByID tests retrieving specific tasks by ID
func TestTaskStorage_GetByID(t *testing.T) {
	storage := NewTaskStorage()

	// Test non-existent task
	_, err := storage.GetByID(999)
	if err == nil {
		t.Error("Expected error when getting non-existent task")
	}

	// Create and store a task
	task, _ := models.NewTask("Test task", 0)
	created, err := storage.Create(task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Test retrieving existing task
	retrieved, err := storage.GetByID(created.ID)
	if err != nil {
		t.Errorf("Unexpected error retrieving task: %v", err)
	}

	if retrieved.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, retrieved.ID)
	}
	if retrieved.Name != created.Name {
		t.Errorf("Expected name %q, got %q", created.Name, retrieved.Name)
	}
	if retrieved.Status != created.Status {
		t.Errorf("Expected status %d, got %d", created.Status, retrieved.Status)
	}
}

// TestTaskStorage_Update tests updating existing tasks
func TestTaskStorage_Update(t *testing.T) {
	storage := NewTaskStorage()

	// Test updating non-existent task
	task, _ := models.NewTask("Non-existent", 0)
	task.ID = 999
	err := storage.Update(task)
	if err == nil {
		t.Error("Expected error when updating non-existent task")
	}

	// Create a task
	original, _ := models.NewTask("Original task", 0)
	created, err := storage.Create(original)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Update the task
	updated, _ := models.NewTask("Updated task", 1)
	updated.ID = created.ID
	err = storage.Update(updated)
	if err != nil {
		t.Errorf("Unexpected error updating task: %v", err)
	}

	// Verify update
	retrieved, err := storage.GetByID(created.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated task: %v", err)
	}

	if retrieved.Name != "Updated task" {
		t.Errorf("Expected updated name %q, got %q", "Updated task", retrieved.Name)
	}
	if retrieved.Status != 1 {
		t.Errorf("Expected updated status 1, got %d", retrieved.Status)
	}
}

// TestTaskStorage_Delete tests deleting tasks
func TestTaskStorage_Delete(t *testing.T) {
	storage := NewTaskStorage()

	// Test deleting non-existent task
	err := storage.Delete(999)
	if err == nil {
		t.Error("Expected error when deleting non-existent task")
	}

	// Create a task
	task, _ := models.NewTask("Task to delete", 0)
	created, err := storage.Create(task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Verify task exists
	_, err = storage.GetByID(created.ID)
	if err != nil {
		t.Fatalf("Task should exist before deletion: %v", err)
	}

	// Delete the task
	err = storage.Delete(created.ID)
	if err != nil {
		t.Errorf("Unexpected error deleting task: %v", err)
	}

	// Verify task is gone
	_, err = storage.GetByID(created.ID)
	if err == nil {
		t.Error("Expected error when getting deleted task")
	}
}

// TestNewTaskStorage_DatabaseBackend tests the database backend detection
func TestNewTaskStorage_DatabaseBackend(t *testing.T) {
	// Set DATABASE_URL to trigger database backend detection
	os.Setenv("DATABASE_URL", "postgres://localhost/testdb")

	// This should panic because database backend is not yet implemented
	defer func() {
		if r := recover(); r != nil {
			// Expected panic - database backend not implemented
			expectedMessage := "Database backend not implemented yet. Remove DATABASE_URL to use memory backend."
			if panicMsg, ok := r.(string); ok {
				if panicMsg != expectedMessage {
					t.Errorf("Expected panic message %q, got %q", expectedMessage, panicMsg)
				}
			} else {
				t.Errorf("Expected string panic message, got %v", r)
			}
		} else {
			t.Error("Expected panic when DATABASE_URL is set, but no panic occurred")
		}
	}()

	// This should panic
	NewTaskStorage()
}

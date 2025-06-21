package models

import (
	"testing"
)

// TestTask_CreateTask tests task creation using table-driven approach
func TestTask_CreateTask(t *testing.T) {
	tests := []struct {
		name      string
		taskName  string
		status    int
		wantValid bool
	}{
		{
			name:      "valid incomplete task",
			taskName:  "Complete project documentation",
			status:    0,
			wantValid: true,
		},
		{
			name:      "valid completed task",
			taskName:  "Review code changes",
			status:    1,
			wantValid: true,
		},
		{
			name:      "empty name",
			taskName:  "",
			status:    0,
			wantValid: false,
		},
		{
			name:      "whitespace only name",
			taskName:  "   ",
			status:    0,
			wantValid: false,
		},
		{
			name:      "invalid status too high",
			taskName:  "Valid task name",
			status:    2,
			wantValid: false,
		},
		{
			name:      "invalid status negative",
			taskName:  "Valid task name",
			status:    -1,
			wantValid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := NewTask(tt.taskName, tt.status)

			if tt.wantValid {
				// Valid cases should not return error
				if err != nil {
					t.Errorf("Expected valid task, got error: %v", err)
					return
				}

				// Verify name and status are set correctly
				if task.Name != tt.taskName {
					t.Errorf("Expected Name to be %q, got %q", tt.taskName, task.Name)
				}

				if task.Status != tt.status {
					t.Errorf("Expected Status to be %d, got %d", tt.status, task.Status)
				}
			} else {
				// Invalid cases should return error
				if err == nil {
					t.Errorf("Expected error for invalid task, got nil")
				}

				if task != nil {
					t.Errorf("Expected nil task for invalid input, got %+v", task)
				}
			}
		})
	}
}

// TestTask_Update tests the Update method of Task using table-driven tests
func TestTask_Update(t *testing.T) {
	tests := []struct {
		name          string
		initialName   string
		initialStatus int
		updateName    string
		updateStatus  int
		wantErr       bool
		wantName      string
		wantStatus    int
	}{
		{
			name:          "valid update to completed",
			initialName:   "Initial Task",
			initialStatus: 0,
			updateName:    "Updated Task",
			updateStatus:  1,
			wantErr:       false,
			wantName:      "Updated Task",
			wantStatus:    1,
		},
		{
			name:          "valid update to incomplete",
			initialName:   "Initial Task",
			initialStatus: 1,
			updateName:    "Another Name",
			updateStatus:  0,
			wantErr:       false,
			wantName:      "Another Name",
			wantStatus:    0,
		},
		{
			name:         "empty name update",
			initialName:  "",
			updateName:   "Initial Task",
			updateStatus: 1,
			wantErr:      false,
			wantName:     "Initial Task",
			wantStatus:   1,
		},
		{
			name:          "whitespace name update",
			initialName:   "Initial Task",
			initialStatus: 0,
			updateName:    "   ",
			updateStatus:  1,
			wantErr:       true,
			wantName:      "Initial Task",
			wantStatus:    0,
		},
		{
			name:          "invalid status negative",
			initialName:   "Initial Task",
			initialStatus: 0,
			updateName:    "Valid Name",
			updateStatus:  -1,
			wantErr:       true,
			wantName:      "Initial Task",
			wantStatus:    0,
		},
		{
			name:          "invalid status too high",
			initialName:   "Initial Task",
			initialStatus: 0,
			updateName:    "Valid Name",
			updateStatus:  200,
			wantErr:       true,
			wantName:      "Initial Task",
			wantStatus:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := &Task{
				Name:   tt.initialName,
				Status: tt.initialStatus,
			}
			err := task.Update(tt.updateName, tt.updateStatus)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}

				// Ensure state did not change
				if task.Name != tt.wantName {
					t.Errorf("Expected Name to remain %q, got %q", tt.wantName, task.Name)
				}
				if task.Status != tt.wantStatus {
					t.Errorf("Expected Status to remain %d, got %d", tt.wantStatus, task.Status)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if task.Name != tt.wantName {
					t.Errorf("Expected Name to be %q, got %q", tt.wantName, task.Name)
				}
				if task.Status != tt.wantStatus {
					t.Errorf("Expected Status to be %d, got %d", tt.wantStatus, task.Status)
				}
			}
		})
	}
}

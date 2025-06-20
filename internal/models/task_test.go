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

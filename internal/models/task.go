package models

import (
	"errors"
	"strings"
)

// Task represents a task in our task management system
// Note: In production, consider using UUID for better security and distributed system compatibility
type Task struct {
	ID     int    `json:"id"`     // Unique identifier (use UUID in production)
	Name   string `json:"name"`   // Task name
	Status int    `json:"status"` // 0 = incomplete, 1 = completed
}

// NewTask creates a new Task with the given name and status.
// It returns an error if the name is empty or contains only whitespace,
// or if the status is not 0 (incomplete) or 1 (complete).
func NewTask(name string, status int) (*Task, error) {
	// Note: Input validation helps prevent injection attacks and ensures data consistency.
	// In production, consider additional validation like length limits, character filtering, etc.
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("task name cannot be empty")
	}

	// Validate status - must be 0 (incomplete) or 1 (completed)
	if status < 0 || status > 1 {
		return nil, errors.New("status must be 0 (incomplete) or 1 (completed)")
	}

	return &Task{
		Name:   name,
		Status: status,
	}, nil
}

// Update modifies the task with new name and status, applying validation.
// This method follows the fetch-modify-save pattern used in production systems.
func (t *Task) Update(name string, status int) error {
	// Validate name
	if strings.TrimSpace(name) == "" {
		return errors.New("task name cannot be empty")
	}

	// Validate status
	if status < 0 || status > 1 {
		return errors.New("status must be 0 (incomplete) or 1 (completed)")
	}

	// Apply changes
	t.Name = name
	t.Status = status
	return nil
}

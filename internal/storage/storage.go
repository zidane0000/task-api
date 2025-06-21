package storage

import (
	"os"
	"task-api/internal/models"
)

// TaskStorage defines the interface for task storage operations.
// This interface allows different storage implementations (in-memory, database, etc.)
// while keeping the business logic decoupled from storage details.
//
// TODO: Future enhancement - Add hybrid storage backend support
// - Auto-detect storage type based on environment (DATABASE_URL presence)
// - Use in-memory storage for development/testing (current implementation)
// - Use database backend for production when DATABASE_URL is configured
type TaskStorage interface {
	// Create stores a new task and assigns it a unique ID.
	// Returns the task with assigned ID or an error if creation fails.
	Create(task *models.Task) (*models.Task, error)

	// GetAll retrieves all tasks from storage.
	// Returns slice of tasks or error if retrieval fails.
	GetAll() ([]*models.Task, error)

	// GetByID retrieves a specific task by its ID.
	// Returns the task or error if not found or retrieval fails.
	GetByID(id int) (*models.Task, error)

	// Update modifies an existing task in storage.
	// Returns error if task doesn't exist or update fails.
	Update(task *models.Task) error

	// Delete removes a task from storage by ID.
	// Returns error if task doesn't exist or deletion fails.
	Delete(id int) error
}

// StoreBackend defines the type of storage backend
type StoreBackend string

const (
	BackendMemory   StoreBackend = "memory"
	BackendDatabase StoreBackend = "database"
)

// AutoDetectBackend automatically detects which backend to use based on environment
func AutoDetectBackend() StoreBackend {
	if os.Getenv("DATABASE_URL") != "" {
		return BackendDatabase
	}
	return BackendMemory
}

// NewTaskStorage creates a TaskStorage based on auto-detected backend
func NewTaskStorage() TaskStorage {
	backend := AutoDetectBackend()

	switch backend {
	case BackendDatabase:
		// TODO: Implement database storage
		panic("Database backend not implemented yet. Remove DATABASE_URL to use memory backend.")
	default:
		return NewInMemoryStorage()
	}
}

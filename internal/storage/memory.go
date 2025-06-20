package storage

import (
	"errors"
	"sync"
	"task-api/internal/models"
)

// InMemoryStorage implements TaskStorage interface using in-memory storage.
// It provides thread-safe operations for storing and retrieving tasks.
type InMemoryStorage struct {
	tasks  map[int]*models.Task // Map of ID to Task
	nextID int                  // Auto-incrementing ID counter
	mutex  sync.RWMutex         // Protects concurrent access
}

// NewInMemoryStorage creates a new in-memory storage instance.
// Returns a storage implementation ready for use.
func NewInMemoryStorage() TaskStorage {
	return &InMemoryStorage{
		tasks:  make(map[int]*models.Task),
		nextID: 1, // Start IDs from 1
	}
}

// Create stores a new task and assigns it a unique ID.
// Returns the task with assigned ID or an error if creation fails.
func (s *InMemoryStorage) Create(task *models.Task) (*models.Task, error) {
	if task == nil {
		return nil, errors.New("task cannot be nil")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Create a copy of the task with assigned ID
	newTask := &models.Task{
		ID:     s.nextID,
		Name:   task.Name,
		Status: task.Status,
	}

	// Store the task
	s.tasks[s.nextID] = newTask
	s.nextID++

	return newTask, nil
}

// GetAll retrieves all tasks from storage.
// Returns slice of tasks or error if retrieval fails.
func (s *InMemoryStorage) GetAll() ([]*models.Task, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	// Create slice to hold all tasks
	tasks := make([]*models.Task, 0, len(s.tasks))

	// Copy all tasks to slice
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetByID retrieves a specific task by its ID.
// Returns the task or error if not found or retrieval fails.
func (s *InMemoryStorage) GetByID(id int) (*models.Task, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	task, exists := s.tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}

	return task, nil
}

// Update modifies an existing task in storage.
// Returns error if task doesn't exist or update fails.
func (s *InMemoryStorage) Update(task *models.Task) error {
	if task == nil {
		return errors.New("task cannot be nil")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if task exists
	_, exists := s.tasks[task.ID]
	if !exists {
		return errors.New("task not found")
	}

	// Update the task
	s.tasks[task.ID] = task

	return nil
}

// Delete removes a task from storage by ID.
// Returns error if task doesn't exist or deletion fails.
func (s *InMemoryStorage) Delete(id int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Check if task exists
	_, exists := s.tasks[id]
	if !exists {
		return errors.New("task not found")
	}

	// Delete the task
	delete(s.tasks, id)
	return nil
}

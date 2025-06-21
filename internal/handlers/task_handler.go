package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"task-api/internal/models"
	"task-api/internal/storage"

	"github.com/go-chi/chi/v5"
)

// TaskHandler handles HTTP requests for task operations
type TaskHandler struct {
	storage storage.TaskStorage
}

// NewTaskHandler creates a new TaskHandler with the given storage
func NewTaskHandler(storage storage.TaskStorage) *TaskHandler {
	return &TaskHandler{
		storage: storage,
	}
}

// GetAllTasks handles GET /tasks - retrieve all tasks
func (h *TaskHandler) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.storage.GetAll()
	if err != nil {
		writeErrorResponse(w, ErrInternalServer)
		return
	}

	if err := writeJSONResponse(w, tasks, http.StatusOK); err != nil {
		writeErrorResponse(w, ErrInternalServer)
		return
	}
}

// CreateTask handles POST /tasks - create a new task
func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeErrorResponse(w, ErrInvalidJSON)
		return
	}

	newTask, err := models.NewTask(task.Name, task.Status)
	if err != nil {
		writeErrorResponse(w, ErrorResponse{
			Err:     err,
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Create task in storage
	createdTask, err := h.storage.Create(newTask)
	if err != nil {
		writeErrorResponse(w, ErrInternalServer)
		return
	}

	if err := writeJSONResponse(w, createdTask, http.StatusCreated); err != nil {
		writeErrorResponse(w, ErrInternalServer)
		return
	}
}

// UpdateTask handles PUT /tasks/{id} - update an existing task
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path using chi
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeErrorResponse(w, ErrInvalidTaskID)
		return
	}

	// Fetch existing task
	existingTask, err := h.storage.GetByID(id)
	if err != nil {
		writeErrorResponse(w, ErrTaskNotFound)
		return
	}

	// Parse input
	var input struct {
		Name   string `json:"name"`
		Status int    `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeErrorResponse(w, ErrInvalidJSON)
		return
	}

	// Apply changes with validation
	if err := existingTask.Update(input.Name, input.Status); err != nil {
		writeErrorResponse(w, ErrorResponse{
			Err:     err,
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		})
		return
	}

	// Save updated task
	if err := h.storage.Update(existingTask); err != nil {
		writeErrorResponse(w, ErrInternalServer)
		return
	}

	if err := writeJSONResponse(w, existingTask, http.StatusOK); err != nil {
		writeErrorResponse(w, ErrInternalServer)
		return
	}
}

// DeleteTask handles DELETE /tasks/{id} - delete a task
func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path using chi
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeErrorResponse(w, ErrInvalidTaskID)
		return
	}

	// Check if task exists
	_, err = h.storage.GetByID(id)
	if err != nil {
		writeErrorResponse(w, ErrTaskNotFound)
		return
	}

	// Delete task from storage
	if err := h.storage.Delete(id); err != nil {
		writeErrorResponse(w, ErrInternalServer)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

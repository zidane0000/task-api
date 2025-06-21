package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"task-api/internal/handlers"
	"task-api/internal/storage"
)

func main() {
	// Initialize storage
	taskStorage := storage.NewTaskStorage()
	log.Println("Storage initialized successfully")

	// Initialize handlers
	taskHandler := handlers.NewTaskHandler(taskStorage)

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)                    // Request logging
	r.Use(middleware.Recoverer)                 // Panic recovery
	r.Use(middleware.Timeout(60 * time.Second)) // Request timeout

	// Routes
	r.Get("/health", healthHandler)
	r.Route("/tasks", func(r chi.Router) {
		r.Get("/", taskHandler.GetAllTasks)
		r.Post("/", taskHandler.CreateTask)
		r.Put("/{id}", taskHandler.UpdateTask)
		r.Delete("/{id}", taskHandler.DeleteTask)
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Printf("Available endpoints:")
	log.Printf("  GET    /health       - Health check")
	log.Printf("  GET    /tasks        - Get all tasks")
	log.Printf("  POST   /tasks        - Create new task")
	log.Printf("  PUT    /tasks/{id}   - Update task")
	log.Printf("  DELETE /tasks/{id}   - Delete task")

	// Create HTTP server with proper timeouts for security
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

// healthHandler provides a simple health check endpoint
// Future: Add storage connectivity, uptime, version checks
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte(`{"status":"healthy","service":"task-api"}`)); err != nil {
		log.Printf("Error writing health check response: %v", err)
	}
}

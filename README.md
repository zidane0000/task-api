# Task API

A RESTful task management API.

## Documentation

- [Requirements](docs/requirements.md) - Detailed project requirements and specifications

## Project Structure

```text
task-api/
├── cmd/server/          # Main application entry point
├── docs/                # Project documentation
├── internal/
│   ├── handlers/        # HTTP handlers/controllers
│   ├── models/          # Data models and structs
│   └── storage/         # Data storage layer
└── tests/               # Test files
```

## Getting Started

### Prerequisites

- Go 1.23+

### Environment Variables

- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - Database connection string (optional also unimplemented, uses in-memory storage if not set)

### API Endpoints

- `GET /health` - Health check endpoint
- `GET /tasks` - Retrieve all tasks
- `POST /tasks` - Create a new task
- `PUT /tasks/{id}` - Update an existing task
- `DELETE /tasks/{id}` - Delete a task

### Testing

```bash
# Run all tests(include e2e, so make sure the server is running)
go test ./...
```

### Running the Application

#### Local Development

```bash
# Clone the repository
git clone https://github.com/zidane0000/task-api
cd task-api

# Install dependencies
go mod download

# Run the server
go run cmd/server/main.go
```

The API will be available at `http://localhost:8080`

```bash
# Health Check if the API is running
curl http://localhost:8080/health    
{"status":"healthy","service":"task-api"}

# Retrieve all tasks (initially empty)
curl http://localhost:8080/tasks 
[]

# Create a New Task
curl -X POST http://localhost:8080/tasks -H "Content-Type: application/json" -d '{"name":"Test","status":0}'
{"id":1,"name":"Test","status":0}

# Update task with ID
curl -X PUT http://localhost:8080/tasks/1 -H "Content-Type: application/json" -d '{"name":"Test - DONE!","status":1}'
{"id":1,"name":"Test - DONE!","status":1}

# Delete task with ID
curl -X DELETE http://localhost:8080/tasks/1
```

#### Docker Deployment

```bash
# Build the Docker image
docker build -t task-api .

# Run the container
docker run -p 8080:8080 task-api
```

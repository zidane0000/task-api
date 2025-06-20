# Task API

A RESTful task management API.

## Requirements

### API Endpoints

The API must provide the following endpoints:

- `GET /tasks` - Retrieve all tasks
- `POST /tasks` - Create a new task
- `PUT /tasks/{id}` - Update an existing task
- `DELETE /tasks/{id}` - Delete a task

### Task Model

Each task must have the following fields:

- `name` (string) - The task name
- `status` (integer) - Task completion status
  - `0` = incomplete
  - `1` = completed

### Technical Requirements

- **Language**: Go 1.18+
- **Testing**: Unit tests with good coverage
- **Containerization**: Dockerfile included
- **Repository**: GitHub ready
- **Storage**: In-memory data storage

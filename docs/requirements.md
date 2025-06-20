# Task API Requirements

## API Endpoints

The API must provide the following endpoints:

- `GET /tasks` - Retrieve all tasks
- `POST /tasks` - Create a new task
- `PUT /tasks/{id}` - Update an existing task
- `DELETE /tasks/{id}` - Delete a task

## Task Model

Each task must have the following fields:

- name
  - type: string
  - description:task name
- status
  - type: integer
  - enum:[0,1]
  - description: 0 represents an incomplete task, while 1 represents a completed task

## Technical Requirements

- Runtime environment should be Go 1.18+
- Provides unit tests
- Provides Dockerfile to run API in Docker
- Manage the codebase on Github and provide us with the repository link
- For data storage, you can use any in-memory mechanism

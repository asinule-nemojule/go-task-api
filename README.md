# Task Manager REST API in Go

RESTful API for task management built in Go.

## Endpoints

- `GET /tasks` - List all tasks
- `GET /tasks/{id}` - Get single task
- `POST /tasks` - Create task
- `DELETE /tasks/{id}` - Delete task
- `GET /health` - Health check

## Run locally

- `go run main.go`




## Run with Docker

- `docker build -t go-task-api .`
- `docker run -p 8080:8080 go-task-api`

# CI/CD

GitHub Actions automatically runs tests on every push.

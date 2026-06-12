# goProject

A Go REST API for workouts, users, and token-based authentication.

## Requirements

- Go 1.26+
- Docker (for PostgreSQL via `docker-compose`)

## Run locally

1. Start Postgres:
   ```bash
   docker-compose up -d db
   ```
2. Start the API:
   ```bash
   go run .
   ```

The server runs on `:8080` by default.

## Test setup

Tests in `internal/store` use a separate Postgres instance on port `5433`.

```bash
docker-compose up -d test_db
go test ./...
```

## API routes

- `GET /health`
- `POST /users`
- `POST /tokens/authentication`
- `GET /workouts/{id}`
- `POST /workouts` (authenticated)
- `PUT /workouts/{id}` (authenticated)
- `DELETE /workouts/{id}` (authenticated)

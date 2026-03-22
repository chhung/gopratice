# lab2

This lab is a small standalone Go HTTP service that uses a repository, service, and middleware structure.

## Run

```bash
set APP_PORT=8080
set SHUTDOWN_TOKEN=your-token
go run ./cmd/app
```

## Configuration

- `APP_PORT`: HTTP listen port, default `8080`
- `SHUTDOWN_TOKEN`: bearer token for `POST /admin/shutdown`
- `HTTP_READ_HEADER_TIMEOUT`: default `2s`
- `HTTP_READ_TIMEOUT`: default `5s`
- `HTTP_WRITE_TIMEOUT`: default `10s`
- `HTTP_IDLE_TIMEOUT`: default `60s`
- `HTTP_SHUTDOWN_TIMEOUT`: default `5s`

## Endpoints

- `GET /messages`
- `POST /messages`
- `GET /healthz`
- `POST /admin/shutdown`

## Structure

- `cmd/app`: application entrypoint
- `internal/app`: HTTP server bootstrap and handlers
- `internal/config`: environment-based application configuration
- `internal/middleware`: cross-cutting HTTP concerns
- `internal/repository`: data access implementation
- `internal/service`: business rules and validation
- `internal/model`: request and response types
- `test`: integration-test placeholder

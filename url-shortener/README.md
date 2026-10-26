
# Go URL Shortener

A URL shortener built with Go, PostgreSQL, and Redis. The application exposes an HTTP API for creating short URLs and redirecting requests to their original destinations. PostgreSQL provides persistent storage, while Redis caches URL mappings to reduce repeated database lookups.

The project uses a layered application architecture, includes unit tests, and runs with Docker Compose.

## Features

- Create short URLs through `POST /shorten`.
- Redirect short codes using HTTP `302 Found`.
- Persist URL mappings in PostgreSQL.
- Cache URL lookups in Redis with a 15-minute TTL.
- Fall back to PostgreSQL when a cache lookup fails.
- Validate submitted URLs and malformed JSON requests.
- Retry short-code generation when a database uniqueness conflict occurs.
- Unit tests for code generation, service logic, caching, and HTTP handlers.
- Run the API, PostgreSQL, and Redis with Docker Compose.
- Initialize the database table automatically for fresh PostgreSQL data directories.
- Persist PostgreSQL data in a named Docker volume.

## Tech Stack

- **Language:** Go
- **HTTP server:** Go `net/http`
- **Database:** PostgreSQL
- **Cache:** Redis
- **PostgreSQL driver:** `github.com/lib/pq`
- **Redis client:** `github.com/redis/go-redis/v9`
- **Containerization:** Docker and Docker Compose

## Architecture

The request flow for URL creation and redirection is:

1. A client sends a request to the Go HTTP API.
2. The handler validates the request and passes work to the service layer.
3. For URL creation, the service generates a short code and stores the mapping in PostgreSQL.
4. For redirection, the service checks Redis first.
5. If Redis has no mapping or is unavailable, the service queries PostgreSQL.
6. A successful database lookup is cached in Redis on a best-effort basis.
7. The API redirects the client to the original URL with HTTP `302 Found`.

## Project Structure

- `internal/`
  - `cache/`
    - `redis_cache.go`
  - `handler/`
    - `url_handler.go`
    - `url_handler_test.go`
  - `repository/`
    - `url_repository.go`
  - `service/`
    - `code_generator.go`
    - `code_generator_test.go`
    - `url_service.go`
    - `url_service_test.go`
- `.dockerignore`
- `.env.example`
- `compose.yaml`
- `Dockerfile`
- `go.mod`
- `go.sum`
- `init.sql`
- `main.go`
- `README.md`

## Component Responsibilities

- **`main.go`:** Initializes dependencies and configures the HTTP server.
- **`internal/handler/`:** Handles HTTP requests, input validation, and responses.
- **`internal/service/`:** Implements URL-shortening logic, code-generation orchestration, and cache/database lookup behavior.
- **`internal/repository/`:** Handles PostgreSQL operations.
- **`internal/cache/`:** Implements Redis caching with TTL support.
- **`Dockerfile`:** Builds the Go application into a minimal runtime image.
- **`compose.yaml`:** Configures the API, PostgreSQL, Redis, health checks, ports, and persistent storage.
- **`init.sql`:** Creates the URL mappings table when PostgreSQL initializes a fresh data directory.
- **`.env.example`:** Provides a template for local environment configuration.

## Prerequisites

For the Docker-based setup:

- Docker Desktop with Docker Compose
- Git

For local development without Docker:

- A Go version compatible with `go.mod`
- PostgreSQL
- Redis

## Getting Started with Docker Compose

### 1. Clone the repository

```bash
git clone YOUR_REPOSITORY_URL
cd url-shortener
```

Replace `YOUR_REPOSITORY_URL` with your actual repository URL.

### 2. Configure environment variables

Copy the example environment file and open it:

```powershell
Copy-Item .env.example .env
notepad .env
```

Replace the placeholder PostgreSQL password with a unique local development password. Do not commit `.env` or publish real credentials.

### 3. Start the services

Build the API image and start the services:

```bash
docker compose up -d --build
```

Check service status:

```bash
docker compose ps
```

Inspect API logs:

```bash
docker compose logs api --tail 50
```

The API should be available at `http://localhost:8080`.

### 4. Service ports

| Service | Host address | Purpose |
|---|---|---|
| API | `localhost:8080` | HTTP endpoints |
| PostgreSQL | `localhost:5433` | Persistent storage |
| Redis | `localhost:6380` | URL lookup cache |

The API connects to PostgreSQL and Redis using their Compose service names and internal container ports. The host ports are for connecting from your computer.

### 5. Database initialization and persistence

The `init.sql` file creates the `urls` table automatically when PostgreSQL initializes a fresh data directory.

The database uses a named Docker volume, `postgres_data`, to persist data across container restarts and recreation. Initialization scripts are not rerun against an already-initialized PostgreSQL data directory, so changing `init.sql` will not automatically update an existing database.

**Warning:** `docker compose down -v` removes the Compose-managed database volume and its stored data. Do not run it unless you intentionally want to delete that data.

## Running Locally Without Docker

Create a PostgreSQL database named `url_shortener`, then execute the schema in `init.sql`.

Ensure PostgreSQL and Redis are running. Configure the environment variables in PowerShell, replacing the placeholders with your own database credentials:

```powershell
$env:URL_SHORTENER_DB_URL = "host=localhost port=5432 user=YOUR_USER password=YOUR_PASSWORD dbname=url_shortener sslmode=disable"
$env:REDIS_ADDR = "localhost:6379"
```

Download dependencies and run the server:

```bash
go mod download
go run .
```

The HTTP server listens on port `8080`. For local development, the default Redis address is `localhost:6379` when `REDIS_ADDR` is not configured.

## API Reference

### Create a Short URL

**Endpoint:** `POST /shorten`

**Content-Type:** `application/json`

Request body:

```json
{
  "url": "https://example.com/products?id=123"
}
```

Successful response â€” `201 Created`:

```json
{
  "short_code": "gJXpSZ",
  "short_url": "http://localhost:8080/gJXpSZ"
}
```

The short code above is an example; actual codes are generated dynamically.

PowerShell example:

```powershell
$body = @{
    url = "https://example.com/products?id=123"
} | ConvertTo-Json

$response = Invoke-RestMethod `
    -Method Post `
    -Uri http://localhost:8080/shorten `
    -ContentType "application/json" `
    -Body $body

$response
```

The API validates submitted URLs and rejects malformed JSON, unsupported URL formats, and requests using the wrong HTTP method.

### Redirect to the Original URL

**Endpoint:** `GET /{shortCode}`

For example:

```bash
curl.exe -i http://localhost:8080/gJXpSZ
```

A successful redirect returns `302 Found` with a `Location` header pointing to the original URL:

```http
HTTP/1.1 302 Found
Location: https://example.com/products?id=123
```

An unknown short code returns `404 Not Found`.

## Caching Behavior

Redis stores short-code-to-original-URL mappings with a 15-minute TTL (900 seconds).

The lookup flow is:

1. Attempt to retrieve the mapping from Redis.
2. Return the cached URL if available.
3. Query PostgreSQL if the cache misses or fails.
4. Populate Redis after a successful database lookup, on a best-effort basis.
5. Return a not-found response if PostgreSQL has no matching mapping.

A Redis cache-write failure does not prevent a successful database lookup from returning the original URL.

## Manual Redis Verification

With Docker Compose running, inspect a cached URL:

```bash
docker compose exec redis redis-cli GET YOUR_SHORT_CODE
```

Check the remaining TTL:

```bash
docker compose exec redis redis-cli TTL YOUR_SHORT_CODE
```

Replace `YOUR_SHORT_CODE` with an actual short code returned by the API.

## Testing

Run all Go tests:

```bash
go test ./...
```

Run service tests with verbose output:

```bash
go test ./internal/service -v
```

Run static analysis:

```bash
go vet ./...
```

The test suite covers short-code generation, duplicate-code retries, service errors, cache hits and misses, cache-failure fallback, request validation, and redirect handling.

## Development and Troubleshooting

Check running services:

```bash
docker compose ps
```

View service logs:

```bash
docker compose logs api
docker compose logs postgres
docker compose logs redis
```

Restart the services:

```bash
docker compose restart
```

Stop the services without removing the named database volume:

```bash
docker compose down
```

If PostgreSQL reports an authentication error, verify that the credentials in `.env` match the database user's configured credentials. Updating `POSTGRES_PASSWORD` in `.env` alone does not change the password of a database user in an already-initialized data directory.

## Future Improvements

- Database migrations and versioned schema management
- Structured logging, request tracing, and graceful server shutdown
- Metrics for cache hits, cache misses, and database fallbacks
- Rate limiting and abuse prevention
- Integration tests against real PostgreSQL and Redis instances
- Production deployment configuration and monitoring
- Custom short codes, expiration policies, and URL analytics

## Learning Objectives

- HTTP routing and status codes in Go
- Layered application architecture
- Relational persistence using PostgreSQL
- Caching and TTL management with Redis
- Error handling and dependency injection
- Unit testing and static analysis
- Docker image building and multi-service orchestration
- Persistent storage and database initialization
- Graceful degradation when the cache is unavailable



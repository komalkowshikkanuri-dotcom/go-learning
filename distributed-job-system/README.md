# Distributed Job Processing System

A backend project built with Go, PostgreSQL, and Redis to explore asynchronous job processing, worker concurrency, retry handling, and containerized deployment.

## Overview

The system accepts job requests through an HTTP API, persists job metadata in PostgreSQL, and places job IDs in a Redis queue. A worker consumes queued jobs, tracks execution attempts, and updates job statuses.

### Architecture

```text
Client
  |
  v
Go HTTP API
  |
  +------> PostgreSQL
  |         - Users
  |         - Jobs
  |         - Job attempts
  |
  +------> Redis queue
                |
                v
          Go Worker
                |
                +----> PostgreSQL
                       - Update job status
                       - Record attempts
```

## Tech Stack

- **Go** - API and worker implementation
- **PostgreSQL** - persistent job and attempt records
- **Redis** - job queue
- **Docker** - containerized application services
- **Docker Compose** - local multi-service environment
- **Go testing** - automated tests for worker, configuration, and service logic

## Features

- HTTP endpoints for submitting and retrieving jobs
- Persistent job records in PostgreSQL
- Redis-backed job queue
- Worker-based asynchronous processing
- Job lifecycle tracking: pending, processing, completed, and failed
- Execution-attempt tracking with timestamps
- Retry handling with a configured maximum of three attempts
- Bounded worker concurrency, with up to three active jobs per worker process
- Environment-based application configuration
- Docker Compose setup for the API, worker, PostgreSQL, and Redis
- Automated tests for key application logic

## Project Structure

```text
distributed-job-system/
|-- cmd/
|   |-- api/
|   |   `-- main.go
|   `-- worker/
|       |-- main.go
|       |-- concurrency_test.go
|       `-- process_job_test.go
|-- internal/
|   |-- config/
|   |-- handler/
|   |-- models/
|   |-- redis/
|   |-- repository/
|   `-- service/
|-- migrations/
|   `-- 001_init.sql
|-- .dockerignore
|-- .gitignore
|-- compose.yaml
|-- Dockerfile
|-- go.mod
`-- go.sum
```

## Getting Started

### Prerequisites

- Docker Desktop with the Docker Compose plugin
- Git, if cloning the repository
- Go, if running or testing the application outside Docker

### 1. Configure environment variables

Create a local `.env` file containing the database password expected by the Compose configuration:

```dotenv
DB_PASSWORD=your_local_password
```

Use a local development password and do not commit `.env` or share its contents.

### 2. Start the services

From the project root, run:

```powershell
docker compose up -d --build
```

Check the service status:

```powershell
docker compose ps
```

View logs:

```powershell
docker compose logs -f api worker
```

The API is exposed at `http://localhost:8080`. PostgreSQL and Redis are configured as Compose services.

### 3. Create a test user

The current API does not include a user-creation endpoint. For local testing, insert a user into the database:

```powershell
docker compose exec db psql -U postgres -d job_system
```

Then execute:

```sql
INSERT INTO users (name) VALUES ('Test User');
```

Exit `psql` with:

```sql
\q
```

If your Compose configuration uses a different PostgreSQL username or database name, use the configured values instead.

## API Endpoints

### `GET /`

Returns a basic API response.

```powershell
Invoke-RestMethod http://localhost:8080/
```

### `POST /job`

Submits a job.

Example request:

```powershell
$body = @{
    user_id = 1
    type    = "docker-success-test"
} | ConvertTo-Json

Invoke-RestMethod `
    -Uri http://localhost:8080/job `
    -Method Post `
    -ContentType "application/json" `
    -Body $body
```

The response includes the job ID, type, and initial status.

The example job type `docker-success-test` exercises the implementation's simulated success path. Other job types currently exercise the simulated failure and retry path; this is not yet a general-purpose job execution engine.

### `GET /jobs`

Lists jobs.

```powershell
Invoke-RestMethod http://localhost:8080/jobs
```

### `GET /jobs/{id}`

Retrieves a job by ID.

```powershell
Invoke-RestMethod http://localhost:8080/jobs/1
```

## Job Lifecycle

The worker updates job statuses as processing proceeds:

```text
Pending
   |
   v
Processing
   |
   +----> Completed
   |
   +----> Failed
              |
              +----> Pending (if another attempt is allowed)
```

Execution attempts are recorded separately in the `job_attempts` table, allowing the application to retain attempt status and timing information.

## Concurrency and Retries

The worker limits the number of concurrently active jobs within a single process. The current limit is three active jobs per worker process.

Retry handling allows another attempt when the configured maximum has not been reached. The current implementation uses simulated success and failure behavior to exercise these paths.

**Important implementation limitations:**

- Concurrency limits are local to each worker process, not global across all worker instances.
- The Redis queue uses destructive consumption; a worker crash after removing a job but before completing it can result in lost work.
- PostgreSQL persistence and Redis enqueueing are separate operations and are not atomic.
- Graceful shutdown and reliable recovery of interrupted jobs are not yet implemented.
- Job execution is simulated rather than dispatching arbitrary real workloads.

These are known areas for future improvement.

## Database Schema

The initial migration, `migrations/001_init.sql`, defines:

- `users` - user records
- `jobs` - job metadata and current status
- `job_attempts` - individual execution-attempt records and timestamps

Foreign keys associate jobs with users and attempts with jobs.

## Running Tests

Run the test suite from the project root:

```powershell
go test ./... -count=1
```

The suite covers worker processing, configuration, and service logic. A passing test suite does not by itself verify every production failure scenario.

## Future Improvements

- Reliable job acknowledgement and recovery
- Graceful worker shutdown
- Atomic job publishing through a transactional outbox
- Authentication and authorization
- More comprehensive integration and failure-injection tests
- Metrics, structured logging, and health endpoints
- Real job handlers and configurable retry policies

## Purpose

This project is a learning exercise in backend engineering and distributed systems concepts, including asynchronous processing, persistence, concurrency, retries, and container orchestration.
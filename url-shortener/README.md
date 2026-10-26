\# URL Shortener in Go



A URL shortener built with Go, PostgreSQL, and Redis. The application exposes an HTTP API to create short URLs and redirects requests to their original destinations. Redis is used to cache URL mappings, while PostgreSQL provides persistent storage.



\## Features



\- Create short URLs through a REST-style HTTP endpoint.

\- Redirect short codes to their original URLs using HTTP `302 Found`.

\- Persist URL mappings in PostgreSQL.

\- Cache URL mappings in Redis with a configurable time-to-live (TTL).

\- Fall back to PostgreSQL when a cache lookup fails.

\- Validate submitted URLs and handle malformed JSON requests.

\- Retry short-code generation when a database uniqueness conflict occurs.

\- Unit tests for URL generation, service logic, cache behavior, and HTTP handlers.



\## Tech Stack



\- \*\*Language:\*\* Go

\- \*\*HTTP server:\*\* Go `net/http`

\- \*\*Database:\*\* PostgreSQL

\- \*\*Cache:\*\* Redis

\- \*\*PostgreSQL driver:\*\* `github.com/lib/pq`

\- \*\*Redis client:\*\* `github.com/redis/go-redis/v9`



\## Project Structure



```text

url-shortener/

├── main.go

├── go.mod

├── go.sum

└── internal/

&#x20;   ├── cache/

&#x20;   │   └── redis\_cache.go

&#x20;   ├── handler/

&#x20;   │   ├── url\_handler.go

&#x20;   │   └── url\_handler\_test.go

&#x20;   ├── repository/

&#x20;   │   └── url\_repository.go

&#x20;   └── service/

&#x20;       ├── code\_generator.go

&#x20;       ├── code\_generator\_test.go

&#x20;       ├── url\_service.go

&#x20;       └── url\_service\_test.go

```



\### Responsibilities



\- `main.go`: Initializes dependencies and configures the HTTP server.

\- `internal/handler/`: Handles HTTP requests, validates inputs, and produces HTTP responses.

\- `internal/service/`: Implements URL-shortening logic, code generation orchestration, and cache/database lookup behavior.

\- `internal/repository/`: Handles PostgreSQL operations.

\- `internal/cache/`: Implements Redis-based caching with TTL support.



\## Architecture



```text

&#x20;                Client

&#x20;                  |

&#x20;                  v

&#x20;             HTTP Handler

&#x20;                  |

&#x20;                  v

&#x20;              URL Service

&#x20;                  |

&#x20;            Redirect lookup

&#x20;                  |

&#x20;                  v

&#x20;              Redis Cache

&#x20;               /       \\

&#x20;          Cache hit   Cache miss

&#x20;              |           |

&#x20;              |           v

&#x20;              |       PostgreSQL

&#x20;              |           |

&#x20;              |           v

&#x20;              |       Populate cache

&#x20;              |           |

&#x20;              +-----+-----+

&#x20;                    |

&#x20;                    v

&#x20;               HTTP 302

&#x20;                    |

&#x20;                    v

&#x20;            Original destination

```



For URL creation, the service generates a short code and stores the mapping in PostgreSQL. During redirects, the service checks Redis first and falls back to PostgreSQL if the cache lookup fails. Successful database lookups are written back to Redis on a best-effort basis.



\## Prerequisites



Install or configure the following:



\- Go compatible with the version declared in `go.mod`

\- PostgreSQL

\- Redis

\- Git



The application defaults to Redis at `localhost:6379` unless `REDIS\_ADDR` is configured.



\## Database Setup



Create a PostgreSQL database named `url\_shortener`.



Create the URL mappings table:



```sql

CREATE TABLE urls (

&#x20;   id BIGSERIAL PRIMARY KEY,

&#x20;   short\_code VARCHAR(10) NOT NULL UNIQUE,

&#x20;   original\_url TEXT NOT NULL,

&#x20;   created\_at TIMESTAMPTZ NOT NULL DEFAULT NOW()

);

```



The table stores each short code, its original URL, and the creation timestamp. The unique constraint on `short\_code` prevents duplicate mappings.



\## Configuration



Configure the PostgreSQL connection using the `URL\_SHORTENER\_DB\_URL` environment variable.



Example for a local development database:



```text

postgres://YOUR\_USER:YOUR\_PASSWORD@localhost:5432/url\_shortener?sslmode=disable

```



Set the environment variables in PowerShell:



```powershell

$env:URL\_SHORTENER\_DB\_URL = "postgres://YOUR\_USER:YOUR\_PASSWORD@localhost:5432/url\_shortener?sslmode=disable"

$env:REDIS\_ADDR = "localhost:6379"

```



Replace the database username and password with your local credentials. Do not commit real passwords or other secrets to Git.



`REDIS\_ADDR` is optional if Redis is accessible at the default address, `localhost:6379`.



\## Running the Application



From the project root:



```powershell

go mod download

go run .

```



The HTTP server listens on port `8080`.



Keep PostgreSQL and Redis running while testing the full application. The application is designed to fall back to PostgreSQL when Redis lookups fail, but PostgreSQL must be available for persistent storage and cache misses.



\## API Reference



\### 1. Create a Short URL



\*\*Endpoint:\*\* `POST /shorten`



\*\*Content-Type:\*\* `application/json`



Request:



```json

{

&#x20; "url": "https://example.com/products?id=123"

}

```



Successful response — `201 Created`:



```json

{

&#x20; "short\_code": "gJXpSZ",

&#x20; "short\_url": "http://localhost:8080/gJXpSZ"

}

```



The example short code is illustrative; the application generates codes dynamically.



Example using PowerShell:



```powershell

$body = @{

&#x20;   url = "https://example.com/products?id=123"

} | ConvertTo-Json



$response = Invoke-RestMethod `

&#x20;   -Method Post `

&#x20;   -Uri http://localhost:8080/shorten `

&#x20;   -ContentType "application/json" `

&#x20;   -Body $body



$response

```



The API validates the submitted URL and rejects malformed JSON, unsupported URL formats, and requests using the wrong HTTP method.



\### 2. Redirect to the Original URL



\*\*Endpoint:\*\* `GET /{shortCode}`



Example:



```powershell

curl.exe -i http://localhost:8080/gJXpSZ

```



Expected response:



```text

HTTP/1.1 302 Found

Location: https://example.com/products?id=123

```



The browser or HTTP client can then follow the redirect to the original destination.



An unknown short code returns `404 Not Found`.



\## Caching Behavior



Redis stores short-code-to-original-URL mappings to avoid repeated database reads.



The configured cache TTL is \*\*15 minutes (900 seconds)\*\*.



The redirect lookup follows this sequence:



1\. Attempt to retrieve the mapping from Redis.

2\. If the mapping is available, return it without querying PostgreSQL.

3\. If the cache lookup misses or fails, query PostgreSQL.

4\. If PostgreSQL returns the mapping, attempt to populate Redis.

5\. If PostgreSQL cannot find the mapping, return a not-found response through the HTTP handler.



Cache writes are best-effort: a Redis write failure does not prevent a successful database lookup from returning the original URL.



\## Testing



Run all tests:



```powershell

go test ./...

```



Run service tests with verbose output:



```powershell

go test ./internal/service -v

```



Run static analysis:



```powershell

go vet ./...

```



The existing tests cover short-code generation, duplicate-code retries, service errors, cache hits, cache misses, cache failure fallback, request validation, and redirect handling.



\### Manual Redis Verification



After creating and redirecting a short URL, check the cache from a terminal where `redis-cli` is available:



```bash

redis-cli GET YOUR\_SHORT\_CODE

redis-cli TTL YOUR\_SHORT\_CODE

```



The first command should return the original URL if it has been cached. The second returns the remaining TTL in seconds.



\## Current Scope and Future Improvements



Potential next improvements include:



\- Docker Compose configuration for PostgreSQL, Redis, and the application.

\- Database migration files and repeatable schema setup.

\- Structured logging and request tracing.

\- Context propagation and graceful server shutdown.

\- Metrics for cache hits, cache misses, and database fallbacks.

\- Rate limiting and abuse prevention.

\- Integration tests using real PostgreSQL and Redis instances.

\- Deployment and production configuration.



\## Learning Objectives



This project demonstrates practical backend engineering concepts:



\- HTTP routing and status codes in Go.

\- Layered application architecture.

\- Relational persistence using PostgreSQL.

\- Caching and TTL management with Redis.

\- Error handling and dependency injection.

\- Unit testing and static analysis.

\- Graceful degradation when the cache is unavailable.






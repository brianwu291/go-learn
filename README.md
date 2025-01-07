# Go Learning Project

REST API implementation with rate limiting, caching, and external API integration.

## Project Structure

```
.
├── cache/              # Cache abstractions
├── constants/          # Application constants
├── db/                 # Database related code
│   ├── migrations/     # SQL migration files
│   ├── postgres/       # PostgreSQL client and configurations
│   ├── redis/          # Redis client and configurations
│   └── scripts/        # Database management scripts
├── handlers/
│   ├── fakestore/     # Fake store API handlers
│   └── financial/     # Financial calculation handlers
├── httpclient/        # HTTP client wrapper
├── middlewares/
│   └── ratelimiter/
├── repos/             # Repository layer
│   └── fakestore/     # Fake store API integration
├── services/          # Business logic
│   ├── fakestore/
│   └── financial/
├── types/             # Shared types
└── utils/             # Utility functions
```

## Features

- PostgreSQL Integration: Robust database layer with connection pooling and migrations
- Redis Caching: High-performance caching layer for data access and rate limiter
- Financial calculations with rate limiting
- Fake Store API integration with Redis caching
- Clean architecture with DI
- Concurrent API requests
- Unit tests and coverage reporting
- Docker Support: Containerized development environment

## Technical Stack

- Golang 1.23.4
- Gin web framework
- Redis 7
- PostgreSQL 15
- Environment variable configuration

## Architecture

- Handler → Service → Repository pattern
- Interface-based dependency injection
- Redis for caching and rate limiting
- Concurrent operations with goroutines
- Error handling abstractions

## Setup

Can either choose Mac or Docker

### With Mac:

1. Install dependencies:
   `$ go mod download`

2. Set up environment variables in .env:

```
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=""
REDIS_DB=0
DB_HOST=postgres
DB_MAX_POOL_CONS=200
DB_PORT=5432
DB_USER=your_username
DB_NAME=your_db_name
DB_PASSWORD=your_db_password
```

3. Run the application:
   `$ go run main.go` or `$ air` for hot reload

### With Docker:

1. Go to `/development-container` folder
2. Run the whole application:
   `docker compose up --build`
   (This include `air` command, so hot-reload supported)

## Testing

1. Run all tests only:
   `$ go test ./...`

2. Run all tests with coverage result:
   `$ go test -coverprofile=coverage.out ./...` (get the coverage file)
   or
   `$ go test -v -cover ./...`

3. Check the test coverage on browser after run `$ go test -coverprofile=coverage.out ./...`:
   `$ go tool cover -html=coverage.out `

## Test rate limiting:

`$ chmod +x shellscripts/ratelimiterchecker.sh`
or
`$ . ./shellscripts/ratelimiterchecker.sh`

## Database Migration

go to `/db` folder and `$. ./scripts/migrate.sh up`

## API Endpoints

### Financial

- `POST /calculate`: Financial calculations
  - Rate limit: 100 req/5min
  - Payload: `{"revenue": 100, "expenses": 50, "taxRate": 0.3}`

### Fake Store

- `GET /fake-store/all/categories`: Get all categories

  - Redis caching

- `GET /fake-store/all/categories/products`: Get all products with category filtering
  - Redis caching
  - Concurrent category fetching

## Learning Goals

- Golang syntax and patterns
- Clean code architecture
- Middleware implementation
- Redis integration
- Postgres integration
- Rate limiting concepts
- Unit testing in Go
- Shell scripting

## License

See LICENSE file.

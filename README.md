# Go Learning Project

A hands-on project to learn Golang, implementing a financial calculation service with Redis-based rate limiting.

## Project Structure

```
.
├── cache/             # Cache interface definitions
├── constants/         # Application constants
├── db/                # Database implementations
│   └── redis/         # Redis client implementation
├── handlers/          # HTTP request handlers
│   └── financial/     # Financial calculation endpoints
├── middlewares/       # HTTP middleware components
│   └── ratelimiter/   # Redis-based rate limiting
├── services/          # Business logic layer
│   └── financial/     # Financial calculation service
├── shellscripts/      # Utility scripts
│   └── ratelimiterchecker.sh  # Rate limit testing script
├── types/             # Common type definitions
└── main.go            # Application entry point
```

## Features

- Financial calculations API
- Redis-based rate limiting
- Clean architecture pattern
- Unit tests for handlers and services
- Shell script for load testing

## Technical Stack

- Golang
- Gin web framework
- Redis for rate limiting
- Environment variable configuration

## Setup

1. Install dependencies:
`$ go mod download`

2. Set up environment variables in .env:
```
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=""
REDIS_DB=0
```

3. Run the application:
`$ go run main.go` or `$ air` for hot reload


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

## API Endpoints
- `POST /calculate`: Financial calculation endpoint
  - With rate limit 100 requests per 5 minutes
  - example req payload:
    ```
    {
      "revenue": 100,
      "expenses": 50,
      "taxRate": 0.3
    }
    ```
- `GET /ping`: Health check endpoint
  - With rate limit: 5 requests per 20 secs

## Learning Goals

- Golang syntax and patterns
- Clean code architecture
- Middleware implementation
- Redis integration
- Rate limiting concepts
- Unit testing in Go
- Shell scripting

## License
See LICENSE file.
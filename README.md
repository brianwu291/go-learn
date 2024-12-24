# Go Learning Project

REST API implementation with rate limiting, caching, and external API integration.

## Project Structure

```
.
├── cache/              # Cache abstractions
├── constants/          # Application constants
├── db/                 # Database clients
│   └── redis/
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
└── types/             # Shared types
```

## Features

- Financial calculations with rate limiting
- Fake Store API integration with Redis caching
- Clean architecture with DI
- Concurrent API requests
- Unit tests and coverage reporting

## Technical Stack

- Golang
- Gin web framework
- Redis for rate limiting
- Environment variable configuration

## Architecture
- Handler → Service → Repository pattern
- Interface-based dependency injection
- Redis for caching and rate limiting
- Concurrent operations with goroutines
- Error handling abstractions

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
### Financial
- `POST /calculate`: Financial calculations
  - Rate limit: 100 req/5min
  - Payload: `{"revenue": 100, "expenses": 50, "taxRate": 0.3}`

### Fake Store
- `GET /products`: Get all products with category filtering
  - Redis caching
  - Concurrent category fetching

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
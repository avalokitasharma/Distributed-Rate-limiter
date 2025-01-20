# HTTP Rate Limiter

A distributed rate limiting system implemented in Go using microservices architecture. This system provides HTTP-based rate limiting capabilities with configurable time windows and Redis-backed storage.

## Features

- HTTP-based rate limiting service
- Configurable time window support (per second, minute, hour, day)
- Configurable rate limits per API endpoint
- Real-time rate limit checking
- Sliding window rate limiting implementation
- Redis as storage for distributed rate limiting
- Microservices architecture with separate admin and checker clients
- Docker support for easy deployment

## System Architecture

The system consists of three main components:

1. **Rate Limiter Server**
   - Core rate limiting service
   - Handles rate limit storage and verification
   - Exposes HTTP endpoints for administration and checking
   - Interfaces with Redis for distributed storage

2. **Admin Client**
   - Sets and manages rate limit rules
   - Configures limits for different time windows
   - Communicates with server via HTTP

3. **Rate limit Checker Client**
   - Verifies if requests are allowed
   - Provides real-time rate limit checking
   - Returns response with allow/deny status

## Requirements

- Go 1.21+
- Docker and Docker Compose
- Redis

## Project Structure

```
rate-limiter/
├── client/
│   ├── admin/
│   │   └── main.go
│   └── rateLimitChecker/
│       └── main.go
├── server/
│   ├── api/
│   │   ├── handlers.go
│   │   └── routes.go
│   ├── config/
│   │   └── config.go
│   ├── models/
│   │   └── models.go
│   ├── redis/
│   │   └── client.go
│   └── main.go
├── docker-compose.yml
├── Dockerfile.admin
├── Dockerfile.checker
├── Dockerfile.server
└── README.md
```

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/rate-limiter.git
cd rate-limiter
```

2. Build and run using Docker Compose:
```bash
docker-compose up --build
```

This will start:
- Redis server on port 6379
- Rate limiter server on port 8080
- Admin and checker clients

## Usage

### Setting Rate Limits

Use the admin client to set rate limits for an API endpoint:

```bash
curl -X POST http://localhost:8080/api/limit \
-H "Content-Type: application/json" \
-d '{
  "api_path": "/api/example",
  "requests_per_second": 10,
  "requests_per_minute": 100,
  "requests_per_hour": 1000,
  "requests_per_day": 10000
}'
```

### Checking Rate Limits

Use the checker client to verify if a request is allowed:

```bash
curl -X POST http://localhost:8080/api/check \
-H "Content-Type: application/json" \
-d '{
  "api_path": "/api/example",
  "client_id": "test-client"
}'
```

Response format:
```json
{
  "allowed": true,
  "reason": ""
}
```

Or when limited:
```json
{
  "allowed": false,
  "reason": "Rate limit exceeded for minute window"
}
```

## Configuration

### Server Configuration

The server can be configured using environment variables:

```bash
REDIS_HOST=localhost    # Redis host (default: localhost)
REDIS_PORT=6379        # Redis port (default: 6379)
SERVER_PORT=8080       # Server HTTP port (default: 8080)
```

### Client Configuration

Both admin and checker clients accept a server URL flag:

```bash
./admin --server=http://localhost:8080
./checker --server=http://localhost:8080
```

## API Endpoints

### Set Rate Limit
- **POST** `/api/limit`
- Sets rate limits for an API endpoint
- Request body:
```json
{
  "api_path": string,
  "requests_per_second": int,
  "requests_per_minute": int,
  "requests_per_hour": int,
  "requests_per_day": int
}
```

### Check Rate Limit
- **POST** `/api/check`
- Checks if a request should be allowed
- Request body:
```json
{
  "api_path": string,
  "client_id": string
}
```

## Implementation Details

### Rate Limiting Algorithm

The system uses a sliding window algorithm implemented with Redis:
1. Each request increments a counter in Redis
2. Counters automatically expire based on their time window
3. Multiple time windows are checked simultaneously
4. First failing check prevents the request

### Redis Storage

Rate limits are stored in Redis using the following key patterns:
- Rate limit configurations: `limit:{api_path}`
- Request counters: `{api_path}:{client_id}:{window}`

### Error Handling

The system provides detailed error responses:
- 400 Bad Request: Invalid input
- 429 Too Many Requests: Rate limit exceeded
- 500 Internal Server Error: Server or Redis issues

## Development

### Running Tests

```bash
go test ./...
```

### Local Development

1. Start Redis:
```bash
docker-compose up redis
```

2. Run the server:
```bash
go run server/main.go
```

3. Run clients:
```bash
go run client/admin/main.go
go run client/checker/main.go
```

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Inspired by Lyft's rate limiter design
- Uses Redis-based distributed rate limiting concepts
- Built with Go's standard HTTP server capabilities
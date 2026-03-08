# PokeServer Go Backend - Clean Architecture

This is a complete rewrite of the PokeServer Go backend following clean architecture principles, with improved error handling, proper dependency injection, and modern Go practices.

## Architecture Overview

The application now follows a clean, layered architecture:

```
┌─────────────────┐
│   HTTP Layer    │  - Handlers, Middleware, Validation
├─────────────────┤
│ Business Layer  │  - Services, Domain Logic
├─────────────────┤
│   Data Layer    │  - Repository, Database
├─────────────────┤
│ External APIs   │  - PokeAPI Client
└─────────────────┘
```

## Key Components

### Configuration (`config.go`)
- Centralized configuration management using Viper
- Environment variable support
- Configuration validation
- Default values for all settings

### Logging (`logger.go`)
- Structured JSON logging using Go's `slog` package
- Context-aware logging
- Request-specific fields
- Error tracking

### Middleware (`middleware.go`)
- **CORS**: Cross-origin resource sharing
- **Logging**: Request/response logging with timing
- **Recovery**: Panic recovery with logging
- Chainable middleware pattern

### Repository Layer (`repository.go`)
- Interface-based design for testability
- Proper error handling with context
- Database connection management
- Structured logging integration

### Service Layer (`service.go`)
- Business logic separation
- Pokemon operations (get, vote, list)
- Error handling and validation
- Database initialization

### Handler Layer (`handlers.go`)
- HTTP request/response handling
- Input validation
- Proper HTTP status codes
- JSON error responses
- Template rendering for web pages

### Validation (`validation.go`)
- Input parameter validation
- Structured error responses
- Type-safe validation functions

### External Client (`pokeclient.go`)
- PokeAPI integration with proper error handling
- No more `log.Fatal` calls
- Structured logging
- HTTP client best practices

## Key Improvements

### 1. Error Handling
- **Before**: `log.Fatal()` everywhere, causing application crashes
- **After**: Proper error propagation with context and structured responses

### 2. Dependency Injection
- **Before**: Global variables and tight coupling
- **After**: Interface-based design with constructor injection

### 3. Configuration Management
- **Before**: Scattered configuration with Viper calls throughout
- **After**: Centralized config struct with validation

### 4. Logging
- **Before**: Basic `log.Print()` statements
- **After**: Structured JSON logging with context and fields

### 5. HTTP Handling
- **Before**: Manual CORS, no validation, poor error responses
- **After**: Middleware-based CORS, validation, proper status codes

### 6. Graceful Shutdown
- **Before**: No signal handling
- **After**: Proper signal handling with graceful shutdown

## API Endpoints

All endpoints now support proper HTTP methods and return structured JSON responses:

### GET /health
Health check endpoint
```json
{"status": "healthy"}
```

### GET /getpokemon
Get a random Pokemon
```json
{
  "id": 25,
  "name": "pikachu",
  "sprites": {
    "front_default": "https://..."
  }
}
```

### GET /getall
Get all Pokemon with vote counts
```json
{
  "pokemon": [...],
  "count": 5
}
```

### GET /vote?id=25&vote=up
Vote for a Pokemon (up/down)
```json
{
  "id": 25,
  "name": "pikachu",
  "vote": 10,
  "url": "https://..."
}
```

## Error Responses

All errors now return structured JSON:
```json
{
  "error": "validation_error",
  "message": "Invalid request parameters",
  "details": [
    {
      "field": "vote",
      "message": "vote is required"
    }
  ]
}
```

## Configuration

The application can be configured via environment variables or a `.env` file:

```yaml
database:
  url: "postgres://user:password@localhost:5432/pokemon?sslmode=disable"
server:
  port: "8080"
pokeapi:
  url: "https://pokeapi.co/api/v2/pokemon?limit="
  max: 1025
```

Environment variables:
- `DATABASE_URL`
- `PORT`
- `POKEAPI_URL`
- `POKEAPI_MAX`

## Testing

The test suite has been updated to work with the new architecture:
- Fixed Pokemon struct inconsistencies
- Added proper dependency injection in tests
- Maintained all existing test functionality

## Running the Application

1. Set up your database URL in `.env` or environment variables
2. Run: `go run .`
3. The server will start with graceful shutdown support
4. Health check available at `/health`

## Migration Notes

The rewrite maintains **full backward compatibility** with the existing API while providing:
- Better error handling
- Improved performance
- Enhanced maintainability
- Proper testing infrastructure
- Production-ready features

All existing endpoints work exactly as before, but now with proper error handling and improved reliability.
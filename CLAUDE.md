# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Nebula-live is a modern backend API service built with Go 1.22+, following Domain-Driven Design (DDD) principles and Clean Architecture. The project uses:

### Core Framework & Libraries
- **Fiber v2.52.9**: High-performance web framework for HTTP APIs
- **EntGo v0.14.1**: Type-safe ORM with code generation
- **Uber Fx v1.24.0**: Dependency injection container for modular architecture
- **Zap v1.28.0**: Structured logging with lumberjack for log rotation
- **Viper v1.20.0**: Configuration management
- **Cobra v1.8.1**: CLI framework

### Database Support
- **PostgreSQL**: Production-ready relational database
- **SQLite**: Development and lightweight deployments via `modernc.org/sqlite`
- **Multi-database**: Switch between databases via configuration

### Development Tools
- **Air**: Hot reload for development
- **Docker**: Containerization with multi-stage builds
- **Docker Compose**: Local development and production orchestration

## Development Commands

### Basic Commands
- **Build**: `go build ./cmd/server`
- **Run**: `go run ./cmd/server`
- **Test**: `go test ./...`
- **Format**: `go fmt ./...`
- **Vet**: `go vet ./...`
- **Mod tidy**: `go mod tidy`

### Development with Hot Reload
- **Install Air**: `go install github.com/cosmtrek/air@latest`
- **Start with hot reload**: `air`

### Docker Commands
- **Build production**: `docker build -t nebula-live .`
- **Build development**: `docker build -f Dockerfile.dev -t nebula-live:dev .`
- **Start dev environment**: `docker-compose -f docker-compose.dev.yml up app-dev`
- **Start production**: `docker-compose up app`
- **Full stack**: `docker-compose --profile postgres --profile redis up`

## Project Structure (DDD Architecture)

```
├── cmd/server/           # Application entry point
├── internal/
│   ├── app/             # Application service layer
│   ├── domain/          # Domain layer (entities, repositories, services)
│   │   ├── entity/      # Domain entities
│   │   ├── repository/  # Repository interfaces
│   │   └── service/     # Domain services
│   └── infrastructure/  # Infrastructure layer
│       ├── config/      # Configuration management
│       ├── logger/      # Logging setup
│       ├── persistence/ # Database implementations
│       └── web/         # HTTP layer (handlers, middleware, routing)
├── pkg/                 # Shared utilities
├── configs/             # Configuration files
└── logs/               # Log files directory
```

## Configuration

Configuration is managed via `configs/config.yaml` and can be overridden with environment variables prefixed with `NEBULA_`.

### Database Configuration Options

#### SQLite (Development & Lightweight)
```yaml
database:
  driver: "sqlite"
  database: "data/nebula_live.db"  # or ":memory:" for in-memory
```

#### PostgreSQL (Production)
```yaml
database:
  driver: "postgres"
  host: "localhost"
  port: 5432
  username: "postgres"
  password: "password"
  database: "nebula_live"
  ssl_mode: "disable"
```

### Configuration Files
- `configs/config.yaml` - Default configuration
- `configs/config-sqlite.yaml` - SQLite example configuration

## Key Design Patterns

- **Modular Architecture**: Fx modules for each layer (infrastructure, persistence, service, handler)
- **Dependency Injection**: Using Fx for clean dependency management with modular providers
- **Domain-Driven Design**: Clear separation between domain, application, and infrastructure layers
- **Clean Architecture**: Dependencies point inward toward the domain
- **Structured Logging**: JSON-formatted logs with proper rotation + global logger for convenience
- **Unified Error Handling**: APIError for consistent error responses across all endpoints

## API Endpoints

### User Management
- `POST /api/v1/users` - Create user
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user
- `GET /api/v1/users` - List users (with pagination: ?page=1&limit=10)

### User Status Management
- `POST /api/v1/users/:id/activate` - Activate user
- `POST /api/v1/users/:id/deactivate` - Deactivate user
- `POST /api/v1/users/:id/ban` - Ban user

### Health Check
- `GET /health` - Application health status
- `GET /api/v1/ping` - API health check

## Database Setup

The application supports both PostgreSQL and SQLite databases:

### Quick Start (SQLite - Recommended for Development)
```bash
# Use SQLite configuration
cp configs/config-sqlite.yaml configs/config.yaml
go run ./cmd/server
```

### Production Setup (PostgreSQL)
1. Ensure PostgreSQL is running
2. Create database: `createdb nebula_live`  
3. Update `configs/config.yaml` with your database settings
4. Default connection: `postgres://postgres:password@localhost:5432/nebula_live?sslmode=disable`

### Database Features
- **Auto-migrations**: Schema migrations run automatically on startup via EntGo
- **Multi-database support**: Switch between PostgreSQL and SQLite via configuration
- **SQLite optimizations**: Foreign keys enabled, WAL mode for better performance

## Logging System

### Global Logger (Recommended)
```go
import "nebula-live/pkg/logger"

logger.Info("Operation successful", zap.String("user_id", "123"))
logger.Error("Operation failed", zap.Error(err))
```

### Dependency Injection Logger (Alternative)
```go
// Constructor injection
func NewService(logger *zap.Logger) Service {
    return &service{logger: logger}
}
```

## Error Handling

All API responses use standardized APIError format:
```go
// Usage in handlers
return c.Status(fiber.StatusBadRequest).JSON(
    errors.NewAPIError(fiber.StatusBadRequest, "Invalid request", "Missing required field")
)

// Response format
{
  "code": 400,
  "error": "Invalid request", 
  "message": "Missing required field"
}
```

## Development Notes

- **EntGo Integration**: Entities are defined in `ent/schema/` and code is auto-generated
- **Modular DI**: Each layer has its own Fx module to avoid parameter explosion
- **Clean Architecture**: Domain layer has no external dependencies
- **Hot Reload**: Use `air` command for automatic restarts during development
- **Docker Support**: Multi-stage builds for production, hot reload for development

## Git Commit Guidelines

This project follows [Conventional Commits](https://www.conventionalcommits.org/) specification for consistent commit messages and automated versioning.

### Commit Message Format
```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Commit Types
- **feat**: New features
- **fix**: Bug fixes
- **docs**: Documentation updates
- **style**: Code formatting (no functional changes)
- **refactor**: Code refactoring
- **perf**: Performance improvements
- **test**: Test-related changes
- **chore**: Build tools, dependency management
- **ci**: CI/CD configuration
- **build**: Build system changes

### Scopes (Optional)
- **api**: API layer changes
- **web**: Web layer (handlers, middleware, routing)
- **domain**: Domain layer (entities, services)
- **infra**: Infrastructure layer
- **config**: Configuration changes
- **db**: Database-related changes
- **docker**: Docker configuration
- **deps**: Dependency updates

### Commit Examples
```bash
feat(api): add user authentication endpoint
fix(db): resolve SQLite connection timeout
docs: update README with Docker instructions
refactor(domain): extract user validation to service
perf(db): add indexes to user queries
chore(deps): update Go dependencies
ci: add GitHub Actions workflow
feat(api)!: change user response format

BREAKING CHANGE: user API now returns different structure
```

### Commit Rules
- Use imperative mood ("add" not "added")
- First letter lowercase
- No period at the end
- Description under 50 characters
- Use English consistently
- Mark breaking changes with `!` or `BREAKING CHANGE:`
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

### Authentication & Security
- **JWT v5.3.0**: JSON Web Token authentication with access and refresh tokens
- **Argon2id**: Secure password hashing algorithm with salt and timing attack protection
- **Authentication Middleware**: Route-level JWT token validation and user context injection

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

### JWT Configuration
```yaml
jwt:
  secret: "your-secret-key-change-this-in-production"
  access_token_ttl: "15m"     # Access token expiration time
  refresh_token_ttl: "168h"   # Refresh token expiration time (7 days)
  issuer: "nebula-live"       # JWT issuer
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
- **JWT Authentication**: Complete JWT token system with access/refresh tokens and middleware protection
- **Password Security**: Argon2id hashing with salt for secure password storage
- **RBAC Authorization**: Role-Based Access Control with fine-grained permission system
- **System Initialization**: Automatic creation of default roles and permissions on startup

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login (returns JWT tokens)
- `GET /api/v1/auth/me` - Get current user information (requires authentication)
- `POST /api/v1/auth/refresh` - Refresh access token using refresh token

### User Management (Requires Admin Role)
⚠️ **All user management endpoints require JWT authentication and admin role**

- `POST /api/v1/users` - Create user
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user
- `GET /api/v1/users` - List users (with pagination: ?page=1&limit=10)

### User Status Management (Requires Admin Role)
- `POST /api/v1/users/:id/activate` - Activate user
- `POST /api/v1/users/:id/deactivate` - Deactivate user
- `POST /api/v1/users/:id/ban` - Ban user

### RBAC Role Management (Requires Admin Role)
- `POST /api/v1/roles` - Create role
- `GET /api/v1/roles/:id` - Get role by ID
- `PUT /api/v1/roles/:id` - Update role
- `DELETE /api/v1/roles/:id` - Delete role
- `GET /api/v1/roles` - List roles (with pagination: ?page=1&limit=10)
- `POST /api/v1/roles/:id/assign` - Assign role to user
- `DELETE /api/v1/roles/:id/users/:userId` - Remove role from user
- `GET /api/v1/roles/users/:userId` - Get user roles

### RBAC Permission Management (Requires Admin Role)
- `POST /api/v1/permissions` - Create permission
- `GET /api/v1/permissions/:id` - Get permission by ID
- `PUT /api/v1/permissions/:id` - Update permission
- `DELETE /api/v1/permissions/:id` - Delete permission
- `GET /api/v1/permissions` - List permissions (with pagination: ?page=1&limit=10)
- `POST /api/v1/permissions/:id/assign` - Assign permission to role
- `DELETE /api/v1/permissions/:id/roles/:roleId` - Remove permission from role
- `GET /api/v1/permissions/roles/:roleId` - Get role permissions
- `GET /api/v1/permissions/users/:userId` - Get user permissions

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

## Authentication System

### JWT Token Management
The application uses a dual-token JWT system:
- **Access Token**: Short-lived (15 minutes) for API access
- **Refresh Token**: Long-lived (7 days) for token renewal

### Password Security
- **Argon2id Algorithm**: Industry-standard password hashing
- **Salt Generation**: Unique salt per password
- **Timing Attack Protection**: Constant-time comparison

### Authentication Middleware
- **RequireAuth**: Mandatory authentication for protected routes
- **OptionalAuth**: Optional authentication for context-aware features
- **User Context**: Authenticated user information injected into request context

### Usage Examples
```go
// Get current user from context
currentUser, exists := auth.GetCurrentUser(c)
if !exists {
    return c.Status(fiber.StatusUnauthorized).JSON(...)
}

// Hash password securely
hashedPassword, err := security.HashPassword("password123")

// Verify password
isValid, err := security.VerifyPassword("password123", hashedPassword)

// Generate JWT tokens
tokenPair, err := jwtManager.GenerateTokenPair(userID, username, email)
```

## RBAC Authorization System

### Role-Based Access Control Overview
The application implements a comprehensive RBAC system with:
- **Roles**: Named collections of permissions (e.g., `admin`, `user`)
- **Permissions**: Specific actions on resources (e.g., `user:read`, `user:write`)
- **User-Role Assignment**: Users can have multiple roles
- **Role-Permission Assignment**: Roles can have multiple permissions

### System Roles and Permissions
**Default Roles:**
- `admin` - Administrator with all system permissions
- `user` - Basic user with read-only permissions

**System Permissions:**
- User management: `user:read`, `user:write`, `user:delete`, `user:manage`
- Role management: `role:read`, `role:write`, `role:delete`, `role:manage`
- Permission management: `permission:read`, `permission:write`, `permission:delete`, `permission:manage`
- System management: `system:manage`

### RBAC Middleware Usage
```go
// Require specific permission
router.Group("/api/v1/content").Use(
    rbacMiddleware.RequirePermission("content", "read"),
)

// Require specific role
router.Group("/api/v1/admin").Use(
    rbacMiddleware.RequireRole("admin"),
)

// Require admin role (shorthand)
router.Group("/api/v1/users").Use(
    rbacMiddleware.RequireAdmin(),
)
```

### RBAC Service Usage
```go
// Check user permissions
hasPermission, err := rbacService.HasPermission(ctx, userID, "user", "write")

// Check user roles
hasRole, err := rbacService.HasRole(ctx, userID, "admin")

// Assign role to user
err := userService.AssignRole(ctx, userID, "admin", assignerID)

// Get user permissions
permissions, err := rbacService.GetUserPermissions(ctx, userID)
```

### System Initialization
- **Automatic Setup**: System roles and permissions are created automatically on first startup
- **Idempotent**: Safe to run multiple times, existing data is preserved
- **Configurable**: System permissions can be modified through the API

## Development Notes

- **EntGo Integration**: Entities are defined in `ent/schema/` and code is auto-generated
- **Modular DI**: Each layer has its own Fx module to avoid parameter explosion
- **Clean Architecture**: Domain layer has no external dependencies
- **Hot Reload**: Use `air` command for automatic restarts during development
- **Docker Support**: Multi-stage builds for production, hot reload for development
- **JWT Security**: All user management endpoints require valid JWT authentication
- **Route Protection**: Authentication middleware automatically validates tokens and injects user context
- **RBAC Integration**: User management requires admin role, fine-grained permissions available
- **System Bootstrap**: Default roles and permissions created automatically on first run

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
- **auth**: Authentication and authorization changes
- **security**: Security-related improvements

### Commit Examples
```bash
feat(api): add user authentication endpoint
feat(auth): implement JWT token system with refresh tokens
feat(security): add Argon2id password hashing
feat(middleware): add JWT authentication middleware
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
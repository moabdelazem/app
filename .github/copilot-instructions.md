# GitHub Copilot Instructions

## Project Overview

This is a **Guest Book API** built with Go, featuring a clean architecture with PostgreSQL database integration. The application provides REST API endpoints for managing guest book messages with full CRUD operations, proper error handling, and comprehensive logging.

## Architecture & Design Patterns

### Clean Architecture
The project follows clean architecture principles with clear separation of concerns:

```
cmd/                    # Application entry points
├── main.go            # Main application with graceful shutdown

internal/              # Private application code
├── config/            # Configuration management with environment variables
├── database/          # Database connection and pooling (pgx driver)
├── handlers/          # HTTP handlers with proper JSON responses
├── logger/            # Structured logging with slog
├── models/            # Domain models and DTOs
├── repository/        # Data access layer with PostgreSQL integration
├── server/            # HTTP server setup with middleware
└── service/           # Business logic layer

pkg/                   # Public library code (future use)
```

### Key Design Principles
- **Dependency Injection**: Dependencies are injected through constructors
- **Interface Segregation**: Each layer defines its own interfaces
- **Single Responsibility**: Each package has a specific purpose
- **Error Handling**: Comprehensive error handling with structured logging
- **Configuration**: Environment-based configuration with defaults

## Technology Stack

### Core Technologies
- **Go 1.24.2**: Programming language
- **PostgreSQL**: Primary database with pgx driver
- **Gorilla Mux**: HTTP router and middleware
- **pgx/v5**: Modern PostgreSQL driver with connection pooling
- **slog**: Structured logging (Go standard library)
- **Docker Compose**: Development database setup

### Key Dependencies
```go
github.com/gorilla/mux v1.8.1      // HTTP routing
github.com/joho/godotenv v1.5.1    // Environment file loading
github.com/jackc/pgx/v5 v5.7.5     // PostgreSQL driver
```

## Database Schema

### Guest Book Messages Table
```sql
CREATE TABLE guest_book_messages (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_guest_book_created_at ON guest_book_messages(created_at DESC);
```

## API Endpoints

### Core Endpoints
| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| `GET` | `/` | API information and documentation | None |
| `GET` | `/health` | Basic health check | None |
| `GET` | `/api/v1/health` | Health check with database connectivity | None |
| `GET` | `/api/v1/guestbook` | Get all messages (supports pagination) | None |
| `POST` | `/api/v1/guestbook` | Create a new message | CreateGuestBookMessage |
| `GET` | `/api/v1/guestbook/{id}` | Get specific message by ID | None |

### Request/Response Examples

#### Create Message (POST /api/v1/guestbook)
```json
// Request
{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "message": "Hello! This is my message in the guest book."
}

// Response (201 Created)
{
  "id": 1,
  "name": "John Doe",
  "email": "john.doe@example.com",
  "message": "Hello! This is my message in the guest book.",
  "created_at": "2025-07-15T10:30:00Z",
  "updated_at": "2025-07-15T10:30:00Z"
}
```

#### Get Messages with Pagination (GET /api/v1/guestbook?page=1&page_size=10)
```json
{
  "messages": [...],
  "pagination": {
    "page": 1,
    "page_size": 10,
    "total": 25,
    "total_pages": 3
  }
}
```

### Error Responses
All errors follow a consistent JSON format:
```json
{
  "error": "Error Type",
  "message": "Detailed error description",
  "path": "/api/v1/guestbook",
  "method": "POST"
}
```

## Configuration

### Environment Variables
```bash
# Application Configuration
PORT=4260
DEBUG=false

# Database Configuration
DB_HOST=localhost
DB_USER=root
DB_PASSWORD=example
DB_NAME=example_db
DB_PORT=5432
DB_SSL_MODE=disable
```

### Configuration Loading
- Uses `.env` file for development
- Falls back to system environment variables
- Provides sensible defaults for all values
- Validates required database connection parameters

## Development Workflow

### Setup Commands
```bash
# Start database
docker-compose -f docker-compose.dev.yml up -d

# Install dependencies
go mod download

# Run application
make run
# or
go run cmd/main.go

# Build application
make build

# Run tests
make test

# Format code
make fmt
```

### Database Development
- Automatic table creation on startup
- Connection pooling with configurable parameters
- Health checks for database connectivity
- Graceful shutdown with connection cleanup

## Code Style & Conventions

### File Organization
- **cmd/**: Application entry points only
- **internal/**: Private packages, not importable by external projects
- **pkg/**: Public packages (future use)
- Each package should have a single responsibility

### Naming Conventions
- **Files**: lowercase with underscores (e.g., `guestbook.go`)
- **Packages**: lowercase, single word when possible
- **Types**: PascalCase (e.g., `GuestBookMessage`)
- **Functions/Methods**: PascalCase for exported, camelCase for private
- **Constants**: PascalCase or UPPER_CASE for package-level

### Error Handling
- Always handle errors explicitly
- Use structured logging for error context
- Return appropriate HTTP status codes
- Provide meaningful error messages to clients

### Testing Guidelines
- Write unit tests for business logic
- Use table-driven tests where applicable
- Mock external dependencies
- Test error conditions

## Middleware & Cross-Cutting Concerns

### Implemented Middleware
1. **Logging Middleware**: Logs all HTTP requests with timing
2. **CORS Middleware**: Handles cross-origin requests
3. **Error Handling**: Consistent error response format

### Security Considerations
- Input validation on all endpoints
- SQL injection prevention through parameterized queries
- CORS configuration for controlled access
- Environment-based configuration (no hardcoded secrets)

## Deployment & Operations

### Docker Support
- PostgreSQL container for development
- Environment-based configuration
- Volume persistence for database data
- Health checks for service monitoring

### Monitoring & Observability
- Structured logging with slog
- HTTP request/response logging
- Database connection health checks
- Graceful shutdown with proper cleanup

## Future Development Guidelines

### Adding New Features
1. **New Endpoints**: Add to handlers package, register in server
2. **Database Changes**: Update repository layer, add migrations
3. **Business Logic**: Implement in service layer
4. **Models**: Define in models package with proper validation

### Scaling Considerations
- Database connection pooling is already configured
- Stateless design allows horizontal scaling
- Structured logging supports centralized log aggregation
- Configuration supports multiple environments

### Code Quality Standards
- Follow Go best practices and idioms
- Maintain test coverage above 80%
- Use go fmt and go vet for code quality
- Document public APIs with proper comments

## Common Patterns & Examples

### Adding a New Endpoint
1. Define model in `internal/models/`
2. Add repository methods in `internal/repository/`
3. Implement business logic in `internal/service/`
4. Create handler in `internal/handlers/`
5. Register route in `internal/server/server.go`

### Database Operations
Always use context for database operations:
```go
func (r *Repository) GetByID(ctx context.Context, id int) (*Model, error) {
    query := `SELECT ... FROM table WHERE id = $1`
    var result Model
    err := r.db.Pool.QueryRow(ctx, query, id).Scan(...)
    if err != nil {
        return nil, fmt.Errorf("failed to get record: %w", err)
    }
    return &result, nil
}
```

### Error Handling Pattern
```go
if err != nil {
    slog.Error("Operation failed", "error", err, "context", value)
    RespondJSON(w, http.StatusInternalServerError, map[string]string{
        "error": "Operation failed",
    })
    return
}
```

## Dependencies & Libraries

### Core Dependencies
- **gorilla/mux**: Mature HTTP router with good middleware support
- **jackc/pgx**: High-performance PostgreSQL driver
- **joho/godotenv**: Environment file loading for development

### Why These Choices
- **pgx over lib/pq**: Better performance, modern API, connection pooling
- **Gorilla Mux**: Simple, reliable, well-documented routing
- **Standard library**: Prefer stdlib (slog, context, http) when possible

## Troubleshooting

### Common Issues
1. **Database Connection**: Check Docker container status and environment variables
2. **Port Conflicts**: Ensure PORT environment variable is set correctly
3. **Migration Issues**: Tables are created automatically on startup
4. **CORS Issues**: Middleware is configured to allow all origins in development

### Debug Mode
Set `DEBUG=true` in environment for verbose logging and detailed error messages.

---

This project demonstrates modern Go development practices with clean architecture, proper error handling, and production-ready patterns. Use this as a reference for maintaining consistency and quality in future development.

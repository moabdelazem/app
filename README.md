# Go API Project

A well-structured Go REST API using Gorilla Mux with proper separation of concerns.

## Project Structure

```
.
├── cmd/                    # Application entry points
│   └── main.go            # Main application
├── internal/              # Private application code
│   ├── config/            # Configuration management
│   ├── handlers/          # HTTP handlers
│   ├── logger/            # Logging setup
│   └── server/            # HTTP server setup
├── pkg/                   # Public library code (currently empty)
├── go.mod                 # Go module file
├── go.sum                 # Go dependencies
├── Makefile              # Build automation
└── README.md             # This file
```

## Features

- Clean architecture with separation of concerns
- Structured logging with slog
- Environment-based configuration
- Graceful shutdown
- JSON response handling with proper error management
- Request logging middleware
- Health check endpoint

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Make (optional, for using Makefile commands)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/moabdelazem/app
cd app
```

2. Install dependencies:
```bash
go mod download
```

### Running the Application

1. Using Go directly:
```bash
go run cmd/main.go
```

2. Using Make (if Makefile is configured):
```bash
make run
```

### Configuration

The application uses environment variables for configuration. You can set these variables in a `.env` file or as system environment variables.

#### Setup Environment Variables

1. Copy the example environment file:
```bash
cp .env.example .env
```

2. Edit `.env` with your preferred settings:
```env
PORT=4260
DEBUG=false
```

#### Available Configuration Options

- `PORT`: Server port (default: 4260)
- `DEBUG`: Enable debug logging (default: false)

#### Environment Variable Priority

1. System environment variables (highest priority)
2. `.env` file variables
3. Default values (lowest priority)

Example with environment variables:
```bash
PORT=8080 DEBUG=true go run cmd/main.go
```

## API Endpoints

### Base Endpoints

- `GET /` - API version information
- `GET /health` - Health check

### API v1 Endpoints

- `GET /api/v1/health` - Health check (API versioned)

## Development

### Project Layout

This project follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout):

- **cmd/**: Main applications for this project
- **internal/**: Private application and library code
- **pkg/**: Library code that's ok to use by external applications

### Adding New Features

1. **Handlers**: Add new HTTP handlers in `internal/handlers/`
2. **Routes**: Register new routes in `internal/server/server.go`
3. **Configuration**: Add configuration options in `internal/config/`
4. **Business Logic**: Add business logic in appropriate packages under `internal/`

### Building

```bash
go build -o bin/api cmd/main.go
```

### Testing

```bash
go test ./...
```

## Contributing

1. Follow Go best practices and conventions
2. Add tests for new features
3. Update documentation as needed
4. Use structured logging with appropriate levels

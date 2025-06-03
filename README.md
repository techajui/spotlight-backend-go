# Spotlight Backend Go

A Go-based backend service for the Spotlight application.

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- PostgreSQL (if running locally)

## Getting Started

1. Clone the repository:
```bash
git clone <repository-url>
cd spotlight-backend-go
```

2. Install dependencies:
```bash
go mod download
```

3. Run the application locally:
```bash
# Using Docker Compose (recommended)
docker-compose up --build

# Or run directly with Go
go run cmd/api/main.go
```

The application will be available at `http://localhost:8080`

## API Endpoints

- `GET /health` - Health check endpoint

## Development

The project structure is organized as follows:

```
spotlight-backend-go/
├── cmd/
│   └── api/          # Main application entry point
├── internal/
│   ├── config/       # Configuration management
│   ├── database/     # Database related code
│   ├── handlers/     # HTTP handlers
│   ├── middleware/   # HTTP middleware
│   └── models/       # Data models
└── pkg/              # Public packages
```

## Environment Variables

The following environment variables can be configured:

- `DB_HOST` - Database host (default: localhost)
- `DB_USER` - Database user (default: postgres)
- `DB_PASSWORD` - Database password (default: postgres)
- `DB_NAME` - Database name (default: spotlight)
- `DB_PORT` - Database port (default: 5432) 
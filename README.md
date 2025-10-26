# GoPilot

[![CI/CD Pipeline](https://github.com/codewithwan/gopilot/actions/workflows/ci.yml/badge.svg)](https://github.com/codewithwan/gopilot/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/codewithwan/gopilot)](https://goreportcard.com/report/github.com/codewithwan/gopilot)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Production-ready REST API service for managing todos, built with Go and following clean architecture principles.

## Features

- 🚀 **RESTful API** with CRUD operations for todos
- 🔐 **JWT Authentication** for secure access
- 🗄️ **PostgreSQL** database with sqlc for type-safe queries
- 📝 **Structured Logging** with zap
- 📊 **Prometheus Metrics** for monitoring
- 🔍 **OpenTelemetry** tracing support
- 📚 **Swagger Documentation** (OpenAPI)
- 🔧 **Viper Configuration** management
- 🐳 **Docker** support with multi-stage build
- 🔄 **Database Migrations** with golang-migrate
- ✅ **Automated Testing** and linting
- 🚦 **CI/CD** with GitHub Actions
- 📦 **Container Registry** (GitHub Container Registry)

## Tech Stack

- **Framework**: Gin
- **Database**: PostgreSQL
- **ORM/Query Builder**: sqlc
- **Configuration**: Viper
- **Logging**: Zap
- **Metrics**: Prometheus
- **Tracing**: OpenTelemetry
- **Documentation**: Swagger/OpenAPI
- **Authentication**: JWT
- **Migrations**: golang-migrate
- **Containerization**: Docker
- **CI/CD**: GitHub Actions

## Project Structure

```
gopilot/
├── cmd/
│   └── server/          # Main application entry point
├── internal/
│   ├── config/          # Configuration management
│   ├── domain/          # Domain models and DTOs
│   ├── handler/         # HTTP handlers
│   ├── middleware/      # Custom middlewares (JWT, etc.)
│   ├── repository/      # Database repository layer
│   │   └── db/          # Generated sqlc code
│   └── service/         # Business logic layer
├── pkg/
│   ├── logger/          # Logger utilities
│   ├── metrics/         # Prometheus metrics
│   └── tracing/         # OpenTelemetry tracing
├── db/
│   ├── migrations/      # Database migration files
│   └── queries/         # SQL queries for sqlc
├── docs/                # Swagger documentation
├── Dockerfile           # Multi-stage Docker build
├── Makefile             # Build automation
├── docker-compose.yml   # Local development setup
└── .github/
    └── workflows/       # CI/CD pipelines
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15+
- Docker (optional)
- Make (optional)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/codewithwan/gopilot.git
cd gopilot
```

2. Install dependencies:
```bash
go mod download
```

3. Install development tools:
```bash
make install-tools
```

4. Set up configuration:
```bash
cp config.yaml.example config.yaml
# Edit config.yaml with your settings
```

5. Set up the database:
```bash
# Start PostgreSQL (using Docker)
docker run -d \
  --name gopilot-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=gopilot \
  -p 5432:5432 \
  postgres:15-alpine

# Run migrations
make migrate-up
```

6. Run the application:
```bash
make run
```

The API will be available at `http://localhost:8080`

### Using Docker Compose

For local development with all dependencies:

```bash
# Start all services
make dev

# Stop all services
make dev-down
```

## API Documentation

Once the application is running, access the Swagger UI at:
```
http://localhost:8080/swagger/index.html
```

### API Endpoints

#### Authentication
- `POST /api/v1/auth/register` - Register a new user
- `POST /api/v1/auth/login` - Login and get JWT token

#### Todos (Protected)
- `GET /api/v1/todos` - List all todos
- `POST /api/v1/todos` - Create a new todo
- `GET /api/v1/todos/:id` - Get a specific todo
- `PUT /api/v1/todos/:id` - Update a todo
- `DELETE /api/v1/todos/:id` - Delete a todo

#### Health & Metrics
- `GET /health` - Health check endpoint
- `GET /metrics` - Prometheus metrics

### Example Usage

1. Register a new user:
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","password":"password123"}'
```

2. Login:
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user1","password":"password123"}'
```

3. Create a todo (with JWT token):
```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{"title":"My Todo","description":"Todo description"}'
```

4. List todos:
```bash
curl -X GET http://localhost:8080/api/v1/todos \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Development

### Building

```bash
# Build the binary
make build

# Build Docker image
make docker-build
```

### Testing

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage
```

### Linting

```bash
make lint
```

### Database Operations

```bash
# Generate sqlc code
make sqlc-generate

# Create a new migration
make migrate-create NAME=add_new_table

# Run migrations
make migrate-up

# Rollback migrations
make migrate-down

# Reset database
make db-reset
```

### Swagger Documentation

```bash
# Generate Swagger docs
make swagger-generate
```

## Configuration

Configuration can be provided via:
1. YAML config file (`config.yaml`)
2. Environment variables (uppercase with underscores)

Example environment variables:
```bash
SERVER_PORT=8080
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres
DATABASE_DBNAME=gopilot
JWT_SECRET=your-secret-key
LOG_LEVEL=info
```

## CI/CD

The project includes a comprehensive GitHub Actions workflow that:
- Runs linting with golangci-lint
- Executes tests with coverage
- Builds the application
- Builds and pushes Docker images to GitHub Container Registry (GHCR)

The workflow runs on:
- Every push to `main` branch
- Every pull request to `main` branch

Docker images are pushed to GHCR only on pushes to the `main` branch.

## Monitoring

### Prometheus Metrics

The application exposes Prometheus metrics at `/metrics`. Key metrics include:
- HTTP request count
- HTTP request duration
- Go runtime metrics

### OpenTelemetry Tracing

To enable tracing, configure the tracing endpoint in your config:
```yaml
tracing:
  enabled: true
  serviceName: "gopilot"
  endpoint: "otel-collector:4317"
```

## Production Deployment

### Using Docker

Pull and run the latest image:
```bash
docker pull ghcr.io/codewithwan/gopilot:latest
docker run -p 8080:8080 \
  -e DATABASE_HOST=your-db-host \
  -e DATABASE_USER=your-db-user \
  -e DATABASE_PASSWORD=your-db-password \
  ghcr.io/codewithwan/gopilot:latest
```

### Environment Variables

Required environment variables for production:
- `DATABASE_HOST`
- `DATABASE_PORT`
- `DATABASE_USER`
- `DATABASE_PASSWORD`
- `DATABASE_DBNAME`
- `JWT_SECRET` (use a strong secret!)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [sqlc](https://github.com/sqlc-dev/sqlc)
- [Viper](https://github.com/spf13/viper)
- [Zap Logger](https://github.com/uber-go/zap)
- [OpenTelemetry](https://opentelemetry.io/)
- [Prometheus](https://prometheus.io/)
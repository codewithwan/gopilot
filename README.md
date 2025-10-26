# GoPilot - Developer Tools API Platform

[![CI/CD Pipeline](https://github.com/codewithwan/gopilot/actions/workflows/ci.yml/badge.svg)](https://github.com/codewithwan/gopilot/actions/workflows/ci.yml)
[![Security Scanning](https://github.com/codewithwan/gopilot/actions/workflows/security.yml/badge.svg)](https://github.com/codewithwan/gopilot/actions/workflows/security.yml)
[![CodeQL](https://github.com/codewithwan/gopilot/actions/workflows/codeql.yml/badge.svg)](https://github.com/codewithwan/gopilot/actions/workflows/codeql.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/codewithwan/gopilot)](https://goreportcard.com/report/github.com/codewithwan/gopilot)
[![codecov](https://codecov.io/gh/codewithwan/gopilot/branch/main/graph/badge.svg)](https://codecov.io/gh/codewithwan/gopilot)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

🧠 **GoPilot** is a production-ready REST API platform that provides a comprehensive collection of modular developer tools. All endpoints are self-contained, require no external dependencies, and can be accessed directly from Swagger UI without writing code.

## 🎯 Vision

Provide a collection of modular REST APIs for various developer needs: URL shortening, encoding, data tools, file sharing, automation, and analytics — accessible directly from Swagger UI without writing code.

## Features

- 🔗 **URL Shortener** - Create short links with custom aliases and expiration
- 📝 **Pastebin/Snippet Storage** - Share code snippets with syntax highlighting
- 🔲 **QR Code Generator** - Generate QR codes for URLs, text, and more
- 🔐 **Hash & Encode** - MD5, SHA, bcrypt, base64, hex encoding
- 🔄 **Data Converter** - Convert between bases, colors, time formats, JSON/YAML
- 🆔 **UUID & Token Generator** - Generate secure UUIDs and tokens
- 📊 **Mock Data Generator** - Lorem ipsum, fake users, random numbers
- 🎨 **JSON/YAML Formatter** - Format and convert structured data
- 🔒 **Crypto Playground** - AES, RSA, HMAC operations
- 🚀 **RESTful API** with comprehensive documentation
- 🔐 **JWT Authentication** for protected resources
- 🗄️ **PostgreSQL** database with type-safe queries
- 📝 **Structured Logging** with Zap
- 📊 **Prometheus Metrics** for monitoring
- 🔍 **OpenTelemetry** tracing support
- 📚 **Interactive Swagger UI** (OpenAPI 3)
- 🔧 **Viper Configuration** management
- 🐳 **Docker** support with multi-stage build
- ✅ **Automated Testing** and linting
- 🚦 **CI/CD** with GitHub Actions

## Tech Stack

- **Language**: Go 1.22+
- **Framework**: Gin (REST router)
- **Database**: PostgreSQL 15+
- **Query Builder**: sqlc (type-safe SQL)
- **Configuration**: Viper (env + yaml)
- **Logging**: Zap (structured logs)
- **Metrics**: Prometheus
- **Tracing**: OpenTelemetry
- **Documentation**: Swagger/OpenAPI 3
- **Authentication**: JWT
- **Migrations**: golang-migrate
- **Containerization**: Docker
- **CI/CD**: GitHub Actions → GHCR
- **License**: MIT

## 🧩 Core API Modules

### 1️⃣ URL Shortener
Create and manage shortened URLs with analytics.

**Endpoints:**
- `POST /v1/shorten` - Create short link (original_url, optional alias, expire_in)
- `GET /s/:code` - Redirect to original URL
- `GET /v1/shorten/:code` - Get statistics

**Features:**
- Base62 ID generator
- Custom aliases
- Expiration support
- Click tracking (referrer, user agent, IP)
- Auto-cleanup of expired links

### 2️⃣ Pastebin / Snippet Storage
Store and share code snippets.

**Endpoints:**
- `POST /v1/paste` - Create paste
- `GET /p/:id` - View paste
- `DELETE /v1/paste/:id` - Delete paste
- `GET /v1/paste/recent` - List recent pastes

**Features:**
- TTL per paste (default 24h)
- Syntax highlighting support
- Public/private mode
- Compression option

### 3️⃣ QR Code Generator
Generate QR codes from text or URLs.

**Endpoints:**
- `POST /v1/qr` - Generate QR code
- `GET /v1/qr/:id` - Get QR code image

**Features:**
- PNG output
- Configurable size
- Persistent storage

### 4️⃣ Hash & Encode Utilities
Hash and encode text data.

**Endpoints:**
- `POST /v1/hash` - Hash text (md5, sha1, sha256, sha512, bcrypt)
- `POST /v1/encode` - Encode/decode (base64, url, hex)
- `POST /v1/generate/password` - Generate secure passwords

**Features:**
- Multiple hash algorithms
- Configurable password generation
- Salt support for hashing

### 5️⃣ Base & Data Converter
Convert data between different formats.

**Endpoints:**
- `POST /v1/convert/base` - Convert number bases (2-64)
- `POST /v1/convert/color` - RGB ↔ HEX conversion
- `POST /v1/convert/time` - Unix ↔ ISO8601 ↔ human readable

**Features:**
- Automatic type detection
- Validation and error handling

### 6️⃣ UUID & Token Generator
Generate UUIDs and random tokens.

**Endpoints:**
- `POST /v1/generate/uuid` - Generate UUID v1/v4/v7
- `POST /v1/generate/token` - Generate random tokens

**Features:**
- Multiple UUID versions
- Custom prefix/suffix support
- Secure random generation

### 7️⃣ Lorem Ipsum & Mock Data
Generate fake data for testing.

**Endpoints:**
- `POST /v1/generate/lorem` - Lorem ipsum text
- `POST /v1/generate/user` - Fake user profiles
- `POST /v1/generate/number` - Random numbers

**Features:**
- Configurable count
- Multiple data formats
- Realistic fake data

### 8️⃣ JSON / YAML Formatter
Format and convert structured data.

**Endpoints:**
- `POST /v1/format/json` - Format/minify JSON
- `POST /v1/format/yaml` - Convert JSON ↔ YAML

**Features:**
- Validation checks
- Configurable indentation
- Syntax validation

### 9️⃣ Crypto Playground
Cryptographic operations for development.

**Endpoints:**
- `POST /v1/crypto/aes` - AES encrypt/decrypt
- `POST /v1/crypto/rsa/keygen` - Generate RSA keypair
- `POST /v1/crypto/rsa` - RSA encrypt/decrypt
- `POST /v1/crypto/hmac` - HMAC sign/verify

**Features:**
- Secure encryption
- Key generation
- Signature verification
- Ephemeral key storage (memory only)

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

Once the application is running, access the interactive Swagger UI at:
```
http://localhost:8080/swagger/index.html
```
or simply:
```
http://localhost:8080/docs
```

### Quick Start API Examples

#### URL Shortener
```bash
# Create a short URL
curl -X POST http://localhost:8080/v1/shorten \
  -H "Content-Type: application/json" \
  -d '{"original_url":"https://example.com","expire_in":24}'

# Access the short URL (redirects)
curl http://localhost:8080/s/abc123
```

#### Pastebin
```bash
# Create a paste
curl -X POST http://localhost:8080/v1/paste \
  -H "Content-Type: application/json" \
  -d '{"content":"console.log(\"Hello\");","syntax":"javascript"}'

# View paste
curl http://localhost:8080/p/paste_id
```

#### QR Code
```bash
# Generate QR code
curl -X POST http://localhost:8080/v1/qr \
  -H "Content-Type: application/json" \
  -d '{"text":"https://example.com","size":256}'

# Download QR image
curl http://localhost:8080/v1/qr/qr_id -o qrcode.png
```

#### Hash & Encode
```bash
# Hash text
curl -X POST http://localhost:8080/v1/hash \
  -H "Content-Type: application/json" \
  -d '{"text":"password123","algorithm":"sha256"}'

# Base64 encode
curl -X POST http://localhost:8080/v1/encode \
  -H "Content-Type: application/json" \
  -d '{"text":"Hello World","operation":"base64-encode"}'
```

#### Generators
```bash
# Generate UUID
curl -X POST http://localhost:8080/v1/generate/uuid \
  -H "Content-Type: application/json" \
  -d '{"version":4,"count":5}'

# Generate password
curl -X POST http://localhost:8080/v1/generate/password \
  -H "Content-Type: application/json" \
  -d '{"length":16,"include_symbols":true}'

# Generate fake user
curl -X POST http://localhost:8080/v1/generate/user \
  -H "Content-Type: application/json" \
  -d '{"count":3}'
```

#### Converters
```bash
# Convert color
curl -X POST http://localhost:8080/v1/convert/color \
  -H "Content-Type: application/json" \
  -d '{"value":"#FF5733","to":"rgb"}'

# Format JSON
curl -X POST http://localhost:8080/v1/format/json \
  -H "Content-Type: application/json" \
  -d '{"json":"{\"key\":\"value\"}","minify":false,"indent":2}'
```

#### Crypto
```bash
# AES encrypt (use strong 32-character keys for AES-256)
curl -X POST http://localhost:8080/v1/crypto/aes \
  -H "Content-Type: application/json" \
  -d '{"operation":"encrypt","text":"secret","key":"strongkey32characterslongxxxx"}'
```

### Legacy Endpoints (Todo App)

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
- `GET /health` - Basic health check
- `GET /healthz` - Kubernetes-style health check
- `GET /readyz` - Readiness check (includes DB connectivity)
- `GET /metrics` - Prometheus metrics

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

# Seed database with sample data (placeholder for future implementation)
make seed
```

### Swagger Documentation

```bash
# Generate Swagger docs
make swagger-generate

# Or manually
swag init -g cmd/server/main.go -o docs
```

### Code Generation

```bash
# Generate all code (sqlc + swagger)
make gen:openapi
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
- HTTP request count and duration
- Database connection pool metrics
- Go runtime metrics

### OpenTelemetry Tracing

To enable distributed tracing, configure the tracing endpoint in your config:
```yaml
tracing:
  enabled: true
  serviceName: "gopilot"
  endpoint: "otel-collector:4317"
```

### Structured Logging

All logs are output in structured JSON format with fields:
- Level (debug, info, warn, error)
- Timestamp
- Message
- Context fields (request_id, user_id, etc.)

Example log output:
```json
{
  "level": "info",
  "ts": 1234567890,
  "msg": "Short URL created",
  "code": "abc123",
  "original_url": "https://example.com"
}
```

## 🖥️ Frontend (UI)

GoPilot uses **Swagger UI** as its primary frontend interface, allowing all APIs to be accessed directly without coding.

### Features:
- ✅ Interactive API documentation at `/docs`
- ✅ Request/response examples for all endpoints
- ✅ Try-it-out functionality for testing
- ✅ Built-in authentication support (JWT)
- ✅ Auto-generated from OpenAPI spec
- ✅ No frontend development needed

## Architecture

GoPilot follows **Clean Architecture** principles:

```
┌─────────────────────────────────────┐
│         HTTP Handlers               │  (Gin Controllers)
├─────────────────────────────────────┤
│         Service Layer               │  (Business Logic)
├─────────────────────────────────────┤
│        Repository Layer             │  (Data Access)
├─────────────────────────────────────┤
│    Database (PostgreSQL + sqlc)     │  (Persistence)
└─────────────────────────────────────┘
```

**Key Principles:**
- Separation of concerns
- Dependency injection
- Interface-based design
- Testability
- Modularity

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
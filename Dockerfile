# Build stage
FROM golang:1.24.7-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gopilot ./cmd/server

# Final stage
FROM alpine:latest

# Install ca-certificates, postgresql-client for migrations
RUN apk --no-cache add ca-certificates postgresql-client

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/gopilot .

# Copy migration files
COPY --from=builder /app/db/migrations ./db/migrations

# Copy entrypoint script
COPY docker-entrypoint.sh .
RUN chmod +x docker-entrypoint.sh

# Copy config example (optional)
COPY --from=builder /app/config.yaml.example ./config.yaml.example

# Expose port
EXPOSE 8080

# Run the entrypoint script
CMD ["./docker-entrypoint.sh"]

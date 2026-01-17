# Build Stage
FROM golang:1.24-alpine AS builder

# Install build dependencies (needed for CGO/SQLite)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
# CGO_ENABLED=1 is required for go-sqlite3
RUN CGO_ENABLED=1 GOOS=linux go build -o ledgerly cmd/main.go

# Run Stage
FROM alpine:latest

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/ledgerly .

# Copy .env.example (optional, user should mount .env or provide env vars)
# We won't copy .env to avoid leaking secrets in the image
COPY .env.example .

# Expose port
EXPOSE 8080

# Command to run
CMD ["./ledgerly"]

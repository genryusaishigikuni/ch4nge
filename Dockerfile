# Build stage
FROM golang:1.23-alpine AS builder

# Install git and build dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files first (for better caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with verbose output
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -a -installsuffix cgo -o main .

# Verify the binary was created
RUN ls -la main

# Final stage
FROM alpine:latest

# Install ca-certificates and curl for health checks
RUN apk --no-cache add ca-certificates tzdata curl

# Create non-root user
RUN addgroup -g 1001 -S appuser && \
    adduser -S -D -H -u 1001 -s /sbin/nologin -G appuser appuser

# Create app directory and uploads directory with proper permissions
RUN mkdir -p /app/uploads/profiles && \
    chown -R appuser:appuser /app && \
    chmod -R 755 /app

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Ensure the binary has correct permissions and verify it exists
RUN ls -la main && \
    chown appuser:appuser main && \
    chmod +x main

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Health check (updated to use curl)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

# Run the application
CMD ["./main"]
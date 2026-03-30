# Build stage
FROM golang:1.22-alpine AS builder

# Install git and ca-certificates (needed for go modules)
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o tfd .

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S tfd && \
    adduser -u 1001 -S tfd -G tfd

# Set up directories
RUN mkdir -p /var/lib/tfd/Pictures \
    /var/lib/tfd/Videos \
    /var/lib/tfd/Music \
    /var/lib/tfd/Voices \
    /var/lib/tfd/Documents \
    /var/lib/tfd && \
    chown -R tfd:tfd /var/lib/tfd

# Create app directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/tfd .

# Copy configuration example
COPY --from=builder /app/config.example.yaml .

# Change to non-root user
USER tfd

# Set environment variables
ENV TFD_CONFIG_PATH=/app/config.yaml

# Health check (optional - checks if process is running)
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD pgrep -f "tfd" || exit 1

# Run the application  
ENTRYPOINT ["./tfd"]

# Default command (can be overridden)
CMD ["--config", "/app/config.yaml"]
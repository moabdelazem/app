# Stage 1: Build the application
FROM golang:1.24-alpine AS builder

# Install build dependencies and CA certificates
RUN apk add --no-cache git ca-certificates tzdata && \
    adduser -D -s /bin/sh -u 1001 appuser

# Set the working directory
WORKDIR /app

# Copy dependency files first for better caching
COPY go.mod go.sum ./

# Download and verify dependencies
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Run tests during build (optional - can be disabled for faster builds)
RUN go test ./... -short

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o main ./cmd/main.go

# Stage 2: Create the runtime image
FROM scratch

# Copy essential files from builder
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group

# Set the working directory
WORKDIR /app

# Copy the compiled binary with proper ownership
COPY --from=builder --chown=1001:1001 /app/main ./main

# Switch to non-root user
USER 1001

# Expose the application port
EXPOSE 4260

# Environment variables with defaults
ENV PORT=4260
ENV DEBUG=false

# Set the entrypoint and default command
ENTRYPOINT ["/app/main"]
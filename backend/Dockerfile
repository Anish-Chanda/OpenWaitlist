FROM golang:1.25.1-alpine AS builder

# Install dependencies for building
RUN apk add --no-cache git ca-certificates protobuf-dev make

# Set working directory
WORKDIR /app

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api ./backend

# Final stage
FROM alpine:latest

# Install ca-certificates and wget for health checks
RUN apk --no-cache add ca-certificates wget

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -S appuser -u 1001 -G appgroup

WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /app/bin/api .

# Change ownership to non-root user
RUN chown appuser:appgroup api

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./api"]

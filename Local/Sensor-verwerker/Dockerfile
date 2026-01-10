# Build stage - Multi-platform support
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

# Kopieer ALLEEN de Go-bestanden, niet de manifests
COPY *.go ./
COPY --chown=app:app . .



# Build for the target platform (ARM64 for Raspberry Pi, AMD64 for standard Linux)
# Docker automatically sets TARGETARCH at build time
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} go build -o sensor-verwerker .

# Runtime stage
FROM alpine:latest
USER root

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user for security
RUN addgroup -S app && adduser -S app -G app

# Set working directory with proper permissions
WORKDIR /home/app

# Copy binary from builder
COPY --from=builder --chown=app:app /app/sensor-verwerker .

# Switch to non-root user
USER app

# Health check (optional but recommended)
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD ps aux | grep sensor-verwerker | grep -v grep || exit 1

ENTRYPOINT ["./sensor-verwerker"]

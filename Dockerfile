# syntax=docker/dockerfile:1

FROM golang:1.22 AS builder
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app ./cmd

# Runtime stage
FROM gcr.io/distroless/base-debian12
WORKDIR /app

# Copy binary
COPY --from=builder /app/app /app/app

# Copy CA certificates for HTTPS calls
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Run as non-root
USER nonroot:nonroot

# Expose port (Cloud Run usa PORT env var, mas 8080 é padrão)
EXPOSE 8080

ENTRYPOINT ["/app/app"]
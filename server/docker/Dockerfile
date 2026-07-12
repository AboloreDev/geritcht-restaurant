FROM golang:1.26.1-alpine AS builder

WORKDIR /app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build all binaries
RUN mkdir -p /app/bin && \
    for cmd in cmd/*/; do \
        if [ -d "$cmd" ]; then \
            binary=$(basename $cmd); \
            echo "Building $binary..."; \
            CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/$binary ./$cmd; \
        fi \
    done

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy all binaries from builder
COPY --from=builder /app/bin/* ./

# Copy static docs (rapiddoc.html, swagger.json/yaml, etc.)
COPY --from=builder /app/docs ./docs

# Default command (can be overridden in docker-compose)
CMD ["./api"]
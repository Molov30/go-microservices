FROM golang:1.25 AS builder

WORKDIR /build

# Copy generated proto files
COPY generated/ ./generated/

# Copy service files
COPY services/gateway/go.mod services/gateway/go.sum ./services/gateway/
WORKDIR /build/services/gateway
RUN go mod download

# Copy the entire service code
WORKDIR /build
COPY services/gateway/ ./services/gateway/

# Build with version info from git
ARG VERSION=dev
ARG BUILD_TIME
WORKDIR /build/services/gateway
RUN BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S') && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" \
    -o gateway ./cmd/

FROM alpine:latest

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/services/gateway/gateway .

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 9002

CMD ["./gateway"]

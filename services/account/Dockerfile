FROM golang:1.25 AS builder

WORKDIR /build

# Copy generated proto files
COPY generated/ ./generated/

# Copy service files
COPY services/account/go.mod services/account/go.sum ./services/account/
WORKDIR /build/services/account
RUN go mod download

# Copy the entire service code
WORKDIR /build
COPY services/account/ ./services/account/

# Build with version info from git
ARG VERSION=dev
ARG BUILD_TIME
WORKDIR /build/services/account
RUN BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S') && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" \
    -o account ./cmd/

FROM alpine:latest

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/services/account/account .

# Create non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 8080 50051

CMD ["./account"]

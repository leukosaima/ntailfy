# Build stage
FROM golang:1.26-alpine@sha256:3ad57304ad93bbec8548a0437ad9e06a455660655d9af011d58b993f6f615648 AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY *.go ./

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ntailfy .

# Final stage
FROM alpine:latest@sha256:4f4ba248d8a2c90a6e52ffdfc194181f7617f9ddaca348d4c550a6b354fc7c2a

# Add ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/ntailfy .

# Run as non-root user
RUN adduser -D -u 1000 ntailfy && \
    chown -R ntailfy:ntailfy /app

USER ntailfy

ENTRYPOINT ["/app/ntailfy"]

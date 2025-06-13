# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.work and go.work.sum first
COPY go.work go.work.sum ./

# Copy all modules
COPY backend/libs/ backend/libs/
COPY backend/internal/ChatBotCore/ backend/internal/ChatBotCore/
COPY backend/external/FacebookConnect/ backend/external/FacebookConnect/

# Download dependencies
RUN go work sync

# Build the application
WORKDIR /app/backend/external/FacebookConnect
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata wget

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/backend/external/FacebookConnect/main .

# Expose port (different from main project)
EXPOSE 8091

# Run the binary
CMD ["./main"]

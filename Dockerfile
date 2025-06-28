FROM golang:1.20-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/log-ingestion-service ./cmd

# Use a minimal alpine image for the final stage
FROM alpine:3.18

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/log-ingestion-service .

# Expose the API port
EXPOSE 8080

# Run the application
CMD ["/app/log-ingestion-service"]
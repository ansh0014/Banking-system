# Start from the official Golang image
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY Bank/go.mod Bank/go.sum ./Bank/

# Download dependencies
RUN cd Bank && go mod download

# Copy the source code
COPY Bank/ ./Bank/

# Build the Go app
RUN cd Bank && go build -o /gobank

# Use a minimal image for the final container
FROM alpine:latest
WORKDIR /root/

# Copy the built binary from the builder
COPY --from=builder /gobank .

# Copy .env file if needed (uncomment if you use .env)
# COPY Bank/.env .

# Expose the application port
EXPOSE 1000

# Run the executable
CMD ["./gobank"]
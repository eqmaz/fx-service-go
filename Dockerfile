# Build stage, using the official Golang image
FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

# Copy the source from the current directory to the working directory inside the container
COPY . .

# Build the Go application
RUN go build -o server ./cmd/server
RUN ls -l /app/server

# Use a minimal image for the final stage
FROM alpine:latest
WORKDIR /app

# Install necessary packages
# libc6-compat is required for the Go binary to run
# ca-certificates is required to make HTTPS requests
RUN apk add --no-cache libc6-compat ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/server .
RUN ls -l /app/server

# Copy the config file
COPY config.json .

# Expose port 8080 to the outside world
EXPOSE 8080

# Ensure the binary has execution permissions
RUN chmod +x /app/server

# Command to run the executable
CMD ["./server", "--config", "./config.json"]

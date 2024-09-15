# Stage 1: Build the Go app
FROM golang:1.20-alpine AS builder

# Install necessary build tools and dependencies
RUN apk add --no-cache git

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the entire source code
COPY . .

# Build the application
RUN go build -o main .

# Stage 2: Run the app in a minimal image
FROM alpine:latest

# Install necessary runtime dependencies
RUN apk --no-cache add ca-certificates

# Set the working directory
WORKDIR /root/

# Copy the compiled Go binary from the builder stage
COPY --from=builder /app/main .

# Set the environment variables for AWS credentials
# You can use docker-compose.yml to pass these values instead
ENV AWS_ACCESS_KEY_ID=AKIAQXHOIJXXON7TKMET
ENV AWS_SECRET_ACCESS_KEY=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
ENV AWS_REGION=eu-north-1

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./main"]
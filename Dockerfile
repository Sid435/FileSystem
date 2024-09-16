# Start from the official Go image
FROM golang:1.20

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code and .env file into the container
COPY . .

# Build the application
RUN go build -o cmd/main .

# Expose port 8080 and 9010
EXPOSE 8080 9010

# Run the binary program produced by `go build`
CMD ["./main"]
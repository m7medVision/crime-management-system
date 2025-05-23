# Use the official Golang image as the base image
FROM golang:1.24-alpine

# Set the Current Working Directory inside the container
WORKDIR /app

# Install required fonts for PDF generation
RUN apk add --no-cache fontconfig curl

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Create necessary directories
RUN mkdir -p ./tmp ./templates

# Build the Go app
RUN go build -o main cmd/api/main.go

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]

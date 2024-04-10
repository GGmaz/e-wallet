# Start from the latest golang base image
FROM golang:1.20-alpine

# Install git
#RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy only necessary files
COPY go.mod go.sum ./

# Download and install the dependencies
RUN go mod download

# Copy the current directory contents into the container at /app
COPY . .

# Build the Go application
RUN go build -o main ./cmd

# Expose ports to the outside world
EXPOSE 8082

# Command to run the executable
CMD ["./main"]
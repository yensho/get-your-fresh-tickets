FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download and cache Go modules
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o gyft .

# Expose the port that the application listens on
EXPOSE 8080

# Set the entry point for the container
CMD ["./gyft"]
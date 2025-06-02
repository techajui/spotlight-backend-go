FROM golang:1.23  

WORKDIR /app

# Install git
RUN apt-get update && apt-get install -y git

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN go build -o main ./cmd/api

# Expose the port
EXPOSE 8080

# Run the application
CMD ["./main"]
FROM golang:1.22-alpine

WORKDIR /app

# Install git, build dependencies, and postgresql-client for migrations
RUN apk add --no-cache git gcc musl-dev postgresql-client

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Expose the port
EXPOSE 8080

# Run migrations and start the application
CMD sh -c "for migration in \$(ls -v migrations/*.sql); do echo \"Running migration: \$migration\"; PGPASSWORD=\$DB_PASSWORD psql -h \$DB_HOST -p \$DB_PORT -U \$DB_USER -d \$DB_NAME -f \$migration; done && ./main serve"
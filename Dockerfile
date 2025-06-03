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

# Create a script to run migrations and start the app
RUN echo '#!/bin/sh\n\
echo "Waiting for database..."\n\
sleep 10\n\
\n\
echo "Running database migrations..."\n\
for file in /app/migrations/*.sql; do\n\
  echo "Running migration: $file"\n\
  PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -U $DB_USER -d $DB_NAME -f "$file"\n\
done\n\
\n\
echo "Running database seed..."\n\
./main seed\n\
\n\
echo "Starting application..."\n\
./main\n\
' > /app/start.sh && chmod +x /app/start.sh

# Expose the port
EXPOSE 8080

# Run the start script
CMD ["/app/start.sh"]
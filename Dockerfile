FROM golang:1.22-alpine

WORKDIR /app

# Install git, build dependencies, and postgresql-client
RUN apk add --no-cache git gcc musl-dev postgresql-client

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api

# Create startup script
RUN echo '#!/bin/sh\n\
# Wait for database to be ready\n\
echo "Waiting for database..."\n\
while ! pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER; do\n\
  sleep 1\n\
done\n\
\n\
# Run migrations\n\
echo "Running database migrations..."\n\
for migration in migrations/*.sql; do\n\
  if [ -f "$migration" ]; then\n\
    echo "Running migration: $migration"\n\
    PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$migration"\n\
  fi\n\
done\n\
\n\
# Start the application\n\
echo "Starting application..."\n\
./main' > /app/start.sh && chmod +x /app/start.sh

# Expose the port
EXPOSE 8080

# Run the startup script
CMD ["/app/start.sh"]
#!/bin/bash

# Exit on error
set -e

# Configuration
PROJECT_ID="spotlight-backend-go-v1"
REGION="us-central1"
INSTANCE_NAME="spotlight-postgres"
DB_NAME="spotlight"
DB_USER="postgres"
DB_PASSWORD="$(gcloud secrets versions access latest --secret=db-password --project=$PROJECT_ID)"
MIGRATIONS_DIR="./migrations"
PROXY_PORT=5432

# Get the instance connection name
SQL_INSTANCE_CONNECTION_NAME="$PROJECT_ID:$REGION:$INSTANCE_NAME"

# Start Cloud SQL Proxy in background
if ! command -v cloud_sql_proxy &> /dev/null; then
  echo "cloud_sql_proxy could not be found. Please install it first."
  exit 1
fi

echo "Killing any existing Cloud SQL Proxy processes..."
pkill -f "cloud_sql_proxy.*$INSTANCE_NAME" || true

echo "Starting Cloud SQL Proxy..."
cloud_sql_proxy --instances=$SQL_INSTANCE_CONNECTION_NAME=tcp:$PROXY_PORT &
PROXY_PID=$!

# Give the proxy a moment to start
echo "Waiting for proxy to start..."
sleep 5

# Wait for database to be ready
echo "Waiting for database to be ready..."
for i in {1..30}; do
  if pg_isready -h localhost -p $PROXY_PORT -U $DB_USER; then
    echo "Database is ready!"
    break
  fi
  if [ $i -eq 30 ]; then
    echo "Database connection timeout after 30 seconds"
    exit 1
  fi
  echo "Waiting for database... ($i/30)"
  sleep 1
done

# Run migrations
echo "Running database migrations..."
PGPASSWORD=$DB_PASSWORD psql -h localhost -p $PROXY_PORT -U $DB_USER -d $DB_NAME -f migrations/complete_schema.sql

# Check if migration was successful
if [ $? -eq 0 ]; then
    echo "Migrations completed successfully!"
else
    echo "Migration failed!"
    exit 1
fi

# Run data/mock data migrations
echo "Running data/mock data migrations..."
for migration in $MIGRATIONS_DIR/002_*.sql; do
  if [ -f "$migration" ]; then
    echo "Running migration: $migration"
    PGPASSWORD=$DB_PASSWORD psql -h localhost -p $PROXY_PORT -U $DB_USER -d $DB_NAME -f "$migration" || echo "Migration $migration failed but continuing..."
  fi
done

# Terminate the proxy
echo "Stopping Cloud SQL Proxy..."
kill $PROXY_PID || true

echo "Migrations completed successfully!" 
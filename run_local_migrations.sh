#!/bin/bash

# Exit on error
set -e

# Set environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgresspotlight@mma
export DB_NAME=spotlight

# Run all migrations in sequence
echo "Running migrations..."
for migration in migrations/*.sql; do
  if [ -f "$migration" ]; then
    echo "Running migration: $migration"
    PGPASSWORD=$DB_PASSWORD psql -h localhost -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$migration"
  fi
done

echo "Migrations completed successfully!" 
#!/bin/bash

# Exit on error
set -e

# Set environment variables
export DB_HOST="localhost"
export DB_PORT="5432"
export DB_USER="postgres"
export DB_PASSWORD="postgresspotlight@mma"
export DB_NAME="spotlight"
export LOCAL_DEV="true"
export JWT_SECRET="spotlight_jwt_secret_key_2024_secure"
export PGPASSWORD=$DB_PASSWORD
export GOOGLE_APPLICATION_CREDENTIALS="$(dirname "$0")/config/client_secret_1089184396463-55fkfc50skejp2cqe448ctrgff40oonj.apps.googleusercontent.com.json"
export GOOGLE_CLIENT_ID="1089184396463-55fkfc50skejp2cqe448ctrgff40oonj.apps.googleusercontent.com"
export GOOGLE_CLIENT_SECRET="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"
export FRONTEND_URL="http://localhost:3000"
export FIREBASE_AUTH_DOMAIN="spotlight-v1.firebaseapp.com"
export FIREBASE_PROJECT_ID="spotlight-v1"
export FIREBASE_REDIRECT_URI="https://spotlight-v1.firebaseapp.com/__/auth/handler"
export GOOGLE_REDIRECT_URI="https://spotlight-v1.firebaseapp.com/__/auth/handler"

# Function to check if psql is available
check_psql() {
    if ! command -v psql &> /dev/null; then
        echo "Error: psql is not installed. Please install PostgreSQL client tools."
        exit 1
    fi
}

# Function to reset database
reset_database() {
    echo "Resetting database..."
    
    # Drop all tables in the database
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO postgres; GRANT ALL ON SCHEMA public TO public;"
    
    if [ $? -eq 0 ]; then
        echo "Database reset successful!"
    else
        echo "Error: Failed to reset database"
        exit 1
    fi
}

# Function to run migrations
run_migrations() {
    echo "Running migrations..."
    
    # Run all migrations in order
    for migration in migrations/*.sql; do
        if [ -f "$migration" ]; then
            echo "Running migration: $migration"
            psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f "$migration"
        fi
    done
}

# Function to start the server
start_server() {
    echo "Starting the server..."
    go run cmd/api/main.go
}

# Main execution
echo "Starting setup process..."

# Check if psql is available
check_psql

# Reset the database
reset_database

# Run migrations
run_migrations

# Run the seed file
echo "Running seed file..."
go run cmd/seed/seed.go

if [ $? -eq 0 ]; then
    echo "Database setup completed successfully!"
else
    echo "Error: Database setup failed"
    exit 1
fi

# Start the server
start_server 
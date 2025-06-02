#!/bin/bash

# Exit on error
set -e

# Configuration
PROJECT_ID="spotlight-backend-go-v1"
INSTANCE_NAME="spotlight-postgres"
DB_NAME="spotlight"
DB_USER="postgres"

echo "🔍 Verifying database configuration..."

# Check if Cloud SQL instance exists
echo "Checking Cloud SQL instance..."
if ! gcloud sql instances describe $INSTANCE_NAME --project=$PROJECT_ID > /dev/null 2>&1; then
    echo "❌ Cloud SQL instance $INSTANCE_NAME not found!"
    exit 1
fi

# Check if database exists
echo "Checking database $DB_NAME..."
if ! gcloud sql databases list --instance=$INSTANCE_NAME --project=$PROJECT_ID | grep -q $DB_NAME; then
    echo "❌ Database $DB_NAME not found!"
    exit 1
fi

# Check if user exists
echo "Checking database user $DB_USER..."
if ! gcloud sql users list --instance=$INSTANCE_NAME --project=$PROJECT_ID | grep -q $DB_USER; then
    echo "❌ Database user $DB_USER not found!"
    exit 1
fi

# Check if db-password secret exists
echo "Checking db-password secret..."
if ! gcloud secrets describe db-password --project=$PROJECT_ID > /dev/null 2>&1; then
    echo "❌ Secret db-password not found!"
    exit 1
fi

# Get the password from secret
DB_PASSWORD=$(gcloud secrets versions access latest --secret=db-password --project=$PROJECT_ID)

# Update database password to match secret
echo "Updating database password..."
gcloud sql users set-password $DB_USER --instance=$INSTANCE_NAME --password="$DB_PASSWORD" --project=$PROJECT_ID

echo "✅ Database configuration verified successfully!" 
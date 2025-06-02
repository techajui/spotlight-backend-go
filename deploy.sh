#!/bin/bash

# Exit on error
set -e

# Configuration
PROJECT_ID="spotlight-backend-go-v1"
REGION="us-central1"
SERVICE_NAME="spotlight-backend-go"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print status messages
print_status() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}! $1${NC}"
}

# Verify database configuration
print_status "Verifying database configuration..."
if ! ./verify-db.sh; then
    print_error "Database verification failed!"
    exit 1
fi

# Build and push the Docker image
print_status "Building and pushing Docker image..."
if ! gcloud builds submit --tag gcr.io/$PROJECT_ID/$SERVICE_NAME .; then
    print_error "Failed to build and push Docker image!"
    exit 1
fi

# Deploy to Cloud Run
print_status "Deploying to Cloud Run..."
if ! gcloud run deploy $SERVICE_NAME \
    --image gcr.io/$PROJECT_ID/$SERVICE_NAME \
    --platform managed \
    --region $REGION \
    --allow-unauthenticated \
    --set-env-vars="DB_HOST=/cloudsql/$PROJECT_ID:$REGION:spotlight-postgres" \
    --set-env-vars="DB_PORT=5432" \
    --set-env-vars="DB_USER=postgres" \
    --set-env-vars="DB_NAME=spotlight" \
    --set-env-vars="DB_SSLMODE=disable" \
    --set-secrets="DB_PASSWORD=db-password:latest" \
    --set-secrets="JWT_SECRET=jwt-secret:latest" \
    --add-cloudsql-instances "$PROJECT_ID:$REGION:spotlight-postgres"; then
    print_error "Failed to deploy to Cloud Run!"
    exit 1
fi

print_status "Deployment completed successfully!"
print_status "Service URL: https://$SERVICE_NAME-$PROJECT_ID.$REGION.run.app" 
#!/bin/bash

# Exit on error
set -e

# Configuration
PROJECT_ID="spotlight-backend-go-v1"
REGION="us-central1"
INSTANCE_NAME="spotlight-postgres"
SERVICE_NAME="spotlight-backend-go"

# Create new project if it doesn't exist
echo "Setting up project..."
gcloud projects create $PROJECT_ID --quiet 2>/dev/null || echo "Project already exists."

# Set as current project
gcloud config set project $PROJECT_ID

# Enable required APIs
echo "Enabling required APIs..."
gcloud services enable \
  run.googleapis.com \
  sqladmin.googleapis.com \
  secretmanager.googleapis.com \
  containerregistry.googleapis.com \
  cloudbuild.googleapis.com \
  iam.googleapis.com

# Configure Docker to use gcloud credentials
echo "Configuring Docker authentication..."
gcloud auth configure-docker --quiet

# Build and push image using Cloud Build
echo "Building and pushing image using Cloud Build..."
gcloud builds submit --tag gcr.io/$PROJECT_ID/$SERVICE_NAME .

# Create Cloud SQL instance (if it doesn't exist)
echo "Setting up Cloud SQL..."
gcloud sql instances create $INSTANCE_NAME \
  --database-version=POSTGRES_13 \
  --cpu=1 \
  --memory=3840MiB \
  --region=$REGION \
  --root-password="$(openssl rand -base64 32)" 2>/dev/null || echo "SQL instance already exists."

# Create database (if it doesn't exist)
echo "Creating database..."
gcloud sql databases create spotlight \
  --instance=$INSTANCE_NAME 2>/dev/null || echo "Database already exists."

# Create secrets in Secret Manager
echo "Creating secrets in Secret Manager..."

# Database password
echo "postgresspotlight@mma" | gcloud secrets create db-password \
  --replication-policy="automatic" \
  --data-file=- 2>/dev/null || \
  echo "postgresspotlight@mma" | gcloud secrets versions add db-password --data-file=-

# JWT secret
echo "$(openssl rand -base64 32)" | gcloud secrets create jwt-secret \
  --replication-policy="automatic" \
  --data-file=- 2>/dev/null || \
  echo "$(openssl rand -base64 32)" | gcloud secrets versions add jwt-secret --data-file=-

# Create a service account for Cloud Run
echo "Creating service account for Cloud Run..."
SERVICE_ACCOUNT="spotlight-api-sa"
gcloud iam service-accounts create $SERVICE_ACCOUNT \
  --display-name="Spotlight API Service Account" 2>/dev/null || echo "Service account already exists."

# Grant necessary permissions
echo "Granting permissions..."
gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:$SERVICE_ACCOUNT@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/secretmanager.secretAccessor"

gcloud projects add-iam-policy-binding $PROJECT_ID \
  --member="serviceAccount:$SERVICE_ACCOUNT@$PROJECT_ID.iam.gserviceaccount.com" \
  --role="roles/cloudsql.client"

# Deploy to Cloud Run
echo "Deploying to Cloud Run..."
gcloud run deploy $SERVICE_NAME \
  --image gcr.io/$PROJECT_ID/$SERVICE_NAME \
  --platform managed \
  --region $REGION \
  --service-account="$SERVICE_ACCOUNT@$PROJECT_ID.iam.gserviceaccount.com" \
  --allow-unauthenticated \
  --memory=512Mi \
  --timeout=300s \
  --set-env-vars="DB_PORT=5432,DB_USER=postgres,DB_NAME=spotlight,DB_SSLMODE=disable,DB_HOST=/cloudsql/$PROJECT_ID:$REGION:$INSTANCE_NAME" \
  --set-secrets="DB_PASSWORD=db-password:latest,JWT_SECRET=jwt-secret:latest" \
  --add-cloudsql-instances="$PROJECT_ID:$REGION:$INSTANCE_NAME"

echo "Deployment completed successfully!" 
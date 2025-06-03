# Set environment variables
$env:DB_HOST = "localhost"
$env:DB_PORT = "5432"
$env:DB_USER = "postgres"
$env:DB_PASSWORD = "postgresspotlight@mma"
$env:DB_NAME = "spotlight"
$env:LOCAL_DEV = "true"
$env:JWT_SECRET = "spotlight_jwt_secret_key_2024_secure"
$env:PGPASSWORD = $env:DB_PASSWORD
$env:GOOGLE_APPLICATION_CREDENTIALS = "$PSScriptRoot\config\client_secret_1089184396463-55fkfc50skejp2cqe448ctrgff40oonj.apps.googleusercontent.com.json"
$env:GOOGLE_CLIENT_ID = "1089184396463-55fkfc50skejp2cqe448ctrgff40oonj.apps.googleusercontent.com"
$env:GOOGLE_CLIENT_SECRET = "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"
$env:FRONTEND_URL = "http://localhost:3000"
$env:FIREBASE_AUTH_DOMAIN = "spotlight-v1.firebaseapp.com"
$env:FIREBASE_PROJECT_ID = "spotlight-v1"
$env:FIREBASE_REDIRECT_URI = "https://spotlight-v1.firebaseapp.com/__/auth/handler"
$env:GOOGLE_REDIRECT_URI = "https://spotlight-v1.firebaseapp.com/__/auth/handler"


# Set psql path
$PSQL_PATH = "C:\Program Files\PostgreSQL\17\pgAdmin 4\runtime\psql.exe"

# Check if psql exists
if (-not (Test-Path $PSQL_PATH)) {
    Write-Error "Error: psql not found at $PSQL_PATH"
    exit 1
}

Write-Host "Using psql at: $PSQL_PATH"

# Reset database
Write-Host "Resetting database..."
& $PSQL_PATH -h $env:DB_HOST -p $env:DB_PORT -U $env:DB_USER -d $env:DB_NAME -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO postgres; GRANT ALL ON SCHEMA public TO public;"
if ($LASTEXITCODE -ne 0) {
    Write-Error "Error: Failed to reset database"
    exit 1
}
Write-Host "Database reset successful!"

# Run migrations
Write-Host "Running migrations..."
Get-ChildItem -Path "migrations\*.sql" | ForEach-Object {
    Write-Host "Running migration: $_"
    & $PSQL_PATH -h $env:DB_HOST -p $env:DB_PORT -U $env:DB_USER -d $env:DB_NAME -f $_.FullName
}

# Run the seed file
Write-Host "Running seed file..."
go run cmd/seed/seed.go
if ($LASTEXITCODE -ne 0) {
    Write-Error "Error: Database setup failed"
    exit 1
}
Write-Host "Database setup completed successfully!"

# Start the server
Write-Host "Starting the server..."
go run cmd/api/main.go 
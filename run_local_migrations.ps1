# Set environment variables
$env:DB_HOST = "localhost"
$env:DB_PORT = "5432"
$env:DB_USER = "postgres"
$env:DB_PASSWORD = "postgresspotlight@mma"
$env:DB_NAME = "spotlight"

# Run all migrations in sequence
Write-Host "Running migrations..."
Get-ChildItem -Path "migrations" -Filter "*.sql" | Sort-Object Name | ForEach-Object {
    Write-Host "Running migration: $($_.Name)"
    $env:PGPASSWORD = $env:DB_PASSWORD
    psql -h $env:DB_HOST -p $env:DB_PORT -U $env:DB_USER -d $env:DB_NAME -f $_.FullName
}

Write-Host "Migrations completed successfully!" 
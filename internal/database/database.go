package database

import (
	"fmt"
	"log"
	"os"
	"spotlight-backend-go/internal/models"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	// Get database configuration from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Debug log the environment variables (excluding password)
	log.Printf("Database configuration - Host: %s, Port: %s, User: %s, DB: %s",
		dbHost, dbPort, dbUser, dbName)

	// Set default port if not provided
	if dbPort == "" {
		dbPort = "5432"
		log.Println("Using default PostgreSQL port: 5432")
	}

	// Validate required environment variables
	if dbHost == "" {
		log.Fatal("Missing required environment variable: DB_HOST")
	}
	if dbUser == "" {
		log.Fatal("Missing required environment variable: DB_USER")
	}
	if dbPassword == "" {
		log.Fatal("Missing required environment variable: DB_PASSWORD")
	}
	if dbName == "" {
		log.Fatal("Missing required environment variable: DB_NAME")
	}

	// Configure GORM logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Create DSN for PostgreSQL
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	log.Printf("Attempting to connect to database at %s:%s", dbHost, dbPort)

	// Open database connection
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(
		&models.User{},
		&models.Event{},
		&models.EventAttendee{},
		&models.Bid{},
		&models.Application{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Seed initial data
	SeedData()

	log.Println("Database connected and migrated successfully!")
	return DB
}

package main

import (
	"log"
	"spotlight-backend-go/internal/database"
)

func main() {
	log.Println("Starting database seeding...")

	// Initialize database
	database.InitDB()

	// Run seed function
	database.SeedData()

	log.Println("Database seeding completed successfully!")
}

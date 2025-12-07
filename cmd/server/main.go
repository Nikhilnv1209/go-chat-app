package main

import (
	"log"

	"chat-app/internal/database"
	"chat-app/internal/models"

	"github.com/joho/godotenv"
)

func main() {
	// 1. Load Environment Variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// 2. Initialize Database
	database.InitDB()
	db := database.GetDB()

	// 3. Auto Migrate Models
	log.Println("Running AutoMigrate...")
	err := db.AutoMigrate(
		&models.User{},
		&models.Group{},
		&models.GroupMember{},
		&models.Message{},
		&models.MessageReceipt{},
		&models.Conversation{},
	)
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}
	log.Println("Database Migration Clean!")

	// 4. Server Start (Placeholder)
	log.Println("Server started on :8080 (Placeholder)")
	select {} // Block forever for now
}

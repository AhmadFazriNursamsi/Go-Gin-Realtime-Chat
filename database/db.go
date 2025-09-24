package database

import (
	"fmt"
	"log"
	"myapi/models"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	var database *gorm.DB
	var err error

	// Retry loop: tunggu database siap
	for i := 0; i < 10; i++ {
		database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("✅ Database connected")
			break
		}
		log.Printf("❌ Database not ready yet, retrying in 3 seconds... (%d/10)\n", i+1)
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect to database after 10 attempts:", err)
	}

	// AutoMigrate semua model
	if err := models.Migrate(database); err != nil {
		log.Fatal("Migration failed:", err)
	}

	DB = database
}

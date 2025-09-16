package database

import (
	"log"
	"myapi/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsk := "host=localhost user=postgres password=Fazri18 dbname=testing port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	database, err := gorm.Open(postgres.Open(dsk), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected")
	// AutoMigrate semua model di sini
	if err := models.Migrate(database); err != nil {
		log.Fatal("Migration failed:", err)
	}
	DB = database
}

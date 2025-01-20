package models

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Setup() {
	DNS := "host=localhost user=root password=root dbname=cryptisecure port=5432 sslmode=disable"
	var err error
	DB, err = gorm.Open(postgres.Open(DNS), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Successfully connected to PostgreSQL database")

	AutoMigrate()
}

func AutoMigrate() {
	DB.AutoMigrate(&Certificate{})
	DB.AutoMigrate(&CerKey{})
}

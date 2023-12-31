package database

import (
	"go-web/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := os.Getenv("DB")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		panic(err)
	}
	log.Println("Connected to Database")

	db.AutoMigrate(&models.User{}, &models.Photo{})

	DB = db
}

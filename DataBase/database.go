package DataBase

import (
	"digikala/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("online-shopping.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect db:", err)
	}

	err = DB.AutoMigrate(&models.Customer{})
	if err != nil {
		log.Fatal("Failed to migrate db:", err)
	}

	err = DB.AutoMigrate(&models.Token{})
	if err != nil {
		log.Fatal("Failed to migrate db:", err)
	}

	err = DB.AutoMigrate(&models.Product{})
	if err != nil {
		log.Fatal("Failed to migrate db:", err)
	}
}

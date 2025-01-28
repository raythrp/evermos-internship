package database

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	errors := godotenv.Load()
	if errors != nil {
		log.Fatal("Error loading .env file")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")

	var err error
	MYSQL := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/evermos?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword)
	DSN := MYSQL

	DB, err = gorm.Open(mysql.Open(DSN), &gorm.Config{})

	if (err != nil) {
		panic("Cannot connect to database")
	}
	if (DB == nil) {
		panic("DB variable at config is null")
	}

	fmt.Println("Connected to database", DB)
}
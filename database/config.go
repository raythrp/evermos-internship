package database

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	_ = godotenv.Load()

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "127.0.0.1:3306"
	}

	var err error
	MYSQL := fmt.Sprintf("%s:%s@tcp(%s)/evermos?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbHost)
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
package database

import (
	"fmt"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	const MYSQL = "root:root@tcp(127.0.0.1:3306)/evermos?charset=utf8mb4&parseTime=True&loc=Local"
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
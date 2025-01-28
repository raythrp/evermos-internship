package main

import (
	"log"
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/raythrp/evermos-internship/database"
	"github.com/raythrp/evermos-internship/routers"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.ConnectDB()
	app := fiber.New()

	routers.RouterApp(app)
	app.Listen(":" + os.Getenv("PORT"))
}
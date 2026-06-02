package main

import (
	"os"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/raythrp/evermos-internship/database"
	"github.com/raythrp/evermos-internship/routers"
)

func main() {
	// .env is optional — in Docker, vars are injected via the environment
	_ = godotenv.Load()

	database.ConnectDB()
	app := fiber.New()

	routers.RouterApp(app)
	app.Listen(":" + os.Getenv("PORT"))
}
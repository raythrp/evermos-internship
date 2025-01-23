package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raythrp/evermos-internship/database"
	"github.com/raythrp/evermos-internship/routers"
)

func main() {
	database.ConnectDB()
	app := fiber.New()

	routers.RouterApp(app)
	app.Listen(":8080")
}
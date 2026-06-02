package middlewares

import (
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

func AuthRequired() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(os.Getenv("JWT_SECRET")),
		ErrorHandler: jwtErrorHandler,
		TokenLookup:  "header:token",
	})
}

func jwtErrorHandler(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"status": false,
		"message": "Unauthorized",
	})
}
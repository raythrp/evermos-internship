package middlewares

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
)

const jwtSecret = "WsItpI3Moq4I0rVwo2fOcbvw8CDgJT9FMrsz9zsqAy3e7PRU8sojZ79jSDtnOuO0bjceupoidsajp3u2019eu[20ihceoijlciuab]"

func AuthRequired() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(jwtSecret),
		ErrorHandler: jwtErrorHandler,
		TokenLookup: "header:token",
	})
}

func jwtErrorHandler(c *fiber.Ctx, err error) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"status": false,
		"message": "Unauthorized",
	})
}
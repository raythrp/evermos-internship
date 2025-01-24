package helpers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func JwtClaimer(c *fiber.Ctx) string {
	credentials := c.Locals("user").(*jwt.Token)
	claims := credentials.Claims.(jwt.MapClaims)
	noTelp := claims["sub"].(string)
	return noTelp
}
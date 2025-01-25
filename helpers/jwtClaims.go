package helpers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	models "github.com/raythrp/evermos-internship/models/entities"
)

func JwtClaimer(c *fiber.Ctx) string {
	credentials := c.Locals("user").(*jwt.Token)
	claims := credentials.Claims.(jwt.MapClaims)
	noTelp := claims["sub"].(string)
	return noTelp
}

func AdminValidator(c *fiber.Ctx, user models.User) error {
	if !user.IsAdmin {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false,
			"message": "Unauthorized",
		})
	}

	return nil
}
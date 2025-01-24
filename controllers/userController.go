package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raythrp/evermos-internship/database"
	"github.com/raythrp/evermos-internship/helpers"
	models "github.com/raythrp/evermos-internship/models/entities"
	"github.com/raythrp/evermos-internship/models/requests"

)

func UserGetMyProfile(c *fiber.Ctx) error {
	noTelp := helpers.JwtClaimer(c)

	// User valid
	var user models.User
	if err := database.DB.Preload("Toko").Where("notelp = ?", noTelp).First(&user).Error; err == nil {
		return c.JSON(fiber.Map{
			"status": true,
			"message": "Succeed to GET data",
			"errors": nil,
			"data": user,
		})
	}

	// User invalid
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"status": false,
		"message": "Unauthorized",
	})
}

func UserUpdateMyProfile(c *fiber.Ctx) error {
	var body requests.UpdateUser
	noTelp := helpers.JwtClaimer(c)

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to PUT data",
			"errors": []string {err.Error()},
			"data": nil,
		})
	}

	var existingUser models.User
	if err := database.DB.Where("notelp = ?", noTelp).First(&existingUser).Error; err != nil {
	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"status":  false,
		"message": "User not found",
		"errors":  err.Error(),
		"data":    nil,
	})
	}

	// Inserting to user table
	formattedTanggalLahir, err := helpers.TimeParserToDate(body.TanggalLahir)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to PUT data",
			"errors": []string {err.Error()},
			"data": nil,
		})
	}

	// Updating fields
	existingUser.Nama = body.Nama
	existingUser.KataSandi = body.KataSandi
	existingUser.NoTelp = body.NoTelp
	existingUser.TanggalLahir = formattedTanggalLahir
	existingUser.Pekerjaan = body.Pekerjaan
	existingUser.Email = body.Email
	existingUser.IDProvinsi = body.IDProvinsi
	existingUser.IDKota = body.IDKota
	if err := database.DB.Save(&existingUser).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to PUT data",
			"errors": []string {err.Error()},
			"data": nil,
		})
	}

	// OK response
	return c.JSON(fiber.Map{
		"status": true,
		"message": "Succeed to PUT data",
		"errors": nil,
		"data": "",
	})
}
package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raythrp/evermos-internship/database"
	"github.com/raythrp/evermos-internship/helpers"
	models "github.com/raythrp/evermos-internship/models/entities"
	"github.com/raythrp/evermos-internship/models/responses"
)

func TokoGetMyToko(c *fiber.Ctx) error {
	noTelp := helpers.JwtClaimer(c)

	// User valid
	var user models.User
	if err := database.DB.Where("notelp = ?", noTelp).First(&user).Error; err != nil {
		return c.JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": err.Error(),
			"data": nil,
		})
	}

	var toko responses.Toko
	if err := database.DB.Where("id_user = ?", user.ID).First(&toko).Error; err == nil {
		return c.JSON(fiber.Map{
			"status": true,
			"message": "Succeed to GET data",
			"errors": nil,
			"data": toko,
		})
	}

	// User invalid
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"status": false,
		"message": "Unauthorized",
	})
}

func TokoGetTokoByID(c *fiber.Ctx) error {
	noTelp := helpers.JwtClaimer(c)

	// User valid
	var user models.User
	if err := database.DB.Where("notelp = ?", noTelp).First(&user).Error; err != nil {
		return c.JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": []string {"Toko tidak ditemukan"},
			"data": nil,
		})
	}

	// Non-admin Token
	helpers.AdminValidator(c, user)

	// Taking Parameters
	tokoID := c.Params("id_toko")
	if tokoID == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": []string {"Toko tidak ditemukan"},
			"data": nil,
		})
	}

	// Getting Toko
	var toko responses.Toko
	if err := database.DB.Where("id = ?", tokoID).First(&toko).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": []string {"Toko tidak ditemukan"},
			"data": nil,
		})
	}

	// OK Response
	return c.JSON(fiber.Map{
		"status": true,
		"message": "Succeed to GET data",
		"errors": nil,
		"data": toko,
	})
}

func TokoGetAllToko(c *fiber.Ctx) error {
	noTelp := helpers.JwtClaimer(c)

	// User valid
	var user models.User
	if err := database.DB.Where("notelp = ?", noTelp).First(&user).Error; err != nil {
		return c.JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": []string {"Toko tidak ditemukan"},
			"data": nil,
		})
	}

	// Non-admin Token
	helpers.AdminValidator(c, user)

	// Pagination Implementation
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	// Nama Query
	nama := c.Query("nama", "")
	if nama != "" {
		var tokos []responses.Toko
		if err := database.DB.Limit(limit).Offset(offset).Where("nama_toko LIKE ?", "%"+nama+"%").Find(&tokos).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status": false,
				"message": "Failed to GET data",
				"errors": []string {"Toko tidak ditemukan"},
				"data": nil,
			})
		}
		return c.JSON(fiber.Map{
			"status": true,
			"message": "Succeed to GET data",
			"errors": nil,
			"data": fiber.Map{
				"page": page,
				"limit": limit,
				"data": tokos,
			},
		})

	}

	var tokos []responses.Toko
	if err := database.DB.Limit(limit).Offset(offset).Find(&tokos).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": []string {"Toko tidak ditemukan"},
			"data": nil,
		})
	}

	// OK Response
	return c.JSON(fiber.Map{
		"status": true,
		"message": "Succeed to GET data",
		"errors": nil,
		"data": fiber.Map{
			"page": page,
			"limit": limit,
			"data": tokos,
		},
	})
}
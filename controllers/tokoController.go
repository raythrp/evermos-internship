package controllers

import (
	"fmt"
	"time"

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

	var toko responses.MyToko
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
	if !user.IsAdmin {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false,
			"message": "Unauthorized",
		})
	}

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
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": []string {"Toko tidak ditemukan"},
			"data": nil,
		})
	}

	// Non-admin Token
	if !user.IsAdmin {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false,
			"message": "Unauthorized",
		})
	}

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

func TokoUpdateProfile(c *fiber.Ctx) error {
	noTelp := helpers.JwtClaimer(c)

	// User valid
	var user models.User
	if err := database.DB.Where("notelp = ?", noTelp).First(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to UPDATE data",
			"errors": nil,
			"data": nil,
		})
	}

	// Non-admin Token
	if !user.IsAdmin {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false,
			"message": "Unauthorized",
		})
	}

	// Taking Parameters
	tokoID := c.Params("id_toko")
	if tokoID == "" {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
			"message": "Failed to UPDATE data",
			"errors": nil,
			"data": nil,
		})
	}

	// Parse form data
	namaToko := c.FormValue("nama_toko")

	// Handle file upload
	file, err := c.FormFile("photo")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to UPDATE data",
			"errors": nil,
			"data": nil,
		})
	}

	// Generate unique filename
	filename := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)

	// Save the file to the server
	savePath := fmt.Sprintf("./uploads/%s", filename)
	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to UPDATE data",
			"errors": nil,
			"data": nil,
		})
	}

	// Finding Toko
	var toko models.Toko
	if err := database.DB.Where("id = ?", tokoID).First(&toko).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
			"message": "Failed to UPDATE data",
			"errors": nil,
			"data": nil,
		})
	}

	// Updating Table
	toko.NamaToko = namaToko
    	toko.UrlFoto = savePath
	if err := database.DB.Save(&toko).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to UPDATE data",
			"errors": nil,
			"data": nil,
		})
	}

	// OK Response
	return c.JSON(fiber.Map{
		"status": true,
		"message": "Succeed to UPDATE data",
		"errors": nil,
		"data": "Update toko succeed",
	})
}
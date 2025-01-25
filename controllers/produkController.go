package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raythrp/evermos-internship/database"
	"github.com/raythrp/evermos-internship/helpers"
	models "github.com/raythrp/evermos-internship/models/entities"
	"github.com/raythrp/evermos-internship/models/responses"
)

func ProdukGetAll(c *fiber.Ctx) error {
	noTelp := helpers.JwtClaimer(c)

	// User valid
	var user models.User
	if err := database.DB.Where("notelp = ?", noTelp).First(&user).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": err.Error(),
			"data": nil,
		})
	}


	var toko responses.Toko
	if err := database.DB.Where("id_user = ?", user.ID).First(&toko).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": err.Error(),
			"data": nil,
		})
	}

	// Pagination Implementation
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	// Get All Products 
	var products []responses.Produk
	if err := database.DB.Limit(limit).Offset(offset).Preload("Toko").Preload("Category").Preload("FotoProduk").Where("id_toko = ?", toko.ID).Find(&products).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": err.Error(),
			"data": nil,
		})
	}

	// OK Response
	return c.JSON(fiber.Map{
		"status": true,
		"message": "Succeed to GET data",
		"errors": nil,
		"data": fiber.Map{
			"data": products,
		},
	})
}

func ProdukGetByID(c *fiber.Ctx) error {
	noTelp := helpers.JwtClaimer(c)

	// User valid
	var user models.User
	if err := database.DB.Where("notelp = ?", noTelp).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": []string {"No Data Product"},
			"data": nil,
		})
	}


	var toko responses.Toko
	if err := database.DB.Where("id_user = ?", user.ID).First(&toko).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": []string {"No Data Product"},
			"data": nil,
		})
	}

	// Pagination Implementation
	produkID := c.Params("id")

	// Get Product
	var product responses.Produk
	if err := database.DB.Preload("Toko").Preload("Category").Preload("FotoProduk").Where("id_toko = ? AND id = ?", toko.ID, produkID).First(&product).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": []string {"No Data Product"},
			"data": nil,
		})
	}

	// OK Response
	return c.JSON(fiber.Map{
		"status": true,
		"message": "Succeed to GET data",
		"errors": nil,
		"data": fiber.Map{
			"data": product,
		},
	})
}
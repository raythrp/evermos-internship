package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/raythrp/evermos-internship/database"
	"github.com/raythrp/evermos-internship/helpers"
	models "github.com/raythrp/evermos-internship/models/entities"
	"github.com/raythrp/evermos-internship/models/responses"
)

func TrxGetAll(c *fiber.Ctx) error {
	noTelp := helpers.JwtClaimer(c)

	// User valid
	var user models.User
	if err := database.DB.Where("notelp = ?", noTelp).First(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
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

      // Get LogProduk
      var product_logs []responses.LogProduk
      if err := database.DB.Preload("Category").Preload("Toko").Find(&product_logs).Error; err != nil {
            c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": nil,
			"data": nil,
		})
      }

      for i := range product_logs {
            var product_photos []responses.FotoProduk
            if err := database.DB.Where("id_produk = ?", product_logs[i].IDProduk).Find(&product_photos).Error; err != nil {
                  c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                        "status": false,
                        "message": "Failed to GET data",
                        "errors": nil,
                        "data": nil,
                  })
            }
            product_logs[i].FotoProduk = product_photos
      }

      // Get DetailTrx
      var transaction_details []responses.DetailTrx
      if err := database.DB.Preload("Toko").Find(&transaction_details).Error; err != nil {
            c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": nil,
			"data": nil,
		})
      }

      // Binding LogProduk to DetailTrx
      for i := range transaction_details {
            for j := range product_logs {
                  if transaction_details[i].IDLogProduk == product_logs[j].ID {
                        transaction_details[i].LogProduk = product_logs[j]
                        // product_logs = append(product_logs[:j], product_logs[j + 1:]...) // Removes taken LogProduk
                  }
            }
      }

      // Get Trx
	var transactions []responses.Trx
	if err := database.DB.Preload("AlamatPengiriman").Find(&transactions).Error; err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": nil,
			"data": nil,
		})
	}

      // Bind DetailTrx to Trx
      for i := range transactions {
            for j := range transaction_details {
                  if transactions[i].ID == transaction_details[j].IDTrx {
                        transactions[i].DetailTrx = append(transactions[i].DetailTrx, transaction_details[j])
                        // transaction_details = append(transaction_details[:j], transaction_details[j + 1:]...) // Removes taken DetailTrx
                  }
            }
      }

      // Apply in-memory pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))       // Default page is 1
	limit, _ := strconv.Atoi(c.Query("limit", "10"))   // Default limit is 10
	start := (page - 1) * limit                        // Start index
	end := start + limit                               // End index

	// Handle slicing to avoid out-of-range errors
	if start > len(transactions) {
		start = len(transactions)
	}
	if end > len(transactions) {
		end = len(transactions)
	}

	// Paginated data
	paginatedTransactions := transactions[start:end]

	// OK Response
	return c.JSON(fiber.Map{
		"status": true,
		"message": "Succeed to GET data",
		"errors": nil,
		"data": fiber.Map{
			"data": paginatedTransactions,
                  "page": page,
                  "limit": limit,
		},
	})
}
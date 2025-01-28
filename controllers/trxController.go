package controllers

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	"github.com/raythrp/evermos-internship/database"
	"github.com/raythrp/evermos-internship/helpers"
	models "github.com/raythrp/evermos-internship/models/entities"
	"github.com/raythrp/evermos-internship/models/requests"
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
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": nil,
			"data": nil,
		})
      }

      for i := range product_logs {
            var product_photos []responses.FotoProduk
            if err := database.DB.Where("id_produk = ?", product_logs[i].IDProduk).Find(&product_photos).Error; err != nil {
                  return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
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
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": nil,
			"data": nil,
		})
      }

      // Binding LogProduk to DetailTrx
      for i := range transaction_details {
            for j := 0; j < len(product_logs); j++ {
                if transaction_details[i].IDLogProduk == product_logs[j].ID {
                    transaction_details[i].LogProduk = product_logs[j]
                    product_logs = append(product_logs[:j], product_logs[j+1:]...) // Removes taken LogProduk
                    j--
                }
            }
        }

      // Get Trx
	var transactions []responses.Trx
	if err := database.DB.Preload("AlamatPengiriman").Find(&transactions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": nil,
			"data": nil,
		})
	}

      // Binding DetailTrx to Trx
      for i := range transactions {
            for j := 0; j < len(transaction_details); j++ {
                if transactions[i].ID == transaction_details[j].IDTrx {
                    transactions[i].DetailTrx = append(transactions[i].DetailTrx, transaction_details[j])
                    transaction_details = append(transaction_details[:j], transaction_details[j+1:]...) // Remove taken DetailTrx
                    j--
                }
            }
        }

      // Apply in-memory pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))      // Default page is 1
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

func TrxGetByID(c *fiber.Ctx) error {
      noTelp := helpers.JwtClaimer(c)
      transactionID := c.Params("id", "")
      if transactionID == "" {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": nil,
			"data": nil,
		})
      }

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

      // Get Trx
	var transaction responses.Trx
	if err := database.DB.Preload("AlamatPengiriman").Where("id = ? AND id_user = ?", transactionID, user.ID).First(&transaction).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false,
			"message": "Unauthorized",
		})
	}

      // Get LogProduk
      var product_logs []responses.LogProduk
      if err := database.DB.Preload("Category").Preload("Toko").Find(&product_logs).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": nil,
			"data": nil,
		})
      }

      for i := range product_logs {
            var product_photos []responses.FotoProduk
            if err := database.DB.Where("id_produk = ?", product_logs[i].IDProduk).Find(&product_photos).Error; err != nil {
                  return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                        "status": false,
                        "message": "Failed to GET data",
                        "errors": []string {"No Data Trx"},
                        "data": nil,
                  })
            }
            product_logs[i].FotoProduk = product_photos
      }

      // Get DetailTrx
      var transaction_details []responses.DetailTrx
      if err := database.DB.Preload("Toko").Find(&transaction_details).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                  "status": false,
                  "message": "Failed to GET data",
                  "errors": []string {"No Data Trx"},
                  "data": nil,
            })
      }

      // Binding LogProduk to DetailTrx
      for i := range transaction_details {
            for j := 0; j < len(product_logs); j++ {
                if transaction_details[i].IDLogProduk == product_logs[j].ID {
                    transaction_details[i].LogProduk = product_logs[j]
                    product_logs = append(product_logs[:j], product_logs[j+1:]...) // Removes taken LogProduk
                    j--
                }
            }
        }

      // Binding DetailTrx to Trx
      for j := 0; j < len(transaction_details); j++ {
            if transaction.ID == transaction_details[j].IDTrx {
                  transaction.DetailTrx = append(transaction.DetailTrx, transaction_details[j])
                  transaction_details = append(transaction_details[:j], transaction_details[j+1:]...) // Remove taken DetailTrx
                  j--
            }
      }

      // OK Response
	return c.JSON(fiber.Map{
		"status": true,
		"message": "Succeed to GET data",
		"errors": nil,
		"data": transaction,
	})
}

func TrxCreate(c *fiber.Ctx) error {
      noTelp := helpers.JwtClaimer(c)
      var body requests.CreateTrx

	// User valid
	var user models.User
	if err := database.DB.Where("notelp = ?", noTelp).First(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to POST data",
			"errors": err.Error(),
			"data": nil,
		})
	}

      // Toko valid
      var toko models.Toko
      if err := database.DB.Where("id_user = ?", user.ID).First(&toko).Error; err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to POST data",
			"errors": err.Error(),
			"data": nil,
		})
      }

      err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to POST data",
			"errors": err.Error(),
			"data": nil,
		})
	}

      // Alamat Valid
      var alamat models.Alamat
      if err := database.DB.Where("id = ?", body.AlamatKirim).First(&alamat).Error; err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to POST data",
			"errors": err.Error(),
			"data": nil,
		})
      }
      if alamat.IDUser != user.ID {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to POST data",
			"errors": "No Selected Alamat Record Found",
			"data": nil,
		})
      }


      // Generate Trx ID
      var newTransactionID = helpers.IDGenerator()
      totalPrice := 0

      // Create New Trx
      newTransaction := models.Trx{
            ID: newTransactionID,
            IDUser: user.ID,
            AlamatPengiriman: body.AlamatKirim,
            HargaTotal: totalPrice,
            MethodBayar: body.MetodeBayar,
            KodeInvoice: helpers.InvoiceCodeGenerator(),
      }
      if err := database.DB.Create(&newTransaction).Error; err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                  "status": false,
                  "message": "Failed to POST data",
                  "errors": []string {err.Error()},
                  "data": nil,
            })
      }

      // Produk and DetailTrx Processing
      for i := range len(body.DetailTrx) {
            // Get Produk
            var produk models.Produk
            if err := database.DB.Where("id = ? AND id_toko = ?", body.DetailTrx[i].ProdukID, toko.ID).First(&produk).Error; err != nil {
                  return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                        "status": false,
                        "message": "Failed to POST data",
                        "errors": err.Error(),
                        "data": nil,
                  })
            }

            // Check Produk Stock Availibility
            if produk.Stok < body.DetailTrx[i].Kuantitas {
                  return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                        "status": false,
                        "message": "Failed to POST data",
                        "errors": []string {"Stok Tidak Cukup"},
                        "data": nil,
                  })
            }
            
            // Produk Stock Reduction
            produk.Stok -= body.DetailTrx[i].Kuantitas
            
            // Bind Produk to LogProduk
            newLogProduk := models.LogProduk{
                  ID: helpers.IDGenerator(),
                  IDProduk: produk.ID,            
                  NamaProduk: produk.NamaProduk,       
                  Slug: produk.Slug,              
                  HargaReseller: produk.HargaReseller,     
                  HargaKonsumen: produk.HargaKonsumen,     
                  Deskripsi: produk.Deskripsi,             
                  IDToko: produk.IDToko,              
                  IDCategory: produk.IDCategory,          
            }

            // Create LogProduk
            if err := database.DB.Create(&newLogProduk).Error; err != nil {
                  return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                        "status": false,
                        "message": "Failed to POST data",
                        "errors": []string {err.Error()},
                        "data": nil,
                  })
            }

            // Create DetailTrx
            numHargaProduk, err := strconv.Atoi(produk.HargaKonsumen)
            if err != nil {
                  return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                        "status": false,
                        "message": "Failed to POST data",
                        "errors": []string {err.Error()},
                        "data": nil,
                  })
            }
            totalProductPrice := int(body.DetailTrx[i].Kuantitas * numHargaProduk)
            totalPrice += totalProductPrice
            newTransactionDetail := models.DetailTrx{
                  IDTrx: newTransactionID,
                  IDLogProduk: newLogProduk.ID,
                  IDToko: toko.ID,
                  Kuantitas: body.DetailTrx[i].Kuantitas,
                  HargaTotal: totalProductPrice,
            }
            if err := database.DB.Create(&newTransactionDetail).Error; err != nil {
                  return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                        "status": false,
                        "message": "Failed to POST data",
                        "errors": []string {err.Error()},
                        "data": nil,
                  })
            }

            // Produk Stock Reduction
            if err := database.DB.Save(&produk).Error; err != nil {
                  return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                        "status": false,
                        "message": "Failed to POST data",
                        "errors": []string {err.Error()},
                        "data": nil,
                  })
            }
      }

      // Update Total Price at Trx
      newTransaction.HargaTotal = totalPrice
      if err := database.DB.Save(&newTransaction).Error; err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                  "status": false,
                  "message": "Failed to POST data",
                  "errors": []string {err.Error()},
                  "data": nil,
            })
      }

      // OK Response
      return c.JSON(fiber.Map{
            "status": true,
            "message": "Succeed to POST data",
            "errors": nil,
            "data": 6,
      })
}
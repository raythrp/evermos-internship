package controllers

import (
	"fmt"
	"strconv"
	"time"

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

func ProdukCreate(c *fiber.Ctx) error {
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

	// Toko valid
	var toko models.Toko
	if err := database.DB.Where("id_user = ?", user.ID).First(&toko).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": []string {"No Data Product"},
			"data": nil,
		})
	}

	// Taking Form Values
	namaProduk := c.FormValue("nama_produk")
	categoryID := c.FormValue("category_id")
	hargaReseller := c.FormValue("harga_reseller")
	hargaKonsumen := c.FormValue("harga_konsumen")
	stok := c.FormValue("stok")
	deskripsi := c.FormValue("deskripsi")

	categoryIDNum, err := strconv.ParseUint(categoryID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to POST data",
			"errors": err.Error(),
			"data": nil,
		})
	}

	stokNum, err := strconv.Atoi(stok)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to POST data",
			"errors": err.Error(),
			"data": nil,
		})
	}

	newProduct := models.Produk{
		NamaProduk: namaProduk,
		Slug: helpers.ConvertToSlug(namaProduk),
		IDCategory: uint(categoryIDNum),
		HargaReseller: hargaReseller,
		HargaKonsumen: hargaKonsumen,
		Stok: stokNum,
		Deskripsi: deskripsi,
		IDToko: toko.ID,
	}

	if err := database.DB.Create(&newProduct).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to POST data",
			"errors": err.Error(),
			"data": nil,
		})
	}

	// Getting Created Produk ID
	var existingProduct models.Produk
	if err := database.DB.Where("nama_produk = ? AND deskripsi = ?", newProduct.NamaProduk, newProduct.Deskripsi).First(&existingProduct).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": nil,
			"data": nil,
		})
	}

	// Handle file upload
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to POST data",
			"errors": err.Error(),
			"data": nil,
		})
	}
	files := form.File["photos"]

	// Save the file to the server
	for _, file := range files {
		// Generate unique filename
		filename := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
		savePath := fmt.Sprintf("./uploads/%s", filename)
		if err := c.SaveFile(file, savePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": false,
				"message": "Failed to POST data",
				"errors": nil,
				"data": nil,
			})
		}

		// Create new FotoProduk
		newFotoProduk := models.FotoProduk{
			IDProduk: existingProduct.ID,
			Url: savePath,
		}
		if err := database.DB.Create(&newFotoProduk).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": false,
				"message": "Failed to POST Product Photo",
				"errors": nil,
				"data": nil,
			})
		}
	}

	// OK Response
	return c.JSON(fiber.Map{
		"status": true,
		"message": "Succeed to POST data",
		"errors": nil,
		"data": 4,
	})
}

func ProdukUpdate(c *fiber.Ctx) error {
	noTelp := helpers.JwtClaimer(c)

	// User valid
	var user models.User
	if err := database.DB.Where("notelp = ?", noTelp).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
			"message": "Failed to PUT data",
			"errors": err.Error(),
			"data": nil,
		})
	}

	// Toko valid
	var toko models.Toko
	if err := database.DB.Where("id_user = ?", user.ID).First(&toko).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
			"message": "Failed to PUT data",
			"errors": err.Error(),
			"data": nil,
		})
	}

	productID := c.Params("id")
	var product models.Produk
	if err := database.DB.Where("id = ? AND id_toko = ?", productID, toko.ID).First(&product).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
			"message": "Failed to PUT data",
			"errors": err.Error(),
			"data": nil,
		})
	}

	// Taking Form Values and Binding filled Keys
	namaProduk := c.FormValue("nama_produk", "")
	if namaProduk != "" {
		product.NamaProduk = namaProduk
		product.Slug = helpers.ConvertToSlug(namaProduk)
	}
	categoryID := c.FormValue("category_id", "")
	if categoryID != "" {
		categoryIDNum, err := strconv.ParseUint(categoryID, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": false,
				"message": "Failed to POST data",
				"errors": err.Error(),
				"data": nil,
			})
		}
		product.IDCategory = uint(categoryIDNum)
	}
	hargaReseller := c.FormValue("harga_reseller", "")
	if hargaReseller != "" {
		product.HargaReseller = hargaReseller
	}
	hargaKonsumen := c.FormValue("harga_konsumen", "")
	if hargaKonsumen != "" {
		product.HargaKonsumen = hargaKonsumen
	}
	stok := c.FormValue("stok", "")
	if stok != "" {
		stokNum, err := strconv.Atoi(stok)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": false,
				"message": "Failed to POST data",
				"errors": err.Error(),
				"data": nil,
			})
		}
		product.Stok = stokNum
	}
	deskripsi := c.FormValue("deskripsi", "")
	if deskripsi != "" {
		product.Deskripsi = deskripsi
	}

	// Saving to Table
	if err := database.DB.Save(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to POST data",
			"errors": err.Error(),
			"data": nil,
		})
	}

	// Handle file upload
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to POST data",
			"errors": err.Error(),
			"data": nil,
		})
	}
	files := form.File["photos"]

	// Save the file to the server
	for _, file := range files {
		// Generate unique filename
		filename := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
		savePath := fmt.Sprintf("./uploads/%s", filename)
		if err := c.SaveFile(file, savePath); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": false,
				"message": "Failed to POST data",
				"errors": nil,
				"data": nil,
			})
		}

		// Create new FotoProduk
		newFotoProduk := models.FotoProduk{
			IDProduk: product.ID,
			Url: savePath,
		}
		if err := database.DB.Create(&newFotoProduk).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": false,
				"message": "Failed to POST Product Photo",
				"errors": nil,
				"data": nil,
			})
		}
	}

	// OK Response
	return c.JSON(fiber.Map{
		"status": true,
		"message": "Succeed to POST data",
		"errors": nil,
		"data": "",
	})
}
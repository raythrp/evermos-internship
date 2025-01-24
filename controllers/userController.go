package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raythrp/evermos-internship/database"
	"github.com/raythrp/evermos-internship/helpers"
	models "github.com/raythrp/evermos-internship/models/entities"
	"github.com/raythrp/evermos-internship/models/requests"
	"github.com/raythrp/evermos-internship/models/responses"
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

func UserGetAlamat(c *fiber.Ctx) error {
	noTelp := helpers.JwtClaimer(c)

	// Invalid
	var user models.User
	var alamat []responses.Alamat
	if err := database.DB.Where("notelp = ?", noTelp).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": err.Error(),
			"data": nil,
		})
	}

	if err := database.DB.Where("id_user = ?", user.ID).Find(&alamat).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": err.Error(),
			"data": nil,
		})
	}

	// OK response
	return c.JSON(fiber.Map{
		"status": true,
		"message": "Succeed to GET data",
		"errors": nil,
		"data": alamat,
	})
}

func UserGetAlamatByID(c *fiber.Ctx) error {
	noTelp := helpers.JwtClaimer(c)
	
	// Taking Parameters
	alamatID := c.Params("id")
	if alamatID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to get data",
			"errors": "No Province ID provided",
			"data": nil,
		})
	}

	// Invalid
	var user models.User
	var alamat responses.Alamat
	if err := database.DB.Where("notelp = ?", noTelp).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": err.Error(),
			"data": nil,
		})
	}

	if err := database.DB.Where("id_user = ? AND id = ?", user.ID, alamatID).First(&alamat).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status": false,
			"message": "Failed to GET data",
			"errors": err.Error(),
			"data": nil,
		})
	}

	// OK response
	return c.JSON(fiber.Map{
		"status": true,
		"message": "Succeed to GET data",
		"errors": nil,
		"data": []responses.Alamat {alamat},
	})
}

func UserCreateAlamat(c *fiber.Ctx) error {
	// Taking Body
	var body requests.CreateAlamat
	noTelp := helpers.JwtClaimer(c)

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to Failed data",
			"errors": []string {err.Error()},
			"data": nil,
		})
	}

	// Getting User Details
	var user models.User
	if err := database.DB.Where("notelp = ?", noTelp).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "User not found",
			"errors":  err.Error(),
			"data":    nil,
		})
	}

	// Creating new Alamat
	newAlamat := models.Alamat{
		IDUser: user.ID,
		JudulAlamat: body.JudulAlamat,
		NamaPenerima: body.NamaPenerima,
		NoTelp:  body.NoTelp,
		DetailAlamat: body.DetailAlamat,
	}
	if err := database.DB.Create(&newAlamat).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to POST data",
			"errors": []string {err.Error()},
			"data": nil,
		})
	}

	// OK response
	return c.JSON(fiber.Map{
		"status": true,
		"message": "Succeed to POST data",
		"errors": nil,
		"data": 1,
	})
}

func UserUpdateAlamatByID(c *fiber.Ctx) error {
	// Taking Parameters
	alamatID := c.Params("id")
	if alamatID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to get data",
			"errors": "No Province ID provided",
			"data": nil,
		})
	}

	// Taking body
	var body requests.UpdateAlamat
	noTelp := helpers.JwtClaimer(c)

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to PUT data",
			"errors": []string {"record not found"},
			"data": nil,
		})	
	}

	// Getting User Details
	var user models.User
	if err := database.DB.Where("notelp = ?", noTelp).First(&user).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to PUT data",
			"errors": []string {"record not found"},
			"data": nil,
		})	
	}

	var existingAlamat models.Alamat
	if err := database.DB.Where("id_user = ? and id = ?", user.ID, alamatID).First(&existingAlamat).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to PUT data",
			"errors": []string {"record not found"},
			"data": nil,
		})	
	}

	// Updating table
	existingAlamat.NamaPenerima = body.NamaPenerima
	existingAlamat.NoTelp = body.NoTelp
	existingAlamat.DetailAlamat = body.DetailAlamat
	if err := database.DB.Save(&existingAlamat).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to PUT data",
			"errors": []string {"record not found"},
			"data": nil,
		})
	}

	return c.JSON(fiber.Map{
		"status": true,
		"message": "Succeed to GET data",
		"errors": nil,
		"data": "",
	})
}
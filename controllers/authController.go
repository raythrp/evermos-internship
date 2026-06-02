package controllers

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/raythrp/evermos-internship/database"
	"github.com/raythrp/evermos-internship/helpers"
	models "github.com/raythrp/evermos-internship/models/entities"
	"github.com/raythrp/evermos-internship/models/requests"
)

func AuthLogin(c *fiber.Ctx) error {
	jwtSecret := os.Getenv("JWT_SECRET")
	var body requests.Login

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false,
			"message": "Failed to POST data",
			"errors": []string {"No Telp atau kata sandi salah"},
			"data": nil,
		})
	}

	// Correct Credentials
	var user models.User
	if err := database.DB.Where("notelp = ? AND kata_sandi = ?", body.NoTelp, body.KataSandi).First(&user).Error; err == nil {

		// Token generation
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": user.NoTelp,                        
			"exp": time.Now().Add(time.Hour * 24 * 7).Unix(), // a week
		  })		  

		s, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": false,
				"message": "Failed to POST data",
				"errors": []string {"No Telp atau kata sandi salah"},
				"data": nil,
			})
		}

		// Get Province and City Detail
		province, err := helpers.GetProvinceDetail(user.IDProvinsi)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": false,
				"message": "Failed to POST data",
				"errors": []string {"No Telp atau kata sandi salah"},
				"data": nil,
			})
		}
		city, err := helpers.GetCityDetail(user.IDKota)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": false,
				"message": "Failed to POST data",
				"errors": []string {"No Telp atau kata sandi salah"},
				"data": nil,
			})
		}

		// Response
		return c.JSON(fiber.Map{
			"status": true,
			"message": "Succeed to POST data",
			"errors": nil,
			"data": fiber.Map{
				"nama": user.Nama,
				"no_telp": user.NoTelp,
				"tanggal_Lahir": user.TanggalLahir.Format("02/01/2006"),
				"tentang": user.Tentang,
				"pekerjaan": user.Pekerjaan,
				"email": user.Email,
				"id_provinsi": fiber.Map{
					"id": user.IDProvinsi,
					"name": province.Name,
				},
				"id_kota": fiber.Map{
					"id": user.IDKota,
					"province_id": province.ID,
					"name": city.Name,
				},
				"token": s,
			},
		})
	}

	// Wrong Credentials
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"status": false,
		"message": "Failed to POST data",
		"errors": []string {"No Telp atau kata sandi salah"},
		"data": nil,
	})
}

func AuthRegister(c *fiber.Ctx) error {
	// Taking body
	var body requests.Register

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to POST data",
			"errors": []string {err.Error()},
			"data": nil,
		})
	}


	// Inserting to user table
	formattedTanggalLahir, err := helpers.TimeParserToDate(body.TanggalLahir)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to POST data",
			"errors": []string {err.Error()},
			"data": nil,
		})
	}
	newUser := models.User{
		Nama: body.Nama,
		KataSandi: body.KataSandi,
		NoTelp: body.NoTelp,
		TanggalLahir: formattedTanggalLahir,
		Pekerjaan: body.Pekerjaan,
		Email: body.Email,
		IDProvinsi: body.IDProvinsi,
		IDKota: body.IDKota,
		IsAdmin: false,
	}
	if err := database.DB.Create(&newUser).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to POST data",
			"errors": []string {err.Error()},
			"data": nil,
		})
	}

	var createdUser models.User
	if err := database.DB.Where("notelp = ? AND kata_sandi = ?", body.NoTelp, body.KataSandi).First(&createdUser).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to POST data",
			"errors": []string {err.Error()},
			"data": nil,
		})
	}

	newToko := models.Toko{
		IDUser: createdUser.ID,
		NamaToko: "Toko " + body.Nama,
	}
	if err := database.DB.Create(&newToko).Error; err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to POST data",
			"errors": []string {err.Error()},
			"data": nil,
		})
	}
	
	return c.JSON(fiber.Map{
		"status": true,
		"message": "Succeed to POST data",
		"errors": nil,
		"data": "Register Succeed",
	})
}
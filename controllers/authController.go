package controllers

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/raythrp/evermos-internship/database"
	"github.com/raythrp/evermos-internship/helpers"
	models "github.com/raythrp/evermos-internship/models/entities"
	"github.com/raythrp/evermos-internship/models/requests"
)

const jwtSecret = "WsItpI3Moq4I0rVwo2fOcbvw8CDgJT9FMrsz9zsqAy3e7PRU8sojZ79jSDtnOuO0bjceupoidsajp3u2019eu[20ihceoijlciuab]"

func AuthLogin(ctx *fiber.Ctx) error {
	var body requests.Login

	err := ctx.BodyParser(&body)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status": false,
			"message": "Failed to POST data",
			"errors": []string {"No Telp atau kata sandi salah"},
			"data": nil,
		})
	}

	// Correct Credentials
	var user models.User
	if err := database.DB.Where("notelp = ? AND kata_sandi = ?", body.NoTelp, body.KataSandi).First(&user).Error; err == nil {
		token := jwt.New((jwt.SigningMethodHS256))
		claims := token.Claims.(jwt.MapClaims)
		claims["sub"] = user.NoTelp
		claims["exp"] = time.Now().Add(time.Hour * 24 * 7) // A week

		s, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": false,
				"message": "Failed to POST data",
				"errors": []string {"No Telp atau kata sandi salah"},
				"data": nil,
			})
		}

		// Get Province and City Detail
		var province models.Province = helpers.GetProvinceDetail(user.IDProvinsi)
		var city models.City = helpers.GetCityDetail(user.IDKota)

		// Response
		return ctx.JSON(fiber.Map{
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
	return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"status": false,
		"message": "Failed to POST data",
		"errors": []string {"No Telp atau kata sandi salah"},
		"data": nil,
	})
}

func
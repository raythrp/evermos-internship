package controllers

import (
	"encoding/json"
	"net/http"
	"github.com/gofiber/fiber/v2"
	models "github.com/raythrp/evermos-internship/models/entities"
)

var WilayahListBaseURL = "https://www.emsifa.com/api-wilayah-indonesia/api"
var WilayahDetailBaseURL = "https://emsifa.github.io/api-wilayah-indonesia/api"

func WilayahGetProvinciesList(c *fiber.Ctx) error {
	url := WilayahListBaseURL + "/provinces.json"

	// Fetch the API
	resp, err := http.Get(url)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to get data",
			"errors": err.Error(),
			"data": nil,
		})
	}
	defer resp.Body.Close()

	// Parse the response body to JSON
	var provincies []models.Province
	err = json.NewDecoder(resp.Body).Decode(&provincies)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to get data",
			"errors": err.Error(),
			"data": nil,
		})
	}

	// OK Response
	return c.JSON(fiber.Map{
		"status": true,
		"message": "Succeed to get data",
		"errors": nil,
		"data": provincies,
	})
}

func WilayahGetCitiesList(c *fiber.Ctx) error {
 	// Taking Parameters
	id := c.Params("prov_id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to get data",
			"errors": "No Province ID provided",
			"data": nil,
		})
	}

	// Fetch the API
	resp, err := http.Get(WilayahListBaseURL + "/regencies/" + id + ".json")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to get data",
			"errors": err.Error(),
			"data": nil,
		})
	}
	defer resp.Body.Close()

	// Parse the response body to JSON
	var cities []models.City
	err = json.NewDecoder(resp.Body).Decode(&cities)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to get data",
			"errors": err.Error(),
			"data": nil,
		})
	}

	// OK Response
	return c.JSON(fiber.Map{
		"status": true,
		"message": "Succeed to get data",
		"errors": nil,
		"data": cities,
	})
}

func WilayahGetProvinceDetail(c *fiber.Ctx) error {
 	// Taking Parameters
	id := c.Params("prov_id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to get data",
			"errors": "No Province ID provided",
			"data": nil,
		})
	}

	// Fetch the API
	resp, err := http.Get(WilayahDetailBaseURL + "/province/" + id + ".json")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to get data",
			"errors": err.Error(),
			"data": nil,
		})
	}
	defer resp.Body.Close()

	// Parse the response body to JSON
	var province models.Province
	err = json.NewDecoder(resp.Body).Decode(&province)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to get data",
			"errors": err.Error(),
			"data": nil,
		})
	}

	// OK Response
	return c.JSON(fiber.Map{
		"status": true,
		"message": "Succeed to get data",
		"errors": nil,
		"data": province,
	})
}

func WilayahGetCityDetail(c *fiber.Ctx) error {
 	// Taking Parameters
	id := c.Params("city_id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status": false,
			"message": "Failed to get data",
			"errors": "No Province ID provided",
			"data": nil,
		})
	}

	// Fetch the API
	resp, err := http.Get(WilayahDetailBaseURL + "/regency/" + id + ".json")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to get data",
			"errors": err.Error(),
			"data": nil,
		})
	}
	defer resp.Body.Close()

	// Parse the response body to JSON
	var city models.City
	err = json.NewDecoder(resp.Body).Decode(&city)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status": false,
			"message": "Failed to get data",
			"errors": err.Error(),
			"data": nil,
		})
	}

	// OK Response
	return c.JSON(fiber.Map{
		"status": true,
		"message": "Succeed to get data",
		"errors": nil,
		"data": city,
	})
}
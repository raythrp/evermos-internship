package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raythrp/evermos-internship/controllers"
)

func RouterApp(c *fiber.App) {
	// For Wilayah
	c.Get("/provcity/listprovincies", controllers.WilayahGetProvinciesList)
	c.Get("/provcity/listcities/:prov_id", controllers.WilayahGetCitiesList)
	c.Get("/provcity/detailprovince/:prov_id", controllers.WilayahGetProvinceDetail)
	c.Get("/provcity/detailcity/:city_id", controllers.WilayahGetCityDetail)

	// For Auth
	c.Post("/auth/login", controllers.AuthLogin)
}
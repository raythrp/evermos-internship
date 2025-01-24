package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/raythrp/evermos-internship/controllers"
	"github.com/raythrp/evermos-internship/middlewares"
)

func RouterApp(c *fiber.App) {
	// For Wilayah
	c.Get("/provcity/listprovincies", controllers.WilayahGetProvinciesList)
	c.Get("/provcity/listcities/:prov_id", controllers.WilayahGetCitiesList)
	c.Get("/provcity/detailprovince/:prov_id", controllers.WilayahGetProvinceDetail)
	c.Get("/provcity/detailcity/:city_id", controllers.WilayahGetCityDetail)

	// For Auth
	c.Post("/auth/login", controllers.AuthLogin)
	c.Post("/auth/register", controllers.AuthRegister)

	// For User
	c.Get("/user", middlewares.AuthRequired(), controllers.UserGetMyProfile)
	c.Put("/user", middlewares.AuthRequired(), controllers.UserUpdateMyProfile)
	c.Get("/user/alamat", middlewares.AuthRequired(), controllers.UserGetAlamat)
	c.Get("/user/alamat/:id", middlewares.AuthRequired(), controllers.UserGetAlamatByID)
	c.Post("/user/alamat", middlewares.AuthRequired(), controllers.UserCreateAlamat)
}
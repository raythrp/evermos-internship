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
	c.Put("/user/alamat/:id", middlewares.AuthRequired(), controllers.UserUpdateAlamatByID)
	c.Delete("/user/alamat/:id", middlewares.AuthRequired(), controllers.UserDeleteAlamatByID)

	// For Toko
	c.Get("/toko/my", middlewares.AuthRequired(), controllers.TokoGetMyToko)
	c.Get("/toko/:id_toko", middlewares.AuthRequired(), controllers.TokoGetTokoByID)
	c.Get("/toko", middlewares.AuthRequired(), controllers.TokoGetAllToko)
	c.Put("/toko/:id_toko", middlewares.AuthRequired(), controllers.TokoUpdateProfile)

	// For Category
	c.Get("/category", middlewares.AuthRequired(), controllers.CategoryGetAll)
	c.Get("/category/:id_category", middlewares.AuthRequired(), controllers.CategoryGetByID)
	c.Post("/category", middlewares.AuthRequired(), controllers.CategoryCreate)
	c.Delete("/category/:id", middlewares.AuthRequired(), controllers.CategoryDelete)

	// For Produk
	c.Get("/product", middlewares.AuthRequired(), controllers.ProdukGetAll)
	c.Get("/product/:id", middlewares.AuthRequired(), controllers.ProdukGetByID)
	c.Post("/product", middlewares.AuthRequired(), controllers.ProdukCreate)
	c.Put("/product/:id", middlewares.AuthRequired(), controllers.ProdukUpdate)
	c.Delete("/product/:id", middlewares.AuthRequired(), controllers.ProdukDelete)

	// For Transactions
	c.Get("/trx", middlewares.AuthRequired(), controllers.TrxGetAll)
	c.Get("/trx/:id", middlewares.AuthRequired(), controllers.TrxGetByID)
	c.Post("/trx", middlewares.AuthRequired(), controllers.TrxCreate)
}
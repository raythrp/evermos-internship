//go:build integration

package tests

import (
	"fmt"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/raythrp/evermos-internship/database"
	models "github.com/raythrp/evermos-internship/models/entities"
	"github.com/raythrp/evermos-internship/routers"
	"github.com/raythrp/evermos-internship/testutils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var testApp *fiber.App

func TestMain(m *testing.M) {
	// Load test env file; fall back to environment variables if absent.
	_ = godotenv.Load(".env.test")

	dsn := os.Getenv("TEST_POSTGRES_DSN")
	if dsn == "" {
		fmt.Println("TEST_POSTGRES_DSN not set; skipping integration tests")
		os.Exit(0)
	}

	// Set JWT secret used by the middleware.
	os.Setenv("JWT_SECRET", testutils.TestJWTSecret)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("integration test: cannot connect to postgres: " + err.Error())
	}
	database.DB = db

	// Auto-migrate all entity tables.
	if err := db.AutoMigrate(
		&models.User{},
		&models.Toko{},
		&models.Category{},
		&models.Produk{},
		&models.Alamat{},
		&models.FotoProduk{},
		&models.Trx{},
		&models.DetailTrx{},
		&models.LogProduk{},
	); err != nil {
		panic("integration test: automigrate failed: " + err.Error())
	}

	// Ensure uploads directory exists relative to project root (parent of tests/).
	os.MkdirAll("../uploads", 0755)

	// Seed base data.
	seedData(db)

	// Build the shared Fiber app.
	testApp = fiber.New()
	routers.RouterApp(testApp)

	code := m.Run()

	// Teardown: drop all test tables.
	dropTables(db)

	os.Exit(code)
}

func seedData(db *gorm.DB) {
	// Admin user (created directly; register endpoint doesn't support isAdmin=true).
	admin := models.User{
		Nama:       "Admin",
		KataSandi:  "adminpass",
		NoTelp:     "08000000000",
		Pekerjaan:  "Admin",
		Email:      "admin@test.com",
		IDProvinsi: "11",
		IDKota:     "1101",
		IsAdmin:    true,
	}
	db.Where("notelp = ?", admin.NoTelp).FirstOrCreate(&admin)
	db.Where("id_user = ?", admin.ID).FirstOrCreate(&models.Toko{IDUser: admin.ID, NamaToko: "Toko Admin"})

	// Regular user (seeded so individual tests can log in without running register first).
	regular := models.User{
		Nama:       "Regular User",
		KataSandi:  "userpass",
		NoTelp:     "08111111111",
		Pekerjaan:  "Developer",
		Email:      "user@test.com",
		IDProvinsi: "11",
		IDKota:     "1101",
		IsAdmin:    false,
	}
	db.Where("notelp = ?", regular.NoTelp).FirstOrCreate(&regular)
	db.Where("id_user = ?", regular.ID).FirstOrCreate(&models.Toko{IDUser: regular.ID, NamaToko: "Toko Regular"})

	// Category.
	db.Where("nama_category = ?", "Elektronik").FirstOrCreate(&models.Category{NamaCategory: "Elektronik"})
}

func dropTables(db *gorm.DB) {
	db.Migrator().DropTable(
		&models.LogProduk{},
		&models.DetailTrx{},
		&models.Trx{},
		&models.FotoProduk{},
		&models.Produk{},
		&models.Alamat{},
		&models.Category{},
		&models.Toko{},
		&models.User{},
	)
}

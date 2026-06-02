package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/raythrp/evermos-internship/database"
	"github.com/raythrp/evermos-internship/helpers"
	"github.com/raythrp/evermos-internship/routers"
	"github.com/raythrp/evermos-internship/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/gofiber/fiber/v2"
)

func setupControllerTest(t *testing.T) (*fiber.App, sqlmock.Sqlmock) {
	t.Helper()
	os.Setenv("JWT_SECRET", testutils.TestJWTSecret)
	gormDB, mock, sqlDB := testutils.SetupMockDB()
	database.DB = gormDB
	t.Cleanup(func() { sqlDB.Close() })
	app := fiber.New()
	routers.RouterApp(app)
	return app, mock
}

func TestAuthRegister_BadBody(t *testing.T) {
	app, _ := setupControllerTest(t)

	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAuthRegister_InvalidDate(t *testing.T) {
	app, _ := setupControllerTest(t)

	body := map[string]string{
		"nama":          "Test User",
		"kata_sandi":    "pass123",
		"no_telp":       "08111111111",
		"tanggal_lahir": "not-a-date",
		"pekerjaan":     "Developer",
		"email":         "test@test.com",
		"id_provinsi":   "11",
		"id_kota":       "1101",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestAuthRegister_Success(t *testing.T) {
	app, mock := setupControllerTest(t)

	// INSERT user
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO.*user").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// SELECT user (re-fetch after insert)
	userRows := sqlmock.NewRows([]string{"id", "notelp", "kata_sandi", "nama", "isAdmin", "tanggal_lahir", "created_at", "updated_at"}).
		AddRow(int64(1), "08111111111", "pass123", "Test User", false, time.Now(), time.Now(), time.Now())
	mock.ExpectQuery("SELECT .* FROM .*user.* WHERE").WillReturnRows(userRows)

	// INSERT toko
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO.*toko").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	body := map[string]string{
		"nama":          "Test User",
		"kata_sandi":    "pass123",
		"no_telp":       "08111111111",
		"tanggal_lahir": "01/01/2000",
		"pekerjaan":     "Developer",
		"email":         "test@test.com",
		"id_provinsi":   "11",
		"id_kota":       "1101",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthLogin_WrongCredentials(t *testing.T) {
	app, mock := setupControllerTest(t)

	// SELECT user → no rows
	mock.ExpectQuery("SELECT .* FROM .*user.* WHERE").
		WillReturnRows(sqlmock.NewRows([]string{"id", "notelp"}))

	body := map[string]string{"no_telp": "08111111111", "kata_sandi": "wrongpass"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestAuthLogin_Success(t *testing.T) {
	// Mock province API
	provinceSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"id": "11", "name": "Aceh"})
	}))
	defer provinceSrv.Close()

	// Mock city API
	citySrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"id": "1101", "province_id": "11", "name": "Kab. Simeulue"})
	}))
	defer citySrv.Close()

	origProvince := helpers.ProvinceBaseURL
	origCity := helpers.CityBaseURL
	helpers.ProvinceBaseURL = provinceSrv.URL
	helpers.CityBaseURL = citySrv.URL
	defer func() {
		helpers.ProvinceBaseURL = origProvince
		helpers.CityBaseURL = origCity
	}()

	app, mock := setupControllerTest(t)

	userRows := sqlmock.NewRows([]string{"id", "notelp", "kata_sandi", "nama", "isAdmin", "id_provinsi", "id_kota", "tanggal_lahir", "created_at", "updated_at"}).
		AddRow(int64(1), "08111111111", "pass123", "Test User", false, "11", "1101", time.Now(), time.Now(), time.Now())
	mock.ExpectQuery("SELECT .* FROM .*user.* WHERE").WillReturnRows(userRows)

	body := map[string]string{"no_telp": "08111111111", "kata_sandi": "pass123"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].(map[string]interface{})
	assert.NotEmpty(t, data["token"])
	assert.NoError(t, mock.ExpectationsWereMet())
}

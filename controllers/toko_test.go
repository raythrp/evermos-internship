package controllers_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raythrp/evermos-internship/testutils"
)

func TestGetMyToko_Unauthorized(t *testing.T) {
	app, _ := setupControllerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/toko/my", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestGetMyToko_Success(t *testing.T) {
	app, mock := setupControllerTest(t)
	token := testutils.GenerateTestToken("08111111111")

	userRows := sqlmock.NewRows([]string{"id", "notelp", "isAdmin", "tanggal_lahir", "created_at", "updated_at"}).
		AddRow(int64(1), "08111111111", false, time.Now(), time.Now(), time.Now())
	mock.ExpectQuery("SELECT .* FROM .*user.* WHERE").WillReturnRows(userRows)

	tokoRows := sqlmock.NewRows([]string{"id", "id_user", "nama_toko", "url_foto", "created_at", "updated_at"}).
		AddRow(int64(10), int64(1), "Toko Test", "", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .* FROM .*toko.* WHERE").WillReturnRows(tokoRows)

	req := httptest.NewRequest(http.MethodGet, "/toko/my", nil)
	req.Header.Set("token", token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllToko_NonAdmin(t *testing.T) {
	app, mock := setupControllerTest(t)
	token := testutils.GenerateTestToken("08111111111")

	// Non-admin user
	userRows := sqlmock.NewRows([]string{"id", "notelp", "isAdmin", "tanggal_lahir", "created_at", "updated_at"}).
		AddRow(int64(1), "08111111111", false, time.Now(), time.Now(), time.Now())
	mock.ExpectQuery("SELECT .* FROM .*user.* WHERE").WillReturnRows(userRows)

	req := httptest.NewRequest(http.MethodGet, "/toko", nil)
	req.Header.Set("token", token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAllToko_AdminSuccess(t *testing.T) {
	app, mock := setupControllerTest(t)
	token := testutils.GenerateTestToken("08000000000")

	userRows := sqlmock.NewRows([]string{"id", "notelp", "isAdmin", "tanggal_lahir", "created_at", "updated_at"}).
		AddRow(int64(99), "08000000000", true, time.Now(), time.Now(), time.Now())
	mock.ExpectQuery("SELECT .* FROM .*user.* WHERE").WillReturnRows(userRows)

	tokoRows := sqlmock.NewRows([]string{"id", "nama_toko", "url_foto"}).
		AddRow(int64(1), "Toko A", "").
		AddRow(int64(2), "Toko B", "")
	mock.ExpectQuery("SELECT .* FROM .*toko.*").WillReturnRows(tokoRows)

	req := httptest.NewRequest(http.MethodGet, "/toko?page=1&limit=10", nil)
	req.Header.Set("token", token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTokoByID_NonAdmin(t *testing.T) {
	app, mock := setupControllerTest(t)
	token := testutils.GenerateTestToken("08111111111")

	userRows := sqlmock.NewRows([]string{"id", "notelp", "isAdmin", "tanggal_lahir", "created_at", "updated_at"}).
		AddRow(int64(1), "08111111111", false, time.Now(), time.Now(), time.Now())
	mock.ExpectQuery("SELECT .* FROM .*user.* WHERE").WillReturnRows(userRows)

	req := httptest.NewRequest(http.MethodGet, "/toko/1", nil)
	req.Header.Set("token", token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.NoError(t, mock.ExpectationsWereMet())
}

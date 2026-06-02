package controllers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/raythrp/evermos-internship/testutils"
)

func TestTrxGetAll_NonAdmin(t *testing.T) {
	app, mock := setupControllerTest(t)
	token := testutils.GenerateTestToken("08111111111")

	mock.ExpectQuery("SELECT .* FROM .*user.* WHERE").WillReturnRows(nonAdminUserRows())

	req := httptest.NewRequest(http.MethodGet, "/trx", nil)
	req.Header.Set("token", token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTrxGetAll_Unauthorized(t *testing.T) {
	app, _ := setupControllerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/trx", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestTrxGetByID_Unauthorized(t *testing.T) {
	app, _ := setupControllerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/trx/1", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestTrxGetByID_NotOwned(t *testing.T) {
	app, mock := setupControllerTest(t)
	token := testutils.GenerateTestToken("08111111111")

	// User found
	mock.ExpectQuery("SELECT .* FROM .*user.* WHERE").WillReturnRows(
		sqlmock.NewRows([]string{"id", "notelp", "isAdmin", "tanggal_lahir", "created_at", "updated_at"}).
			AddRow(int64(2), "08111111111", false, time.Now(), time.Now(), time.Now()),
	)

	// Trx query returns no rows (trx belongs to different user)
	mock.ExpectQuery("SELECT .* FROM .*trx.* WHERE").WillReturnRows(
		sqlmock.NewRows([]string{"id", "id_user", "harga_total", "kode_invoice", "method_bayar", "created_at", "updated_at"}),
	)

	req := httptest.NewRequest(http.MethodGet, "/trx/1", nil)
	req.Header.Set("token", token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTrxCreate_Unauthorized(t *testing.T) {
	app, _ := setupControllerTest(t)

	req := httptest.NewRequest(http.MethodPost, "/trx", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestTrxCreate_InsufficientStock(t *testing.T) {
	app, mock := setupControllerTest(t)
	token := testutils.GenerateTestToken("08111111111")

	// User found
	mock.ExpectQuery("SELECT .* FROM .*user.* WHERE").WillReturnRows(
		sqlmock.NewRows([]string{"id", "notelp", "isAdmin", "tanggal_lahir", "created_at", "updated_at"}).
			AddRow(int64(1), "08111111111", false, time.Now(), time.Now(), time.Now()),
	)

	// Toko found
	mock.ExpectQuery("SELECT .* FROM .*toko.* WHERE").WillReturnRows(
		sqlmock.NewRows([]string{"id", "id_user", "nama_toko", "created_at", "updated_at"}).
			AddRow(int64(10), int64(1), "Toko Test", time.Now(), time.Now()),
	)

	// Alamat found and owned
	mock.ExpectQuery("SELECT .* FROM .*alamat.* WHERE").WillReturnRows(
		sqlmock.NewRows([]string{"id", "id_user", "judul_alamat", "created_at", "updated_at"}).
			AddRow(int64(5), int64(1), "Rumah", time.Now(), time.Now()),
	)

	// Create Trx
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO.*trx").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Product found with insufficient stock (stok=1, requested=5)
	mock.ExpectQuery("SELECT .* FROM .*produk.* WHERE").WillReturnRows(
		sqlmock.NewRows([]string{"id", "id_toko", "nama_produk", "harga_konsumen", "stok", "created_at", "updated_at"}).
			AddRow(int64(100), int64(10), "Produk A", "50000", 1, time.Now(), time.Now()),
	)

	body := map[string]interface{}{
		"method_bayar": "transfer",
		"alamat_kirim": 5,
		"detail_trx": []map[string]interface{}{
			{"product_id": 100, "kuantitas": 5},
		},
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/trx", bytes.NewBuffer(b))
	req.Header.Set("token", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.False(t, result["status"].(bool))
	assert.NoError(t, mock.ExpectationsWereMet())
}

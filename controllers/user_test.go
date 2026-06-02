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

func TestGetMyProfile_Unauthorized(t *testing.T) {
	app, _ := setupControllerTest(t)

	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestGetAlamat_Success(t *testing.T) {
	app, mock := setupControllerTest(t)
	token := testutils.GenerateTestToken("08111111111")

	userRows := sqlmock.NewRows([]string{"id", "notelp", "nama", "isAdmin", "tanggal_lahir", "created_at", "updated_at"}).
		AddRow(int64(1), "08111111111", "Test User", false, time.Now(), time.Now(), time.Now())
	mock.ExpectQuery("SELECT .* FROM .*user.* WHERE").WillReturnRows(userRows)

	alamatRows := sqlmock.NewRows([]string{"id", "id_user", "judul_alamat", "nama_penerima", "no_telp", "detail_alamat", "created_at", "updated_at"}).
		AddRow(int64(1), int64(1), "Rumah", "Test User", "08111111111", "Jl. Test No. 1", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .* FROM .*alamat.* WHERE").WillReturnRows(alamatRows)

	req := httptest.NewRequest(http.MethodGet, "/user/alamat", nil)
	req.Header.Set("token", token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetAlamatByID_Success(t *testing.T) {
	app, mock := setupControllerTest(t)
	token := testutils.GenerateTestToken("08111111111")

	userRows := sqlmock.NewRows([]string{"id", "notelp", "nama", "isAdmin", "tanggal_lahir", "created_at", "updated_at"}).
		AddRow(int64(1), "08111111111", "Test User", false, time.Now(), time.Now(), time.Now())
	mock.ExpectQuery("SELECT .* FROM .*user.* WHERE").WillReturnRows(userRows)

	alamatRows := sqlmock.NewRows([]string{"id", "id_user", "judul_alamat", "nama_penerima", "no_telp", "detail_alamat", "created_at", "updated_at"}).
		AddRow(int64(1), int64(1), "Rumah", "Test User", "08111111111", "Jl. Test No. 1", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .* FROM .*alamat.* WHERE").WillReturnRows(alamatRows)

	req := httptest.NewRequest(http.MethodGet, "/user/alamat/1", nil)
	req.Header.Set("token", token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateAlamat_Success(t *testing.T) {
	app, mock := setupControllerTest(t)
	token := testutils.GenerateTestToken("08111111111")

	userRows := sqlmock.NewRows([]string{"id", "notelp", "nama", "isAdmin", "tanggal_lahir", "created_at", "updated_at"}).
		AddRow(int64(1), "08111111111", "Test User", false, time.Now(), time.Now(), time.Now())
	mock.ExpectQuery("SELECT .* FROM .*user.* WHERE").WillReturnRows(userRows)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO.*alamat").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	body := map[string]string{
		"judul_alamat":  "Kantor",
		"nama_penerima": "Test User",
		"no_telp":       "08111111111",
		"detail_alamat": "Jl. Kantor No. 2",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/user/alamat", bytes.NewBuffer(b))
	req.Header.Set("token", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateAlamat_Unauthorized(t *testing.T) {
	app, _ := setupControllerTest(t)

	req := httptest.NewRequest(http.MethodPost, "/user/alamat", nil)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestDeleteAlamat_Success(t *testing.T) {
	app, mock := setupControllerTest(t)
	token := testutils.GenerateTestToken("08111111111")

	userRows := sqlmock.NewRows([]string{"id", "notelp", "isAdmin", "tanggal_lahir", "created_at", "updated_at"}).
		AddRow(int64(1), "08111111111", false, time.Now(), time.Now(), time.Now())
	mock.ExpectQuery("SELECT .* FROM .*user.* WHERE").WillReturnRows(userRows)

	alamatRows := sqlmock.NewRows([]string{"id", "id_user", "created_at", "updated_at"}).
		AddRow(int64(1), int64(1), time.Now(), time.Now())
	mock.ExpectQuery("SELECT .* FROM .*alamat.* WHERE").WillReturnRows(alamatRows)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM.*alamat").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	req := httptest.NewRequest(http.MethodDelete, "/user/alamat/1", nil)
	req.Header.Set("token", token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NoError(t, mock.ExpectationsWereMet())
}

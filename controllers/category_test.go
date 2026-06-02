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

func adminUserRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "notelp", "isAdmin", "tanggal_lahir", "created_at", "updated_at"}).
		AddRow(int64(99), "08000000000", true, time.Now(), time.Now(), time.Now())
}

func nonAdminUserRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "notelp", "isAdmin", "tanggal_lahir", "created_at", "updated_at"}).
		AddRow(int64(1), "08111111111", false, time.Now(), time.Now(), time.Now())
}

func TestCategoryGetAll_NonAdmin(t *testing.T) {
	app, mock := setupControllerTest(t)
	token := testutils.GenerateTestToken("08111111111")

	mock.ExpectQuery("SELECT .* FROM .*user.* WHERE").WillReturnRows(nonAdminUserRows())

	req := httptest.NewRequest(http.MethodGet, "/category", nil)
	req.Header.Set("token", token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryGetAll_AdminSuccess(t *testing.T) {
	app, mock := setupControllerTest(t)
	token := testutils.GenerateTestToken("08000000000")

	mock.ExpectQuery("SELECT .* FROM .*user.* WHERE").WillReturnRows(adminUserRows())

	catRows := sqlmock.NewRows([]string{"id", "nama_category", "created_at", "updated_at"}).
		AddRow(int64(1), "Elektronik", time.Now(), time.Now()).
		AddRow(int64(2), "Fashion", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .* FROM .*category.*").WillReturnRows(catRows)

	req := httptest.NewRequest(http.MethodGet, "/category", nil)
	req.Header.Set("token", token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryGetByID_AdminSuccess(t *testing.T) {
	app, mock := setupControllerTest(t)
	token := testutils.GenerateTestToken("08000000000")

	mock.ExpectQuery("SELECT .* FROM .*user.* WHERE").WillReturnRows(adminUserRows())

	catRows := sqlmock.NewRows([]string{"id", "nama_category", "created_at", "updated_at"}).
		AddRow(int64(1), "Elektronik", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .* FROM .*category.* WHERE").WillReturnRows(catRows)

	req := httptest.NewRequest(http.MethodGet, "/category/1", nil)
	req.Header.Set("token", token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryCreate_AdminSuccess(t *testing.T) {
	app, mock := setupControllerTest(t)
	token := testutils.GenerateTestToken("08000000000")

	mock.ExpectQuery("SELECT .* FROM .*user.* WHERE").WillReturnRows(adminUserRows())

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO.*category").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	body := map[string]string{"nama_category": "Olahraga"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/category", bytes.NewBuffer(b))
	req.Header.Set("token", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryCreate_NonAdmin(t *testing.T) {
	app, mock := setupControllerTest(t)
	token := testutils.GenerateTestToken("08111111111")

	mock.ExpectQuery("SELECT .* FROM .*user.* WHERE").WillReturnRows(nonAdminUserRows())

	body := map[string]string{"nama_category": "Olahraga"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/category", bytes.NewBuffer(b))
	req.Header.Set("token", token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCategoryDelete_AdminSuccess(t *testing.T) {
	app, mock := setupControllerTest(t)
	token := testutils.GenerateTestToken("08000000000")

	mock.ExpectQuery("SELECT .* FROM .*user.* WHERE").WillReturnRows(adminUserRows())

	catRows := sqlmock.NewRows([]string{"id", "nama_category", "created_at", "updated_at"}).
		AddRow(int64(1), "Elektronik", time.Now(), time.Now())
	mock.ExpectQuery("SELECT .* FROM .*category.* WHERE").WillReturnRows(catRows)

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM.*category").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	req := httptest.NewRequest(http.MethodDelete, "/category/1", nil)
	req.Header.Set("token", token)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.NoError(t, mock.ExpectationsWereMet())
}

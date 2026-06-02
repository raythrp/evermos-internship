//go:build integration

package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_Register_Success(t *testing.T) {
	body := map[string]string{
		"nama":          "Regular User",
		"kata_sandi":    "userpass",
		"no_telp":       "08111111111",
		"tanggal_lahir": "15/06/1995",
		"pekerjaan":     "Developer",
		"email":         "user@test.com",
		"id_provinsi":   "11",
		"id_kota":       "1101",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := testApp.Test(req, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.True(t, result["status"].(bool))
}

func TestIntegration_Register_InvalidDate(t *testing.T) {
	body := map[string]string{
		"nama":          "Bad Date User",
		"kata_sandi":    "pass123",
		"no_telp":       "09999999999",
		"tanggal_lahir": "not-a-date",
		"pekerjaan":     "Tester",
		"email":         "baddate@test.com",
		"id_provinsi":   "11",
		"id_kota":       "1101",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := testApp.Test(req, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestIntegration_Login_Success(t *testing.T) {
	body := map[string]string{
		"no_telp":    "08000000000",
		"kata_sandi": "adminpass",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := testApp.Test(req, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.True(t, result["status"].(bool))
	data := result["data"].(map[string]interface{})
	assert.NotEmpty(t, data["token"])
}

func TestIntegration_Login_WrongPassword(t *testing.T) {
	body := map[string]string{
		"no_telp":    "08000000000",
		"kata_sandi": "wrongpassword",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	resp, err := testApp.Test(req, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

// loginAs is a helper that logs in and returns the JWT token.
func loginAs(t *testing.T, noTelp, password string) string {
	t.Helper()
	body := map[string]string{"no_telp": noTelp, "kata_sandi": password}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := testApp.Test(req, 10000)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	return result["data"].(map[string]interface{})["token"].(string)
}

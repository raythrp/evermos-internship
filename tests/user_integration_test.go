//go:build integration

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_GetMyProfile(t *testing.T) {
	token := loginAs(t, "08000000000", "adminpass")

	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	req.Header.Set("token", token)

	resp, err := testApp.Test(req, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.True(t, result["status"].(bool))
	data := result["data"].(map[string]interface{})
	assert.Equal(t, "08000000000", data["no_telp"])
}

func TestIntegration_GetMyProfile_Unauthorized(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user", nil)
	resp, err := testApp.Test(req, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestIntegration_AlamatLifecycle(t *testing.T) {
	token := loginAs(t, "08111111111", "userpass")

	// Create
	body := map[string]string{
		"judul_alamat":  "Rumah",
		"nama_penerima": "Regular User",
		"no_telp":       "08111111111",
		"detail_alamat": "Jl. Merdeka No. 1",
	}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/user/alamat", bytes.NewBuffer(b))
	req.Header.Set("token", token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := testApp.Test(req, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get all alamat
	req2 := httptest.NewRequest(http.MethodGet, "/user/alamat", nil)
	req2.Header.Set("token", token)
	resp2, err := testApp.Test(req2, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	var listResult map[string]interface{}
	json.NewDecoder(resp2.Body).Decode(&listResult)
	alamatList := listResult["data"].([]interface{})
	require.NotEmpty(t, alamatList)

	// Get the first alamat ID
	firstAlamat := alamatList[0].(map[string]interface{})
	alamatID := int(firstAlamat["id"].(float64))

	// Get by ID
	req3 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/user/alamat/%d", alamatID), nil)
	req3.Header.Set("token", token)
	resp3, err := testApp.Test(req3, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp3.StatusCode)

	// Update
	updateBody := map[string]string{
		"nama_penerima": "Updated Name",
		"no_telp":       "08222222222",
		"detail_alamat": "Jl. Baru No. 99",
	}
	ub, _ := json.Marshal(updateBody)
	req4 := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/user/alamat/%d", alamatID), bytes.NewBuffer(ub))
	req4.Header.Set("token", token)
	req4.Header.Set("Content-Type", "application/json")
	resp4, err := testApp.Test(req4, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp4.StatusCode)

	// Delete
	req5 := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/user/alamat/%d", alamatID), nil)
	req5.Header.Set("token", token)
	resp5, err := testApp.Test(req5, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp5.StatusCode)
}

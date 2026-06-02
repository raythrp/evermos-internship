//go:build integration

package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_GetMyToko(t *testing.T) {
	token := loginAs(t, "08000000000", "adminpass")

	req := httptest.NewRequest(http.MethodGet, "/toko/my", nil)
	req.Header.Set("token", token)

	resp, err := testApp.Test(req, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestIntegration_GetAllToko_AdminOnly(t *testing.T) {
	userToken := loginAs(t, "08111111111", "userpass")
	req := httptest.NewRequest(http.MethodGet, "/toko", nil)
	req.Header.Set("token", userToken)
	resp, err := testApp.Test(req, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	adminToken := loginAs(t, "08000000000", "adminpass")
	req2 := httptest.NewRequest(http.MethodGet, "/toko?page=1&limit=10", nil)
	req2.Header.Set("token", adminToken)
	resp2, err := testApp.Test(req2, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)
}

func TestIntegration_GetTokoByID(t *testing.T) {
	adminToken := loginAs(t, "08000000000", "adminpass")

	// Get the toko ID via /toko/my.
	req := httptest.NewRequest(http.MethodGet, "/toko/my", nil)
	req.Header.Set("token", adminToken)
	resp, err := testApp.Test(req, 10000)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var myTokoResult map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&myTokoResult)
	tokoID := int(myTokoResult["data"].(map[string]interface{})["id"].(float64))

	req2 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/toko/%d", tokoID), nil)
	req2.Header.Set("token", adminToken)
	resp2, err := testApp.Test(req2, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)
}

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

func TestIntegration_Category_NonAdminBlocked(t *testing.T) {
	token := loginAs(t, "08111111111", "userpass")

	req := httptest.NewRequest(http.MethodGet, "/category", nil)
	req.Header.Set("token", token)

	resp, err := testApp.Test(req, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestIntegration_CategoryLifecycle(t *testing.T) {
	token := loginAs(t, "08000000000", "adminpass")

	// Create
	body := map[string]string{"nama_category": "Makanan"}
	b, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/category", bytes.NewBuffer(b))
	req.Header.Set("token", token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := testApp.Test(req, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get all — find the new one
	req2 := httptest.NewRequest(http.MethodGet, "/category", nil)
	req2.Header.Set("token", token)
	resp2, err := testApp.Test(req2, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	var listResult map[string]interface{}
	json.NewDecoder(resp2.Body).Decode(&listResult)
	categories := listResult["data"].([]interface{})
	require.NotEmpty(t, categories)

	// Find the "Makanan" category
	var makananID int
	for _, c := range categories {
		cat := c.(map[string]interface{})
		if cat["nama_category"] == "Makanan" {
			makananID = int(cat["id"].(float64))
		}
	}
	require.NotZero(t, makananID, "Makanan category should exist")

	// Get by ID
	req3 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/category/%d", makananID), nil)
	req3.Header.Set("token", token)
	resp3, err := testApp.Test(req3, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp3.StatusCode)

	// Delete
	req4 := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/category/%d", makananID), nil)
	req4.Header.Set("token", token)
	resp4, err := testApp.Test(req4, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp4.StatusCode)
}

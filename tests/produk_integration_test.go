//go:build integration

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// createTestProduct creates a product via the API and returns its ID.
func createTestProduct(t *testing.T, token string, categoryID int) int {
	t.Helper()

	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	w.WriteField("nama_produk", "Test Produk")
	w.WriteField("category_id", fmt.Sprintf("%d", categoryID))
	w.WriteField("harga_reseller", "45000")
	w.WriteField("harga_konsumen", "50000")
	w.WriteField("stok", "10")
	w.WriteField("deskripsi", "Deskripsi produk test")
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/product", &buf)
	req.Header.Set("token", token)
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := testApp.Test(req, 10000)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Retrieve the product via GetAll to find its ID.
	req2 := httptest.NewRequest(http.MethodGet, "/product", nil)
	req2.Header.Set("token", token)
	resp2, err := testApp.Test(req2, 10000)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)

	var listResult map[string]interface{}
	json.NewDecoder(resp2.Body).Decode(&listResult)
	products := listResult["data"].(map[string]interface{})["data"].([]interface{})
	require.NotEmpty(t, products)

	return int(products[0].(map[string]interface{})["id"].(float64))
}

func TestIntegration_ProdukLifecycle(t *testing.T) {
	token := loginAs(t, "08111111111", "userpass")

	// Get an existing category ID.
	adminToken := loginAs(t, "08000000000", "adminpass")
	req := httptest.NewRequest(http.MethodGet, "/category", nil)
	req.Header.Set("token", adminToken)
	resp, err := testApp.Test(req, 10000)
	require.NoError(t, err)
	var catResult map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&catResult)
	cats := catResult["data"].([]interface{})
	require.NotEmpty(t, cats)
	categoryID := int(cats[0].(map[string]interface{})["id"].(float64))

	// Create product.
	produkID := createTestProduct(t, token, categoryID)
	assert.NotZero(t, produkID)

	// Get by ID.
	req2 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/product/%d", produkID), nil)
	req2.Header.Set("token", token)
	resp2, err := testApp.Test(req2, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	// Update (partial).
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	mw.WriteField("stok", "20")
	mw.Close()
	req3 := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/product/%d", produkID), &buf)
	req3.Header.Set("token", token)
	req3.Header.Set("Content-Type", mw.FormDataContentType())
	resp3, err := testApp.Test(req3, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp3.StatusCode)

	// Delete.
	req4 := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/product/%d", produkID), nil)
	req4.Header.Set("token", token)
	resp4, err := testApp.Test(req4, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp4.StatusCode)
}

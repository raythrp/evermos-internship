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

func TestIntegration_TrxGetAll_NonAdmin(t *testing.T) {
	token := loginAs(t, "08111111111", "userpass")
	req := httptest.NewRequest(http.MethodGet, "/trx", nil)
	req.Header.Set("token", token)
	resp, err := testApp.Test(req, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestIntegration_TrxGetAll_Admin(t *testing.T) {
	token := loginAs(t, "08000000000", "adminpass")
	req := httptest.NewRequest(http.MethodGet, "/trx?page=1&limit=10", nil)
	req.Header.Set("token", token)
	resp, err := testApp.Test(req, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestIntegration_TrxCreateAndGet(t *testing.T) {
	userToken := loginAs(t, "08111111111", "userpass")
	adminToken := loginAs(t, "08000000000", "adminpass")

	// Get category.
	catReq := httptest.NewRequest(http.MethodGet, "/category", nil)
	catReq.Header.Set("token", adminToken)
	catResp, _ := testApp.Test(catReq, 10000)
	var catResult map[string]interface{}
	json.NewDecoder(catResp.Body).Decode(&catResult)
	cats := catResult["data"].([]interface{})
	require.NotEmpty(t, cats)
	categoryID := int(cats[0].(map[string]interface{})["id"].(float64))

	// Create a product under the regular user's toko.
	produkID := createTestProduct(t, userToken, categoryID)

	// Create an alamat for the regular user.
	alamatBody := map[string]string{
		"judul_alamat":  "Rumah Trx",
		"nama_penerima": "Regular User",
		"no_telp":       "08111111111",
		"detail_alamat": "Jl. Trx No. 1",
	}
	ab, _ := json.Marshal(alamatBody)
	alamatReq := httptest.NewRequest(http.MethodPost, "/user/alamat", bytes.NewBuffer(ab))
	alamatReq.Header.Set("token", userToken)
	alamatReq.Header.Set("Content-Type", "application/json")
	alamatResp, _ := testApp.Test(alamatReq, 10000)
	require.Equal(t, http.StatusOK, alamatResp.StatusCode)

	// Get alamat ID.
	listAlamatReq := httptest.NewRequest(http.MethodGet, "/user/alamat", nil)
	listAlamatReq.Header.Set("token", userToken)
	listAlamatResp, _ := testApp.Test(listAlamatReq, 10000)
	var alamatListResult map[string]interface{}
	json.NewDecoder(listAlamatResp.Body).Decode(&alamatListResult)
	alamatList := alamatListResult["data"].([]interface{})
	require.NotEmpty(t, alamatList)
	alamatID := int(alamatList[0].(map[string]interface{})["id"].(float64))

	// Create transaction.
	trxBody := map[string]interface{}{
		"method_bayar": "transfer",
		"alamat_kirim": alamatID,
		"detail_trx": []map[string]interface{}{
			{"product_id": produkID, "kuantitas": 2},
		},
	}
	tb, _ := json.Marshal(trxBody)
	trxReq := httptest.NewRequest(http.MethodPost, "/trx", bytes.NewBuffer(tb))
	trxReq.Header.Set("token", userToken)
	trxReq.Header.Set("Content-Type", "application/json")
	trxResp, err := testApp.Test(trxReq, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, trxResp.StatusCode)

	// GetAll as admin should show the transaction.
	allReq := httptest.NewRequest(http.MethodGet, "/trx?page=1&limit=10", nil)
	allReq.Header.Set("token", adminToken)
	allResp, err := testApp.Test(allReq, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, allResp.StatusCode)

	var allResult map[string]interface{}
	json.NewDecoder(allResp.Body).Decode(&allResult)
	data := allResult["data"].(map[string]interface{})
	trxList := data["data"].([]interface{})
	assert.NotEmpty(t, trxList)

	// GetByID for the regular user.
	firstTrx := trxList[0].(map[string]interface{})
	trxID := int(firstTrx["id"].(float64))
	req2 := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/trx/%d", trxID), nil)
	req2.Header.Set("token", userToken)
	resp2, err := testApp.Test(req2, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)
}

func TestIntegration_TrxCreate_InsufficientStock(t *testing.T) {
	userToken := loginAs(t, "08111111111", "userpass")
	adminToken := loginAs(t, "08000000000", "adminpass")

	// Get category.
	catReq := httptest.NewRequest(http.MethodGet, "/category", nil)
	catReq.Header.Set("token", adminToken)
	catResp, _ := testApp.Test(catReq, 10000)
	var catResult map[string]interface{}
	json.NewDecoder(catResp.Body).Decode(&catResult)
	categoryID := int(catResult["data"].([]interface{})[0].(map[string]interface{})["id"].(float64))

	produkID := createTestProduct(t, userToken, categoryID)

	// Get alamat.
	listAlamatReq := httptest.NewRequest(http.MethodGet, "/user/alamat", nil)
	listAlamatReq.Header.Set("token", userToken)
	listAlamatResp, _ := testApp.Test(listAlamatReq, 10000)
	var alamatResult map[string]interface{}
	json.NewDecoder(listAlamatResp.Body).Decode(&alamatResult)
	alamatList := alamatResult["data"].([]interface{})
	require.NotEmpty(t, alamatList)
	alamatID := int(alamatList[0].(map[string]interface{})["id"].(float64))

	// Request more than available stock (stock=10, request=999).
	trxBody := map[string]interface{}{
		"method_bayar": "transfer",
		"alamat_kirim": alamatID,
		"detail_trx": []map[string]interface{}{
			{"product_id": produkID, "kuantitas": 999},
		},
	}
	tb, _ := json.Marshal(trxBody)
	req := httptest.NewRequest(http.MethodPost, "/trx", bytes.NewBuffer(tb))
	req.Header.Set("token", userToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := testApp.Test(req, 10000)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

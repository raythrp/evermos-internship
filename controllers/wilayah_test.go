package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raythrp/evermos-internship/controllers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWilayahGetProvinciesList_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/provinces.json", r.URL.Path)
		json.NewEncoder(w).Encode([]map[string]string{
			{"id": "11", "name": "Aceh"},
			{"id": "12", "name": "Sumatera Utara"},
		})
	}))
	defer srv.Close()

	orig := controllers.WilayahListBaseURL
	controllers.WilayahListBaseURL = srv.URL
	defer func() { controllers.WilayahListBaseURL = orig }()

	app, _ := setupControllerTest(t)
	req := httptest.NewRequest(http.MethodGet, "/provcity/listprovincies", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	assert.True(t, result["status"].(bool))
	data := result["data"].([]interface{})
	assert.Len(t, data, 2)
}

func TestWilayahGetCitiesList_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/regencies/11.json", r.URL.Path)
		json.NewEncoder(w).Encode([]map[string]string{
			{"id": "1101", "province_id": "11", "name": "Kab. Simeulue"},
		})
	}))
	defer srv.Close()

	orig := controllers.WilayahListBaseURL
	controllers.WilayahListBaseURL = srv.URL
	defer func() { controllers.WilayahListBaseURL = orig }()

	app, _ := setupControllerTest(t)
	req := httptest.NewRequest(http.MethodGet, "/provcity/listcities/11", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestWilayahGetCitiesList_MissingProvID(t *testing.T) {
	app, _ := setupControllerTest(t)
	// Route requires :prov_id param; calling without it hits 404 (route not matched)
	req := httptest.NewRequest(http.MethodGet, "/provcity/listcities/", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	// Fiber returns 404 when route param is missing
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestWilayahGetProvinceDetail_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/province/11.json", r.URL.Path)
		json.NewEncoder(w).Encode(map[string]string{"id": "11", "name": "Aceh"})
	}))
	defer srv.Close()

	orig := controllers.WilayahDetailBaseURL
	controllers.WilayahDetailBaseURL = srv.URL
	defer func() { controllers.WilayahDetailBaseURL = orig }()

	app, _ := setupControllerTest(t)
	req := httptest.NewRequest(http.MethodGet, "/provcity/detailprovince/11", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestWilayahGetCityDetail_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/regency/1101.json", r.URL.Path)
		json.NewEncoder(w).Encode(map[string]string{"id": "1101", "province_id": "11", "name": "Kab. Simeulue"})
	}))
	defer srv.Close()

	orig := controllers.WilayahDetailBaseURL
	controllers.WilayahDetailBaseURL = srv.URL
	defer func() { controllers.WilayahDetailBaseURL = orig }()

	app, _ := setupControllerTest(t)
	req := httptest.NewRequest(http.MethodGet, "/provcity/detailcity/1101", nil)

	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

package helpers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/raythrp/evermos-internship/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetCityDetail_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/regency/1101.json", r.URL.Path)
		json.NewEncoder(w).Encode(map[string]string{"id": "1101", "province_id": "11", "name": "Kabupaten Simeulue"})
	}))
	defer srv.Close()

	orig := helpers.CityBaseURL
	helpers.CityBaseURL = srv.URL
	defer func() { helpers.CityBaseURL = orig }()

	city, err := helpers.GetCityDetail("1101")
	require.NoError(t, err)
	assert.Equal(t, "1101", city.ID)
	assert.Equal(t, "11", city.ProvinceID)
	assert.Equal(t, "Kabupaten Simeulue", city.Name)
}

func TestGetCityDetail_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer srv.Close()

	orig := helpers.CityBaseURL
	helpers.CityBaseURL = srv.URL
	defer func() { helpers.CityBaseURL = orig }()

	_, err := helpers.GetCityDetail("1101")
	assert.Error(t, err)
}

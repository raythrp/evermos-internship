package helpers

import (
	"encoding/json"
	"net/http"

	models "github.com/raythrp/evermos-internship/models/entities"
)

func GetCityDetail(id string) models.City {
	var city models.City
	url := "https://emsifa.github.io/api-wilayah-indonesia/api/regency/" + id + ".json"
	provinceResp, err := http.Get(url)
	if err != nil {
		return city
	}
	defer provinceResp.Body.Close()

	// Parse the response body to JSON
	err = json.NewDecoder(provinceResp.Body).Decode(&city)
	if err != nil {
		return city
	}
	return city
}
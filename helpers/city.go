package helpers

import (
	"encoding/json"
	"net/http"

	models "github.com/raythrp/evermos-internship/models/entities"
)

func GetCityDetail(id string) (models.City, error) {
	var city models.City
	url := "https://emsifa.github.io/api-wilayah-indonesia/api/regency/" + id + ".json"
	cityResp, err := http.Get(url)
	if err != nil {
		return city, err
	}
	defer cityResp.Body.Close()

	// Parse the response body to JSON
	err = json.NewDecoder(cityResp.Body).Decode(&city)
	if err != nil {
		return city, err
	}
	return city, nil
}
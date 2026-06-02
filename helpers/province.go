package helpers

import (
	"encoding/json"
	"net/http"

	models "github.com/raythrp/evermos-internship/models/entities"
)

var ProvinceBaseURL = "https://emsifa.github.io/api-wilayah-indonesia/api"

func GetProvinceDetail(id string) (models.Province, error) {
	var province models.Province
	url := ProvinceBaseURL + "/province/" + id + ".json"
	provinceResp, err := http.Get(url)
	if err != nil {
		return province, err
	}
	defer provinceResp.Body.Close()

	// Parse the response body to JSON
	err = json.NewDecoder(provinceResp.Body).Decode(&province)
	if err != nil {
		return province, err
	}
	return province, nil
}
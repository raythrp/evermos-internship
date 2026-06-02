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

func TestGetProvinceDetail_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/province/11.json", r.URL.Path)
		json.NewEncoder(w).Encode(map[string]string{"id": "11", "name": "Aceh"})
	}))
	defer srv.Close()

	orig := helpers.ProvinceBaseURL
	helpers.ProvinceBaseURL = srv.URL
	defer func() { helpers.ProvinceBaseURL = orig }()

	province, err := helpers.GetProvinceDetail("11")
	require.NoError(t, err)
	assert.Equal(t, "11", province.ID)
	assert.Equal(t, "Aceh", province.Name)
}

func TestGetProvinceDetail_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error"))
	}))
	defer srv.Close()

	orig := helpers.ProvinceBaseURL
	helpers.ProvinceBaseURL = srv.URL
	defer func() { helpers.ProvinceBaseURL = orig }()

	// A 500 response with non-JSON body should cause a decode error.
	_, err := helpers.GetProvinceDetail("99")
	assert.Error(t, err)
}

func TestGetProvinceDetail_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer srv.Close()

	orig := helpers.ProvinceBaseURL
	helpers.ProvinceBaseURL = srv.URL
	defer func() { helpers.ProvinceBaseURL = orig }()

	_, err := helpers.GetProvinceDetail("11")
	assert.Error(t, err)
}

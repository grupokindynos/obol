package main

import (
	"encoding/json"
	"github.com/grupokindynos/common/coin-factory"
	"github.com/grupokindynos/obol/config"
	"github.com/grupokindynos/obol/models"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func init() {
	_ = godotenv.Load()
}

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestSimpleRates(t *testing.T) {
	App := GetApp()
	for _, coin := range coinfactory.Coins {
		w := performRequest(App, "GET", "/simple/"+coin.Info.Tag)
		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		value, exists := response["status"]
		assert.Nil(t, err)
		assert.True(t, exists)
		assert.Equal(t, float64(1), value)
		ratesData := response["data"]
		var ratesArray []models.Rate
		ratesBytes, err := json.Marshal(ratesData)
		assert.Nil(t, err)
		err = json.Unmarshal(ratesBytes, &ratesArray)
		assert.Nil(t, err)
		assert.NotZero(t, len(ratesArray))
	}
}

func TestComplexRates(t *testing.T) {
	Coins := coinfactory.Coins
	App := GetApp()
	for _, coin := range Coins {
		w := performRequest(App, "GET", "/complex/POLIS/"+coin.Info.Tag)
		var firstResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &firstResponse)
		if coin.Info.Tag == "POLIS" {
			assert.Equal(t, http.StatusInternalServerError, w.Code)
			assert.Equal(t, config.ErrorNoC2CWithSameCoin.Error(), firstResponse["error"])
			continue
		}
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Nil(t, err)
		firstValue, firstExist := firstResponse["status"]
		assert.True(t, firstExist)
		assert.Equal(t, float64(1), firstValue)
		firstResData := firstResponse["data"]
		assert.NotZero(t, firstResData)

		w2 := performRequest(App, "GET", "/complex/POLIS/"+coin.Info.Tag+"?amount=1")
		assert.Equal(t, http.StatusOK, w2.Code)
		var secondResponse map[string]interface{}
		err = json.Unmarshal(w2.Body.Bytes(), &secondResponse)
		assert.Nil(t, err)
		secondValue, secondExist := firstResponse["status"]
		assert.True(t, secondExist)
		assert.Equal(t, float64(1), secondValue)
		secondResData := secondResponse["data"]
		assert.NotZero(t, secondResData)
		// TODO temp disable
		//assert.Equal(t, firstResData, secondResData)
	}
}

func TestNonExistingRoute(t *testing.T) {
	App := GetApp()
	w := performRequest(App, "GET", "/none")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

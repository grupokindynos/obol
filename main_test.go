package main

import (
	"encoding/json"
	"github.com/grupokindynos/obol/models/coin-factory"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestSimpleRates(t *testing.T) {
	Coins := coinfactory.CoinFactory
	App := GetApp()
	for _, coin := range Coins {
		w := performRequest(App, "GET", "/simple/"+coin.Tag)
		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		value, exists := response["status"]
		assert.Nil(t, err)
		assert.True(t, exists)
		assert.Equal(t, float64(1), value)
		ratesData := response["data"]
		var ratesMap map[string]float64
		ratesBytes, err := json.Marshal(ratesData)
		assert.Nil(t, err)
		err = json.Unmarshal(ratesBytes, &ratesMap)
		assert.Nil(t, err)
		assert.NotZero(t, ratesMap["BTC"])
	}
}

func TestComplexRates(t *testing.T) {
	Coins := coinfactory.CoinFactory
	App := GetApp()
	for _, coin := range Coins {
		w := performRequest(App, "GET", "/complex/polis/"+coin.Tag)
		w2 := performRequest(App, "GET", "/complex/" + coin.Tag + "/polis?amount=1")
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Equal(t, http.StatusOK, w2.Code)
		var firstResponse map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &firstResponse)
		assert.Nil(t, err)
		firstValue, firstExist := firstResponse["status"]
		assert.True(t, firstExist)
		assert.Equal(t, float64(1), firstValue)
		firstResData := firstResponse["data"]
		assert.NotZero(t, firstResData)

		var secondResponse map[string]interface{}
		err = json.Unmarshal(w2.Body.Bytes(), &secondResponse)
		assert.Nil(t, err)
		secondValue, secondExist := firstResponse["status"]
		assert.True(t, secondExist)
		assert.Equal(t, float64(1), secondValue)
		secondResData := firstResponse["data"]
		assert.NotZero(t, secondResData)

		assert.Equal(t, firstResData, secondResData)
	}
}

func TestNonExistingRoute(t *testing.T) {
	App := GetApp()
	w := performRequest(App, "GET", "/none")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

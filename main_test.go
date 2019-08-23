package main

import (
	"encoding/json"
	coin_factory "github.com/grupokindynos/obol/models/coin-factory"
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
	Coins := coin_factory.CoinFactory
	App := GetApp()
	for _, coin := range Coins {
		w := performRequest(App, "GET", "/simple/" + coin.Tag)
		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		value, exists := response["status"]
		assert.Nil(t, err)
		assert.True(t, exists)
		assert.Equal(t, float64(1), value)
	}
}

func TestNonExistingRoute(t *testing.T) {
	App := GetApp()
	w := performRequest(App, "GET", "/none")
	assert.Equal(t, http.StatusNotFound, w.Code)
}

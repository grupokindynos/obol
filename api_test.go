package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
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

func TestPolisRates(t *testing.T) {
	Body := gin.H{"status":  float64(1)}
	App := GetApp()
	w := performRequest(App, "GET", "/simple/polis")
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	value, exists := response["status"]
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, Body["status"], value)
}

func TestBtcRates(t *testing.T) {
	Body := gin.H{"status":  float64(1)}
	App := GetApp()
	w := performRequest(App, "GET", "/simple/btc")
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	value, exists := response["status"]
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, Body["status"], value)
}

func TestDashRates(t *testing.T) {
	Body := gin.H{"status":  float64(1)}
	App := GetApp()
	w := performRequest(App, "GET", "/simple/dash")
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	value, exists := response["status"]
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, Body["status"], value)
}

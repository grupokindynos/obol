package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/obol/config"
	"github.com/grupokindynos/obol/services"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRateController_GetCoinRateFromCoinToCoin(t *testing.T) {
	rateCtrl := loadRateCtrl()
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "fromcoin", Value: "polis"}, gin.Param{Key: "tocoin", Value: "dash"}}
	c.Request, _ = http.NewRequest("GET", "", nil)
	rateCtrl.GetCoinRateFromCoinToCoin(c)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Nil(t, response["error"])
	assert.Equal(t, float64(1), response["status"])
	assert.NotNil(t, response["data"])
}

func TestRateController_GetCoinRateFromCoinToCoinWithAmount(t *testing.T) {
	rateCtrl := loadRateCtrl()
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "fromcoin", Value: "polis"}, gin.Param{Key: "tocoin", Value: "dash"}}
	c.Request, _ = http.NewRequest("GET", "?amount=1000", nil)
	rateCtrl.GetCoinRateFromCoinToCoin(c)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Nil(t, response["error"])
	assert.Equal(t, float64(1), response["status"])
	assert.NotNil(t, response["data"])
}

func TestRateController_GetCoinRateFromCoinToCoinWithAmountInvalid(t *testing.T) {
	rateCtrl := loadRateCtrl()
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "fromcoin", Value: "polis"}, gin.Param{Key: "tocoin", Value: "dash"}}
	c.Request, _ = http.NewRequest("GET", "?amount=test", nil)
	rateCtrl.GetCoinRateFromCoinToCoin(c)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Nil(t, response["data"])
	assert.Equal(t, config.ErrorInvalidAmountOnC2C.Error(), response["error"])
}

func TestRateController_GetCoinRates(t *testing.T) {
	rateCtrl := loadRateCtrl()
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "coin", Value: "polis"}}
	rateCtrl.GetCoinRates(c)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Nil(t, response["error"])
	assert.Equal(t, float64(1), response["status"])
	assert.NotNil(t, response["data"])
}

func loadRateCtrl() *RateController {
	rc := RateController{RateService: services.InitRateService()}
	return &rc
}

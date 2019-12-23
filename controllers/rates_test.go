package controllers

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/obol/config"
	"github.com/grupokindynos/obol/models"
	"github.com/grupokindynos/obol/services"
	"github.com/grupokindynos/obol/services/exchanges/binance"
	"github.com/grupokindynos/obol/services/exchanges/bittrex"
	"github.com/grupokindynos/obol/services/exchanges/crex24"
	"github.com/grupokindynos/obol/services/exchanges/graviex"
	"github.com/grupokindynos/obol/services/exchanges/kucoin"
	"github.com/grupokindynos/obol/services/exchanges/novaexchange"
	"github.com/grupokindynos/obol/services/exchanges/southxhcange"
	"github.com/grupokindynos/obol/services/exchanges/stex"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func init() {
	_ = godotenv.Load("../.env")
}

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

func TestRateController_GetCoinRatesError(t *testing.T) {
	rateCtrl := loadRateCtrl()
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "coin", Value: "non-existing-coin"}}
	rateCtrl.GetCoinRates(c)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Nil(t, response["data"])
	assert.Equal(t, float64(-1), response["status"])
	assert.Equal(t, config.ErrorCoinNotAvailable.Error(), response["error"])
}

func TestRateController_GetCoinRatesFromCoinToCoinErrorFirstCoinInvalid(t *testing.T) {
	rateCtrl := loadRateCtrl()
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "fromcoin", Value: "non-existing-coin"}, gin.Param{Key: "tocoin", Value: "polis"}}
	rateCtrl.GetCoinRateFromCoinToCoin(c)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Nil(t, response["data"])
	assert.Equal(t, float64(-1), response["status"])
	assert.Equal(t, config.ErrorCoinNotAvailable.Error(), response["error"])
}

func TestRateController_GetCoinRatesFromCoinToCoinErrorSecondCoinInvalid(t *testing.T) {
	rateCtrl := loadRateCtrl()
	resp := httptest.NewRecorder()
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(resp)
	c.Params = gin.Params{gin.Param{Key: "fromcoin", Value: "polis"}, gin.Param{Key: "tocoin", Value: "non-existing-coin"}}
	rateCtrl.GetCoinRateFromCoinToCoin(c)
	var response map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Nil(t, err)
	assert.Nil(t, response["data"])
	assert.Equal(t, float64(-1), response["status"])
	assert.Equal(t, config.ErrorCoinNotAvailable.Error(), response["error"])
}

func loadRateCtrl() *RateController {
	rs := &services.RateSevice{
		FiatRates: &models.FiatRates{
			Rates:       nil,
			LastUpdated: time.Time{},
		},
		BittrexService:      bittrex.InitService(),
		BinanceService:      binance.InitService(),
		Crex24Service:       crex24.InitService(),
		StexService:         stex.InitService(),
		SouthXChangeService: southxhcange.InitService(),
		NovaExchangeService: novaexchange.InitService(),
		KuCoinService:       kucoin.InitService(),
		GraviexService:      graviex.InitService(),
	}
	err := rs.LoadFiatRates()
	if err != nil {
		panic(err)
	}
	rc := RateController{RateService: rs, RatesCache: make(map[string]CoinRate)}
	return &rc
}

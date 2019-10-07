package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/coin-factory"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/obol/config"
	"github.com/grupokindynos/obol/models"
	"github.com/grupokindynos/obol/services"
	"math"
	"sort"
	"strconv"
	"time"
)

const RatesCacheTimeFrame = 60 * 15 // 15 minutes

type CoinRate struct {
	LastUpdated int64
	Rates       []models.Rate
}

// RateController is the main type for serving the API routes.
type RateController struct {
	RateService *services.RateSevice
	RatesCache  map[string]CoinRate
}

// GetCoinRates will return a rate map based on the selected coin
func (rc *RateController) GetCoinRates(c *gin.Context) {
	coin := c.Param("coin")
	coinData, err := coinfactory.GetCoin(coin)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	if rc.RatesCache[coinData.Tag].LastUpdated+RatesCacheTimeFrame > time.Now().Unix() {
		responses.GlobalResponseError(rc.RatesCache[coinData.Tag].Rates, err, c)
		return
	}
	rates, err := rc.RateService.GetCoinRates(coinData, false)
	sort.Slice(rates, func(i, j int) bool {
		return rates[i].Code < rates[j].Code
	})
	rc.RatesCache[coinData.Tag] = CoinRate{
		LastUpdated: time.Now().Unix(),
		Rates:       rates,
	}
	responses.GlobalResponseError(rc.RatesCache[coinData.Tag].Rates, err, c)
	return
}

// GetCoinRateFromCoinToCoin will return the rate converting from the first coin to the second coin
// There is also the option of sending the amount trough a query
func (rc *RateController) GetCoinRateFromCoinToCoin(c *gin.Context) {
	fromcoin := c.Param("fromcoin")
	tocoin := c.Param("tocoin")
	fromCoinData, err := coinfactory.GetCoin(fromcoin)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	toCoinData, err := coinfactory.GetCoin(tocoin)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	amount := c.Query("amount")
	wall := c.Query("orders")
	if amount != "" {
		amountNum, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			responses.GlobalResponseError(nil, config.ErrorInvalidAmountOnC2C, c)
			return
		}
		rates, err := rc.RateService.GetCoinToCoinRatesWithAmount(fromCoinData, toCoinData, amountNum, wall)
		responses.GlobalResponseError(rates, nil, c)
		return
	}
	rates, err := rc.RateService.GetCoinToCoinRates(fromCoinData, toCoinData)
	trunk := math.Floor(rates*1e8) / 1e8
	responses.GlobalResponseError(trunk, err, c)
	return
}

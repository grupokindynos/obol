package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/obol/config"
	"github.com/grupokindynos/obol/models/coin-factory"
	"github.com/grupokindynos/obol/services"
	"sort"
	"strconv"
)

// RateController is the main type for serving the API routes.
type RateController struct {
	RateService *services.RateSevice
}

// GetCoinRates will return a rate map based on the selected coin
func (rc *RateController) GetCoinRates(c *gin.Context) {
	coin := c.Param("coin")
	coinData, err := coinfactory.GetCoin(coin)
	if err != nil {
		config.GlobalResponse(nil, err, c)
		return
	}
	rates, err := rc.RateService.GetCoinRates(coinData, false)
	sort.Slice(rates, func(i, j int) bool {
		return rates[i].Code < rates[j].Code
	})
	config.GlobalResponse(rates, err, c)
	return
}

// GetCoinRateFromCoinToCoin will return the rate converting from the first coin to the second coin
// There is also the option of sending the amount trough a query
func (rc *RateController) GetCoinRateFromCoinToCoin(c *gin.Context) {
	fromcoin := c.Param("fromcoin")
	tocoin := c.Param("tocoin")
	fromCoinData, err := coinfactory.GetCoin(fromcoin)
	if err != nil {
		config.GlobalResponse(nil, err, c)
		return
	}
	toCoinData, err := coinfactory.GetCoin(tocoin)
	if err != nil {
		config.GlobalResponse(nil, err, c)
		return
	}
	amount := c.Query("amount")
	if amount != "" {
		amountNum, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			config.GlobalResponse(nil, config.ErrorInvalidAmountOnC2C, c)
			return
		}
		rates, err := rc.RateService.GetCoinToCoinRatesWithAmount(fromCoinData, toCoinData, amountNum)
		config.GlobalResponse(rates, err, c)
		return
	}
	rates, err := rc.RateService.GetCoinToCoinRates(fromCoinData, toCoinData)
	config.GlobalResponse(rates, err, c)
	return
}

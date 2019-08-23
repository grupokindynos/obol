package main

import (
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/obol/config"
	coin_factory "github.com/grupokindynos/obol/models/coin-factory"
	"github.com/grupokindynos/obol/services"
	"strconv"
)

type RateController struct {
	RateService *services.RateSevice
}

func (rc *RateController) GetCoinRates(c *gin.Context) {
	coin := c.Param("coin")
	coinData, err := coin_factory.GetCoin(coin)
	if err != nil {
		config.GlobalResponse(nil, err, c)
		return
	}
	rates, err := rc.RateService.GetCoinRates(coinData)
	config.GlobalResponse(rates, err, c)
	return
}

func (rc *RateController) GetCoinRateFromCoinToCoin(c *gin.Context) {
	fromcoin := c.Param("fromcoin")
	tocoin := c.Param("tocoin")
	fromCoinData, err := coin_factory.GetCoin(fromcoin)
	if err != nil {
		config.GlobalResponse(nil, err, c)
		return
	}
	toCoinData, err := coin_factory.GetCoin(tocoin)
	if err != nil {
		config.GlobalResponse(nil, err, c)
		return
	}
	query := c.Request.URL.Query()
	amount, reqAmount := query["amount"]
	if reqAmount {
		amountNum, err := strconv.ParseFloat(amount[0], 64)
		if err != nil {
			config.GlobalResponse(nil, config.ErrorUnableToParseStringToFloat, c)
			return
		}
		rates, err := rc.RateService.GetCoinToCoinRatesWithAmount(fromCoinData, toCoinData, amountNum)
		config.GlobalResponse(rates, err, c)
		return
	} else {
		rates, err := rc.RateService.GetCoinToCoinRates(fromCoinData, toCoinData)
		config.GlobalResponse(rates, err, c)
		return
	}
}

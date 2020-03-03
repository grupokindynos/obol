package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/coin-factory"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/obol/config"
	"github.com/grupokindynos/obol/models"
	"github.com/grupokindynos/obol/services"
	"errors"
	"strconv"
	"strings"
	"time"
)

const RatesCacheTimeFrame = 60 * 15 // 15 minutes

type CoinRate struct {
	LastUpdated int64
	Rates       map[string] models.RateV2
}

// RateController is the main type for serving the API routes.
type RateController struct {
	RateService *services.RateSevice
	RatesCache  map[string]CoinRate
}

// Returns rate from crypto to specific FIAT currency
func (rc *RateController) GetCoinToFIATRate(c *gin.Context) {
	fromcoin := c.Param("fromcoin")
	tocoin := c.Param("tocoin")
	tocoin = strings.ToUpper(tocoin)
	fromCoinData, err := coinfactory.GetCoin(fromcoin)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	if rc.RatesCache[fromCoinData.Info.Tag].LastUpdated+RatesCacheTimeFrame > time.Now().Unix() {
		if rate, ok := rc.RatesCache[fromCoinData.Info.Tag].Rates[tocoin]; ok {
			responses.GlobalResponseError(rate, err, c)
			return
		}
		responses.GlobalResponseError(nil, errors.New("FIAT coin not found"), c)
		return
	}
	rates, err := rc.RateService.GetCoinRates(fromCoinData, false)
	rc.RatesCache[fromCoinData.Info.Tag] = CoinRate{
		LastUpdated: time.Now().Unix(),
		Rates:       rates,
	}
	if rate, ok := rc.RatesCache[fromCoinData.Info.Tag].Rates[tocoin]; ok {
		responses.GlobalResponseError(rate, err, c)
		return
	}
	responses.GlobalResponseError(nil, errors.New("FIAT coin not found"), c)
	return
}

// GetCoinRates will return a rate map based on the selected coin
func (rc *RateController) GetCoinRatesV2(c *gin.Context) {
	coin := c.Param("coin")
	coinData, err := coinfactory.GetCoin(coin)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	if rc.RatesCache[coinData.Info.Tag].LastUpdated+RatesCacheTimeFrame > time.Now().Unix() {
		responses.GlobalResponseError(rc.RatesCache[coinData.Info.Tag].Rates, err, c)
		return
	}
	rates, err := rc.RateService.GetCoinRates(coinData, false)
	rc.RatesCache[coinData.Info.Tag] = CoinRate{
		LastUpdated: time.Now().Unix(),
		Rates:       rates,
	}
	responses.GlobalResponseError(rc.RatesCache[coinData.Info.Tag].Rates, err, c)
	return
}

// GetCoinRates will return a rate map based on the selected coin
func (rc *RateController) GetCoinRates(c *gin.Context) {
	coin := c.Param("coin")
	coinData, err := coinfactory.GetCoin(coin)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	if rc.RatesCache[coinData.Info.Tag].LastUpdated+RatesCacheTimeFrame > time.Now().Unix() {
		ratesV1 := convertToV1Array(rc.RatesCache[coinData.Info.Tag].Rates)
		responses.GlobalResponseError(ratesV1, err, c)
		return
	}
	rates, err := rc.RateService.GetCoinRates(coinData, false)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	ratesV1 := convertToV1Array(rates)
	rc.RatesCache[coinData.Info.Tag] = CoinRate{
		LastUpdated: time.Now().Unix(),
		Rates:       rates,
	}
	responses.GlobalResponseError(ratesV1, err, c)
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
	amountReq := c.Query("amount")
	if amountReq != "" {
		amountNum, err := strconv.ParseFloat(amountReq, 64)
		if err != nil {
			responses.GlobalResponseError(nil, config.ErrorInvalidAmountOnC2C, c)
			return
		}
		rates, err := rc.RateService.GetCoinToCoinRatesWithAmount(fromCoinData, toCoinData, amountNum)
		responses.GlobalResponseError(rates, err, c)
		return
	}
	rates, err := rc.RateService.GetCoinToCoinRates(fromCoinData, toCoinData)
	responses.GlobalResponseError(rates, err, c)
	return
}

// GetCoinLiquidity will return the liquidity available for the selected coin.
func (rc *RateController) GetCoinLiquidity(c *gin.Context) {
	coin := c.Param("coin")
	coinConfig, err := coinfactory.GetCoin(coin)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	liquidity, err := rc.RateService.GetCoinLiquidity(coinConfig)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	responses.GlobalResponseError(liquidity, err, c)
	return
}

func convertToV1Array(ratesMap map[string]models.RateV2) (rates []models.Rate){
	for code, rate := range ratesMap {
		newRate := models.Rate{
			Code: code,
			Name: rate.Name,
			Rate: rate.Rate,
		}
		rates = append(rates, newRate)
	}
	return
}

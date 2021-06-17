package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/grupokindynos/common/coin-factory"
	"github.com/grupokindynos/common/responses"
	"github.com/grupokindynos/obol/config"
	"github.com/grupokindynos/obol/models"
	"github.com/grupokindynos/obol/services"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const RatesCacheTimeFrame = 60 * 15 // 15 minutes

type CoinRate struct {
	LastUpdated int64
	Rates       map[string]models.RateV2
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
	exchange := c.Query("exchange")
	fromCoinData, err := coinfactory.GetCoin(fromcoin)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	if rc.RatesCache[fromCoinData.Info.Tag].LastUpdated+RatesCacheTimeFrame > time.Now().Unix() {
		if rate, ok := rc.RatesCache[fromCoinData.Info.Tag].Rates[tocoin]; ok {
			responses.GlobalResponseError(rate.Rate, err, c)
			return
		}
		responses.GlobalResponseError(nil, errors.New("FIAT coin not found"), c)
		return
	}
	rates, err := rc.RateService.GetCoinRates(fromCoinData, exchange,false, false)
	rc.RatesCache[fromCoinData.Info.Tag] = CoinRate{
		LastUpdated: time.Now().Unix(),
		Rates:       rates,
	}
	if rate, ok := rc.RatesCache[fromCoinData.Info.Tag].Rates[tocoin]; ok {
		responses.GlobalResponseError(rate.Rate, err, c)
		return
	}
	responses.GlobalResponseError(nil, errors.New("FIAT coin not found"), c)
	return
}

// GetCoinRatesV2 GetCoinRates will return a rate map based on the selected coin
func (rc *RateController) GetCoinRatesV2(c *gin.Context) {
	coin := c.Param("coin")
	if coin == "POLISBSC" {
		coin = "POLIS"
	}
	coinData, err := coinfactory.GetCoin(coin)
	exchange := c.Query("exchange")
	log.Println("exchange debug: ", exchange)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	if rc.RatesCache[coinData.Info.Tag].LastUpdated+RatesCacheTimeFrame > time.Now().Unix() {
		responses.GlobalResponseError(rc.RatesCache[coinData.Info.Tag].Rates, err, c)
		return
	}
	rates, err := rc.RateService.GetCoinRates(coinData, exchange, false, true)
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
	//if coin == "POLISBSC" || coin == "polisbsc" {
	//	coin = "POLIS"
	//}
	exchange := c.Query("exchange")
	coinData, err := coinfactory.GetCoin(coin)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	if rc.RatesCache[coinData.Info.Tag].LastUpdated+RatesCacheTimeFrame > time.Now().Unix() {
		ratesV1 := convertToV1Array(rc.RatesCache[coinData.Info.Tag].Rates)
		sort.Slice(ratesV1, func(i, j int) bool {
			return ratesV1[i].Code < ratesV1[j].Code
		})
		responses.GlobalResponseError(ratesV1, err, c)
		return
	}
	rates, err := rc.RateService.GetCoinRates(coinData, exchange,false, true)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	ratesV1 := convertToV1Array(rates)
	sort.Slice(ratesV1, func(i, j int) bool {
		return ratesV1[i].Code < ratesV1[j].Code
	})
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
	exchange := c.Query("exchange")
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
		rates, err := rc.RateService.GetCoinToCoinRatesWithAmount(fromCoinData, toCoinData, amountNum, exchange)
		responses.GlobalResponseError(rates, err, c)
		return
	}
	rates, err := rc.RateService.GetCoinToCoinRates(fromCoinData, toCoinData, exchange)
	responses.GlobalResponseError(rates, err, c)
	return
}

func (rc *RateController) GetCoinToCoinRateWithExchangeMargin(c *gin.Context) {
	fromCoin := c.Param("fromCoin")
	toCoin := c.Param("toCoin")
	exchange := c.Query("exchange")
	amountStr := c.Query("amount")

	if amountStr == "" {
		responses.GlobalResponseError(nil, errors.New("missing parameter amount"), c)
		return
	}
	if exchange == "" {
		responses.GlobalResponseError(nil, errors.New("missing parameter exchange"), c)
		return
	}

	fromCoinInfo, err := coinfactory.GetCoin(fromCoin)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	toCoinInfo, err := coinfactory.GetCoin(toCoin)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	rateMargin, err := getExchangeRateMargin(exchange)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	amount *= rateMargin
	rate, err := rc.RateService.GetCoinToCoinRatesWithAmount(fromCoinInfo, toCoinInfo, amount, exchange)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}

	rate.Amount /= rateMargin
	responses.GlobalResponse(rate, c)
	return
}


// GetCoinLiquidity will return the liquidity available for the selected coin.
func (rc *RateController) GetCoinLiquidity(c *gin.Context) {
	coin := c.Param("coin")
	exchange := c.Query("exchange")
	coinConfig, err := coinfactory.GetCoin(coin)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	liquidity, err := rc.RateService.GetCoinLiquidity(coinConfig, exchange)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	responses.GlobalResponseError(liquidity, err, c)
	return
}

func (rc *RateController) GetNodeProvider(c *gin.Context) {
	coin := c.Param("coin")
	coinConfig, err := coinfactory.GetCoin(coin)
	if err != nil {
		responses.GlobalResponseError(nil, err, c)
		return
	}
	if coinConfig.Info.Token {
		coinConfig.Info.Tag = "ETH"
	}
	provider := os.Getenv(coinConfig.Info.Tag + "_NODE")
	responses.GlobalResponseError(provider, err, c)
	return
}

func convertToV1Array(ratesMap map[string]models.RateV2) (rates []models.Rate) {
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

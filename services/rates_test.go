package services

import (
	coin_factory "github.com/grupokindynos/obol/models/coin-factory"
	"github.com/stretchr/testify/assert"
	"testing"
)

var rateService = InitRateService()

func TestInitRateService(t *testing.T) {
	assert.NotNil(t, rateService.FiatRates)
	assert.NotNil(t, rateService.BinanceService)
	assert.NotNil(t, rateService.Crex24Service)
	assert.NotNil(t, rateService.CryptoBridgeService)
	assert.NotNil(t, rateService.StexService)
	assert.IsType(t, &RateSevice{}, rateService)
}

func TestRateSevice_GetBtcRates(t *testing.T) {
	btcRate, err := rateService.GetBtcRates()
	assert.Nil(t, err)
	assert.NotZero(t, len(btcRate))
}

func TestRateSevice_GetBtcMxnRate(t *testing.T) {
	rate, err := rateService.GetBtcMxnRate()
	assert.Nil(t, err)
	assert.NotNil(t, rate)
}

func TestRateSevice_GetCoinOrdersWall(t *testing.T) {
	for _, coin := range coin_factory.CoinFactory {
		if coin.Tag == "BTC" { continue }
		orders, err := rateService.GetCoinOrdersWall(&coin)
		assert.Nil(t, err)
		assert.NotZero(t, len(orders))
	}
}

func TestRateSevice_GetCoinRates(t *testing.T) {
	for _, coin := range coin_factory.CoinFactory {
		rates, err := rateService.GetCoinRates(&coin)
		assert.Nil(t, err)
		assert.NotZero(t, len(rates))
	}
}

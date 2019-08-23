package services

import (
	"github.com/grupokindynos/obol/config"
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

func TestRateSevice_GetCoinToCoinRatesWithAmount(t *testing.T) {
	polis := &coin_factory.Polis
	dash := &coin_factory.Dash
	rate, err := rateService.GetCoinToCoinRatesWithAmount(polis, dash, 100)
	assert.Nil(t, err)
	assert.NotZero(t, rate)
}

func TestRateSevice_GetCoinToCoinRates(t *testing.T) {
	polis := &coin_factory.Polis
	dash := &coin_factory.Dash
	rate, err := rateService.GetCoinToCoinRates(polis, dash)
	assert.Nil(t, err)
	assert.NotZero(t, rate)
}

func TestRateSevice_GetCoinToCoinRatesWithAmountSameParams(t *testing.T) {
	polis := &coin_factory.Polis
	rate, err := rateService.GetCoinToCoinRatesWithAmount(polis, polis, 100)
	assert.Zero(t, rate)
	assert.Equal(t, config.ErrorNoC2CWithSameCoin, err)
}

func TestRateSevice_GetCoinToCoinRatesSameParams(t *testing.T) {
	polis := &coin_factory.Polis
	rate, err := rateService.GetCoinToCoinRates(polis, polis)
	assert.Zero(t, rate)
	assert.Equal(t, config.ErrorNoC2CWithSameCoin, err)
}

func TestRateSevice_GetCoinToCoinRatesWithAmountBTC(t *testing.T) {
	polis := &coin_factory.Polis
	btc := &coin_factory.Bitcoin
	rate, err := rateService.GetCoinToCoinRatesWithAmount(btc, polis, 100)
	assert.Zero(t, rate)
	assert.Equal(t, config.ErrorNoC2CWithBTC, err)
}

func TestRateSevice_GetCoinToCoinRatesBTC(t *testing.T) {
	polis := &coin_factory.Polis
	btc := &coin_factory.Bitcoin
	rate, err := rateService.GetCoinToCoinRates(btc, polis)
	assert.Zero(t, rate)
	assert.Equal(t, config.ErrorNoC2CWithBTC, err)
}

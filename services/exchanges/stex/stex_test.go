package stex

import (
	"github.com/grupokindynos/obol/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

var service = InitService()

func TestService_CoinRate(t *testing.T) {
	rate, err := service.CoinRate("xsg")
	assert.Nil(t, err)
	assert.NotNil(t, rate)
	assert.NotZero(t, rate)
}

func TestService_CoinRateError(t *testing.T) {
	rate, err := service.CoinRate("non-existing")
	assert.Zero(t, rate)
	assert.Equal(t, config.ErrorUnknownIdForCoin, err)
}

func TestService_CoinMarketOrders(t *testing.T) {
	orders, err := service.CoinMarketOrders("xsg")
	assert.Nil(t, err)
	assert.NotNil(t, orders)
	assert.NotZero(t, len(orders))
}

func TestService_CoinMarketOrdersError(t *testing.T) {
	rate, err := service.CoinMarketOrders("non-existing")
	assert.Zero(t, rate)
	assert.Equal(t, config.ErrorUnknownIdForCoin, err)
}

func TestInitService(t *testing.T) {
	assert.NotNil(t, service.BaseRateURL)
	assert.NotNil(t, service.MarketRateURL)
}

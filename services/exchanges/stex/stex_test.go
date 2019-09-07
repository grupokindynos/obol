package stex

import (
	"github.com/grupokindynos/obol/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

var service = InitService()

func TestService_CoinMarketOrders(t *testing.T) {
	orders, err := service.CoinMarketOrders("xsg")
	if err != nil {
		assert.Zero(t, len(orders))
		assert.NotNil(t, err)
	} else {
		assert.NotNil(t, orders)
		assert.NotZero(t, len(orders))
	}
}

func TestService_CoinMarketOrdersError(t *testing.T) {
	rate, err := service.CoinMarketOrders("non-existing")
	assert.Zero(t, len(rate))
	assert.Equal(t, config.ErrorUnknownIdForCoin, err)
}

func TestInitService(t *testing.T) {
	assert.NotNil(t, service.BaseRateURL)
	assert.NotNil(t, service.MarketRateURL)
}

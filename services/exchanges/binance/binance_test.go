package binance

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var service = InitService()

func TestService_CoinRate(t *testing.T) {
	rate, err := service.CoinRate("grs")
	assert.Nil(t, err)
	assert.NotNil(t, rate)
	assert.NotZero(t, rate)
}

func TestService_CoinMarketOrders(t *testing.T) {
	orders, err := service.CoinMarketOrders("grs")
	assert.Nil(t, err)
	assert.NotNil(t, orders)
	assert.NotZero(t, len(orders))
}

func TestInitService(t *testing.T) {
	assert.NotNil(t, service.BaseRateURL)
	assert.NotNil(t, service.MarketRateURL)
}
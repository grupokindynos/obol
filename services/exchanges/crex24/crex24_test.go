package crex24

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var service = InitService()

func TestService_CoinMarketOrders(t *testing.T) {
	orders, err := service.CoinMarketOrders("mnp")
	if err != nil {
		assert.Nil(t, len(orders))
		assert.NotNil(t, err)
	} else {
		assert.NotNil(t, orders)
		assert.NotZero(t, len(orders))
	}
}

func TestInitService(t *testing.T) {
	assert.NotNil(t, service.MarketRateURL)
}

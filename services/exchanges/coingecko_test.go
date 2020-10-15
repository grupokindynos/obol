package exchanges

import (
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

var coingecko = NewCoinGecko()

func TestService_CoinMarketOrders(t *testing.T) {
	orders, err := coingecko.GetSimplePriceToBtcAsRate("POLIS")
	assert.Nil(t, err)
	floatPrice, _ := orders["buy"][0].Price.Float64()
	assert.GreaterOrEqual(t, floatPrice, float64(0))
	log.Println(floatPrice)
}
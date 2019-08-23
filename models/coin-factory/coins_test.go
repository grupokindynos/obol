package coinfactory

import (
	"github.com/grupokindynos/obol/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCoinFactory(t *testing.T) {
	for _, coin := range CoinFactory {
		newCoin, err := GetCoin(coin.Tag)
		assert.Nil(t, err)
		assert.IsType(t, &Coin{}, newCoin)
	}
}

func TestNoCoin(t *testing.T) {
	newCoin, err := GetCoin("NONEXISTINGCOIN")
	assert.Equal(t, config.ErrorCoinNotAvailable, err)
	assert.Nil(t, newCoin)
}

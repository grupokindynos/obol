package services

import (
	"github.com/grupokindynos/common/coin-factory"
	"github.com/grupokindynos/common/coin-factory/coins"
	"github.com/grupokindynos/obol/config"
	"github.com/grupokindynos/obol/models"
	"github.com/grupokindynos/obol/services/exchanges/binance"
	"github.com/grupokindynos/obol/services/exchanges/bitrue"
	"github.com/grupokindynos/obol/services/exchanges/bittrex"
	"github.com/grupokindynos/obol/services/exchanges/crex24"
	"github.com/grupokindynos/obol/services/exchanges/graviex"
	"github.com/grupokindynos/obol/services/exchanges/kucoin"
	"github.com/grupokindynos/obol/services/exchanges/novaexchange"
	"github.com/grupokindynos/obol/services/exchanges/southxhcange"
	"github.com/grupokindynos/obol/services/exchanges/stex"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func init() {
	_ = godotenv.Load("../.env")
}

func TestInitRateService(t *testing.T) {
	rateService := loadRateService()
	assert.NotNil(t, rateService.FiatRates)
	assert.NotNil(t, rateService.BinanceService)
	assert.NotNil(t, rateService.Crex24Service)
	assert.NotNil(t, rateService.StexService)
	assert.IsType(t, &RateSevice{}, rateService)
}

func TestRateSevice_GetBtcRates(t *testing.T) {
	rateService := loadRateService()
	btcRate, err := rateService.GetBtcRates()
	assert.Nil(t, err)
	assert.NotZero(t, len(btcRate))
}

func TestRateSevice_GetBtcMxnRate(t *testing.T) {
	rateService := loadRateService()
	rate, err := rateService.GetBtcMxnRate()
	assert.Nil(t, err)
	assert.NotNil(t, rate)
}

func TestRateSevice_GetCoinOrdersWall(t *testing.T) {
	rateService := loadRateService()
	for _, coin := range coinfactory.Coins {
		if coin.Info.Tag == "BTC" {
			continue
		}
		orders, err := rateService.GetCoinOrdersWall(coin)
		// TODO handle possible error types
		if err != nil {
			assert.NotNil(t, err)
			assert.Zero(t, len(orders))
		} else {
			assert.Nil(t, err)
			assert.NotZero(t, len(orders))
		}
	}
}

func TestRateSevice_GetCoinRates(t *testing.T) {
	rateService := loadRateService()
	for _, coin := range coinfactory.Coins {
		rates, err := rateService.GetCoinRates(coin, false)
		// TODO handle possible error types
		if err != nil {
			assert.NotNil(t, err)
			assert.Zero(t, len(rates))
		} else {
			assert.Nil(t, err)
			assert.NotZero(t, len(rates))
		}
	}
}

func TestRateSevice_GetCoinToCoinRatesWithAmount(t *testing.T) {
	rateService := loadRateService()
	polis, _ := coinfactory.GetCoin("polis")
	dash, _ := coinfactory.GetCoin("dash")
	rate, err := rateService.GetCoinToCoinRatesWithAmount(polis, dash, 100, "buy")
	assert.Nil(t, err)
	assert.NotZero(t, rate)
}

func TestRateSevice_GetCoinToCoinRatesBTCFirst(t *testing.T) {
	rateService := loadRateService()
	polis, _ := coinfactory.GetCoin("polis")
	btc, _ := coinfactory.GetCoin("btc")
	rate, err := rateService.GetCoinToCoinRates(btc, polis)
	assert.Nil(t, err)
	assert.NotZero(t, rate)
}

func TestRateSevice_GetCoinToCoinRatesBTCSecond(t *testing.T) {
	rateService := loadRateService()
	polis, _ := coinfactory.GetCoin("polis")
	btc, _ := coinfactory.GetCoin("btc")
	rate, err := rateService.GetCoinToCoinRates(polis, btc)
	assert.Nil(t, err)
	assert.NotZero(t, rate)
}

func TestRateSevice_GetCoinToCoinRates(t *testing.T) {
	rateService := loadRateService()
	polis, _ := coinfactory.GetCoin("polis")
	dash, _ := coinfactory.GetCoin("dash")
	rate, err := rateService.GetCoinToCoinRates(polis, dash)
	assert.Nil(t, err)
	assert.NotZero(t, rate)
}

func TestRateSevice_GetCoinToCoinRatesWithAmountSameParams(t *testing.T) {
	rateService := loadRateService()
	polis, _ := coinfactory.GetCoin("polis")
	rate, err := rateService.GetCoinToCoinRatesWithAmount(polis, polis, 100, "buy")
	assert.Zero(t, rate)
	assert.Equal(t, config.ErrorNoC2CWithSameCoin, err)
}

func TestRateSevice_GetCoinToCoinRatesSameParams(t *testing.T) {
	rateService := loadRateService()
	polis, _ := coinfactory.GetCoin("polis")
	rate, err := rateService.GetCoinToCoinRates(polis, polis)
	assert.Zero(t, rate)
	assert.Equal(t, config.ErrorNoC2CWithSameCoin, err)
}

func TestNoServiceForCoin(t *testing.T) {
	rateService := loadRateService()
	mockCoin := &coins.Coin{Info: coins.CoinInfo{Tag: "FaKeCOIN", Name: "FakeCoin"}}
	_, err := rateService.GetCoinOrdersWall(mockCoin)
	assert.Equal(t, config.ErrorNoServiceForCoin, err)
	_, err = rateService.GetCoinRates(mockCoin, false)
	assert.Equal(t, config.ErrorNoServiceForCoin, err)
}

func loadRateService() *RateSevice {
	rs := &RateSevice{
		FiatRates: &models.FiatRates{
			Rates:       nil,
			LastUpdated: time.Time{},
		},
		BittrexService:      bittrex.InitService(),
		BinanceService:      binance.InitService(),
		Crex24Service:       crex24.InitService(),
		StexService:         stex.InitService(),
		SouthXChangeService: southxhcange.InitService(),
		NovaExchangeService: novaexchange.InitService(),
		KuCoinService:       kucoin.InitService(),
		GraviexService:      graviex.InitService(),
		BitrueService:       bitrue.InitService(),
	}
	err := rs.LoadFiatRates()
	if err != nil {
		panic(err)
	}
	return rs
}

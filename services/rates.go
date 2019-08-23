package services

import (
	"encoding/json"
	"github.com/grupokindynos/obol/config"
	"github.com/grupokindynos/obol/models"
	coin_factory "github.com/grupokindynos/obol/models/coin-factory"
	"github.com/grupokindynos/obol/services/exchanges/binance"
	"github.com/grupokindynos/obol/services/exchanges/crex24"
	"github.com/grupokindynos/obol/services/exchanges/cryptobridge"
	"github.com/grupokindynos/obol/services/exchanges/stex"
	"io/ioutil"
	"math"
	"net/http"
)

type Exchange interface {
	CoinRate(coin string) (rate float64, err error)
	CoinMarketOrders(coin string) (orders []models.MarketOrder, err error)
}

type RateSevice struct {
	BinanceService      *binance.Service
	CryptoBridgeService *cryptobridge.Service
	Crex24Service       *crex24.Service
	StexService         *stex.Service
}

func (rs *RateSevice) GetCoinRates(coin *coin_factory.Coin) (rates map[string]models.Rate, err error) {
	btcRates, err := rs.GetBitpayRates()
	if err != nil {
		return rates, err
	}
	if coin.Tag == "BTC" {
		btcRatesMap := make(map[string]models.Rate)
		for _, rate := range btcRates {
			btcRatesMap[rate.Code] = rate
		}
		return btcRatesMap, nil
	} else {
		rate := rs.GetCoinExchangeRate(coin)
		newRates := make(map[string]models.Rate)
		for _, singleRate := range btcRates {
			newRate := models.Rate{
				Code: singleRate.Code,
				Name: singleRate.Name,
			}
			if singleRate.Code == "BTC" || singleRate.Code == "BCH" || singleRate.Code == "ETH" {
				newRate.Rate = math.Floor((rate*singleRate.Rate)*1e8) / 1e8
			} else {
				newRate.Rate = math.Floor((rate*singleRate.Rate)*10000) / 10000
			}
			newRates[singleRate.Code] = newRate
		}
		return newRates, nil
	}
}

func (rs *RateSevice) GetCoinToCoinRates(coinFrom *coin_factory.Coin, coinTo *coin_factory.Coin) (rate float64, err error) {
	coinFromRates, err := rs.GetCoinRates(coinFrom)
	if err != nil {
		return rate, err
	}
	coinToRates, err := rs.GetCoinRates(coinTo)
	if err != nil {
		return rate, err
	}
	coinFromCommonRate := coinFromRates["BTC"].Rate
	coinToCommonRate := coinToRates["BTC"].Rate
	rate = math.Floor(coinToCommonRate/coinFromCommonRate*1e8) / 1e8
	return rate, nil
}

func (rs *RateSevice) GetCoinToCoinRatesWithAmount(coinFrom *coin_factory.Coin, coinTo *coin_factory.Coin, amount float64) (rate float64, err error) {
	coinFromMarkets, err := rs.GetCoinOrdersWallet(coinFrom)
	if err != nil {
		return rate, err
	}
	coinToRates, err := rs.GetCoinRates(coinTo)
	if err != nil {
		return rate, err
	}
	coinToBTCRate := coinToRates["BTC"].Rate
	var countedAmount float64
	var pricesSum float64
	for _, order := range coinFromMarkets {
		if countedAmount+order.Amount >= amount {
			diff := math.Abs((countedAmount + order.Amount) - amount)
			newAmount := order.Amount - diff
			countedAmount += newAmount
			percentaje := newAmount / amount
			pricesSum += order.Price * percentaje
		} else {
			countedAmount += order.Amount
			percentaje := order.Amount / amount
			pricesSum += order.Price * percentaje
		}
		if countedAmount >= amount {
			break
		}
	}
	finaleRate := math.Floor((pricesSum/coinToBTCRate)*1e8) / 1e8
	return finaleRate, nil
}

func (rs *RateSevice) GetBitpayRates() (rates []models.Rate, err error) {
	res, err := http.Get(config.BitpayRatesURL)
	if err != nil {
		return rates, err
	} else {
		defer func() {
			_ = res.Body.Close()
		}()
		contents, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return rates, err
		}
		var BitpayRates models.BitpayRates
		err = json.Unmarshal(contents, &BitpayRates)
		if err != nil {
			return rates, err
		}
		return BitpayRates.Data, err
	}
}

func (rs *RateSevice) GetCoinExchangeRate(coin *coin_factory.Coin) float64 {
	var service Exchange
	switch coin.Exchange {
	case "binance":
		service = rs.BinanceService
	case "cryptobridge":
		service = rs.CryptoBridgeService
	case "crex24":
		service = rs.Crex24Service
	case "stex":
		service = rs.StexService
	}
	if service != nil {
		rate, err := service.CoinRate(coin.Tag)
		if err != nil {
			return 0
		}
		return rate
	}

	return 0
}

func (rs *RateSevice) GetCoinOrdersWallet(coin *coin_factory.Coin) ([]models.MarketOrder, error) {
	var service Exchange
	switch coin.Exchange {
	case "binance":
		service = rs.BinanceService
	case "cryptobridge":
		service = rs.CryptoBridgeService
	case "crex24":
		service = rs.Crex24Service
	case "stex":
		service = rs.StexService
	}
	if service != nil {
		orders, err := service.CoinMarketOrders(coin.Tag)
		if err != nil {
			return []models.MarketOrder{}, err
		}
		return orders, nil
	}
	return []models.MarketOrder{}, nil
}

func InitRateService() *RateSevice {
	rs := &RateSevice{
		BinanceService:      binance.InitService(),
		CryptoBridgeService: cryptobridge.InitService(),
		Crex24Service:       crex24.InitService(),
		StexService:         stex.InitService(),
	}
	return rs
}

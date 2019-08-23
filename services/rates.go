package services

import (
	"encoding/json"
	"fmt"
	"github.com/grupokindynos/obol/config"
	"github.com/grupokindynos/obol/models"
	coinfactory "github.com/grupokindynos/obol/models/coin-factory"
	"github.com/grupokindynos/obol/models/exchanges"
	"github.com/grupokindynos/obol/services/exchanges/binance"
	"github.com/grupokindynos/obol/services/exchanges/bittrex"
	"github.com/grupokindynos/obol/services/exchanges/crex24"
	"github.com/grupokindynos/obol/services/exchanges/cryptobridge"
	"github.com/grupokindynos/obol/services/exchanges/stex"
	"io/ioutil"
	"math"
	"strconv"
	"time"
)

const UpdateFiatRatesTimeFrame = 60 * 60 // 1 Hour timeframe

type Exchange interface {
	CoinRate(coin string) (rate float64, err error)
	CoinMarketOrders(coin string) (orders []models.MarketOrder, err error)
}

var FiatRates *models.FiatRates

type RateSevice struct {
	FiatRates           *models.FiatRates
	BittrexService 		*bittrex.Service
	BinanceService      *binance.Service
	CryptoBridgeService *cryptobridge.Service
	Crex24Service       *crex24.Service
	StexService         *stex.Service
}

func (rs *RateSevice) GetCoinRates(coin *coinfactory.Coin) (rates map[string]float64, err error) {
	btcRates, err := rs.GetBtcRates()
	if err != nil {
		return rates, err
	}
	if coin.Tag == "BTC" {
		btcRatesMap := make(map[string]float64)
		for code, rate := range btcRates {
			btcRatesMap[code] = rate
		}
		return btcRatesMap, nil
	} else {
		rate, err := rs.GetCoinExchangeRate(coin)
		newRates := make(map[string]float64)
		for code, singleRate := range btcRates {
			if code == "BTC" {
				newRates[code] = math.Floor((rate*singleRate)*1e8) / 1e8
			} else {
				newRates[code] = math.Floor((rate*singleRate)*10000) / 10000
			}
		}
		return newRates, err
	}
}

func (rs *RateSevice) GetCoinToCoinRates(coinFrom *coinfactory.Coin, coinTo *coinfactory.Coin) (rate float64, err error) {
	coinFromRates, err := rs.GetCoinRates(coinFrom)
	if err != nil {
		return rate, err
	}
	coinToRates, err := rs.GetCoinRates(coinTo)
	if err != nil {
		return rate, err
	}
	coinFromCommonRate := coinFromRates["BTC"]
	coinToCommonRate := coinToRates["BTC"]
	rate = math.Floor(coinToCommonRate/coinFromCommonRate*1e8) / 1e8
	return rate, nil
}

func (rs *RateSevice) GetCoinToCoinRatesWithAmount(coinFrom *coinfactory.Coin, coinTo *coinfactory.Coin, amount float64) (rate float64, err error) {
	coinFromMarkets, err := rs.GetCoinOrdersWall(coinFrom)
	if err != nil {
		return rate, err
	}
	coinToRates, err := rs.GetCoinRates(coinTo)
	if err != nil {
		return rate, err
	}
	coinToBTCRate := coinToRates["BTC"]
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

func (rs *RateSevice) GetCoinExchangeRate(coin *coinfactory.Coin) (float64, error) {
	var service Exchange
	switch coin.Exchange {
	case "binance":
		service = rs.BinanceService
	case "bittrex":
		service = rs.BittrexService
	case "cryptobridge":
		service = rs.CryptoBridgeService
	case "crex24":
		service = rs.Crex24Service
	case "stex":
		service = rs.StexService
	}
	if service != nil {
		rate, err := service.CoinRate(coin.Tag)
		return rate, err
	}
	return 0, config.ErrorNoServiceForCoin
}

func (rs *RateSevice) GetCoinOrdersWall(coin *coinfactory.Coin) (orders []models.MarketOrder, err error) {
	var service Exchange
	switch coin.Exchange {
	case "binance":
		service = rs.BinanceService
	case "bittrex":
		service = rs.BittrexService
	case "cryptobridge":
		service = rs.CryptoBridgeService
	case "crex24":
		service = rs.Crex24Service
	case "stex":
		service = rs.StexService
	}
	if service != nil {
		orders, err = service.CoinMarketOrders(coin.Tag)
		return orders, err
	}
	return orders, config.ErrorNoServiceForCoin
}

func (rs *RateSevice) GetBtcMxnRate() (float64, error) {
	res, err := config.HttpClient.Get("https://api.bitso.com/v3/ticker/?book=btc_mxn")
	if err != nil {
		return 0, err
	} else {
		defer func() {
			_ = res.Body.Close()
		}()
		contents, _ := ioutil.ReadAll(res.Body)
		var bitsoRates exchanges.BitsoRates
		err = json.Unmarshal(contents, &bitsoRates)
		rate, err := strconv.ParseFloat(bitsoRates.Payload.Last, 64)
		return rate, err
	}
}

func (rs *RateSevice) GetBtcRates() (rates map[string]float64, err error) {
	if rs.FiatRates.LastUpdated.Unix()+UpdateFiatRatesTimeFrame > time.Now().Unix() {
		loadFiatRates()
	}
	mxnRate, err := rs.GetBtcMxnRate()
	rates = make(map[string]float64)
	for key, rate := range FiatRates.Rates {
		rates[key] = rate * mxnRate
	}
	rates["BTC"] = 1
	return rates, err
}

func loadFiatRates() {
	res, err := config.HttpClient.Get(config.OpenRatesURL)
	if err != nil {
		fmt.Println("unable to load fiat rates")
	} else {
		defer func() {
			_ = res.Body.Close()
		}()
		contents, _ := ioutil.ReadAll(res.Body)
		var fiatRates models.OpenRates
		_ = json.Unmarshal(contents, &fiatRates)
		rateBytes, _ := json.Marshal(fiatRates.Rates)
		ratesMap := make(map[string]float64)
		_ = json.Unmarshal(rateBytes, &ratesMap)
		FiatRates = &models.FiatRates{
			Rates:       ratesMap,
			LastUpdated: time.Now(),
		}
	}
}

func InitRateService() *RateSevice {
	loadFiatRates()
	rs := &RateSevice{
		FiatRates:           FiatRates,
		BittrexService: bittrex.InitService(),
		BinanceService:      binance.InitService(),
		CryptoBridgeService: cryptobridge.InitService(),
		Crex24Service:       crex24.InitService(),
		StexService:         stex.InitService(),
	}
	return rs
}

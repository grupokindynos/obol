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
	"github.com/grupokindynos/obol/services/exchanges/novaexchange"
	"github.com/grupokindynos/obol/services/exchanges/southxhcange"
	"github.com/grupokindynos/obol/services/exchanges/stex"
	"io/ioutil"
	"math"
	"strconv"
	"time"
)

const UpdateFiatRatesTimeFrame = 60 * 60 // 1 Hour timeframe

//Exchange is the interface to make sure all exchange services have the same properties
type Exchange interface {
	CoinMarketOrders(coin string) (orders []models.MarketOrder, err error)
}

var FiatRates *models.FiatRates

type RateSevice struct {
	FiatRates           *models.FiatRates
	BittrexService      *bittrex.Service
	BinanceService      *binance.Service
	CryptoBridgeService *cryptobridge.Service
	Crex24Service       *crex24.Service
	StexService         *stex.Service
	SouthXChangeService *southxhcange.Service
	NovaExchangeService *novaexchange.Service
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
	}
	ratesWall, err := rs.GetCoinOrdersWall(coin)
	if err != nil {
		return rates, err
	}
	newRates := make(map[string]float64)
	for code, singleRate := range btcRates {
		if code == "BTC" {
			newRates[code] = math.Floor((ratesWall[0].Price*singleRate)*1e8) / 1e8
		} else {
			newRates[code] = math.Floor((ratesWall[0].Price*singleRate)*10000) / 10000
		}
	}
	return newRates, err
}

func (rs *RateSevice) GetCoinToCoinRates(coinFrom *coinfactory.Coin, coinTo *coinfactory.Coin) (rate float64, err error) {
	if coinFrom.Tag == "BTC" {
		coinRates, err := rs.GetCoinRates(coinTo)
		return coinRates["BTC"], err
	}
	if coinTo.Tag == "BTC" {
		coinRates, err := rs.GetCoinRates(coinFrom)
		return 1 / coinRates["BTC"], err
	}
	if coinFrom.Tag == coinTo.Tag {
		return rate, config.ErrorNoC2CWithSameCoin
	}
	coinFromRates, err := rs.GetCoinRates(coinFrom)
	coinToRates, err := rs.GetCoinRates(coinTo)
	coinFromCommonRate := coinFromRates["BTC"]
	coinToCommonRate := coinToRates["BTC"]
	return coinToCommonRate / coinFromCommonRate, err
}

func (rs *RateSevice) GetCoinToCoinRatesWithAmount(coinFrom *coinfactory.Coin, coinTo *coinfactory.Coin, amount float64) (rate float64, err error) {
	if coinFrom.Tag == "BTC" {
		return rate, config.ErrorNoC2CWithBTC
	}
	if coinFrom.Tag == coinTo.Tag {
		return rate, config.ErrorNoC2CWithSameCoin
	}
	// First get the orders wall from the coin we are converting
	coinFromMarkets, err := rs.GetCoinOrdersWall(coinFrom)
	// Get BTC rate of the coin.
	coinToRates, err := rs.GetCoinRates(coinTo)
	coinToBTCRate := coinToRates["BTC"]
	// Init vars for loop
	var countedAmount float64
	var pricesSum float64
	// Looping against values on exchange to make a approachable rate based on the amount.
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
	priceTrunk := math.Floor(pricesSum*1e8) / 1e8
	finaleRate := coinToBTCRate / priceTrunk
	return finaleRate, err
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
	case "novaexchange":
		service = rs.NovaExchangeService
	case "stex":
		service = rs.StexService
	case "southxchange":
		service = rs.SouthXChangeService
	}
	if service != nil {
		orders, err = service.CoinMarketOrders(coin.Tag)
		if err != nil {
			var fallBackService Exchange
			switch coin.FallBackExchange {
			case "binance":
				fallBackService = rs.BinanceService
			case "bittrex":
				fallBackService = rs.BittrexService
			case "cryptobridge":
				fallBackService = rs.CryptoBridgeService
			case "crex24":
				fallBackService = rs.Crex24Service
			case "novaexchange":
				fallBackService = rs.NovaExchangeService
			case "stex":
				fallBackService = rs.StexService
			case "southxchange":
				fallBackService = rs.SouthXChangeService
			}
			if fallBackService != nil {
				fallBackOrders, err := fallBackService.CoinMarketOrders(coin.Tag)
				return fallBackOrders, err
			}
		}
		return orders, err
	}
	return orders, config.ErrorNoServiceForCoin
}

func (rs *RateSevice) GetBtcMxnRate() (float64, error) {
	res, err := config.HttpClient.Get("https://api.bitso.com/v3/ticker/?book=btc_mxn")
	if err != nil {
		return 0, config.ErrorRequestTimeout
	}
	defer func() {
		_ = res.Body.Close()
	}()
	contents, _ := ioutil.ReadAll(res.Body)
	var bitsoRates exchanges.BitsoRates
	err = json.Unmarshal(contents, &bitsoRates)
	rate, err := strconv.ParseFloat(bitsoRates.Payload.Last, 64)
	return rate, err
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
		BittrexService:      bittrex.InitService(),
		BinanceService:      binance.InitService(),
		CryptoBridgeService: cryptobridge.InitService(),
		Crex24Service:       crex24.InitService(),
		StexService:         stex.InitService(),
		SouthXChangeService: southxhcange.InitService(),
		NovaExchangeService: novaexchange.InitService(),
	}
	return rs
}

package services

import (
	"encoding/json"
	"fmt"
	"github.com/grupokindynos/common/coin-factory/coins"
	"github.com/grupokindynos/obol/config"
	"github.com/grupokindynos/obol/models"
	"github.com/grupokindynos/obol/models/exchanges"
	"github.com/grupokindynos/obol/services/exchanges/binance"
	"github.com/grupokindynos/obol/services/exchanges/bittrex"
	"github.com/grupokindynos/obol/services/exchanges/crex24"
	"github.com/grupokindynos/obol/services/exchanges/cryptobridge"
	"github.com/grupokindynos/obol/services/exchanges/kucoin"
	"github.com/grupokindynos/obol/services/exchanges/novaexchange"
	"github.com/grupokindynos/obol/services/exchanges/southxhcange"
	"github.com/grupokindynos/obol/services/exchanges/stex"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"time"
)

// UpdateFiatRatesTimeFrame is the time frame to update fiat rates
const UpdateFiatRatesTimeFrame = 60 * 60 // 1 Hour timeframe

//Exchange is the interface to make sure all exchange services have the same properties
type Exchange interface {
	CoinMarketOrders(coin string) (orders map[string][]models.MarketOrder, err error)
}

// FiatRates is the model of the OpenRates fiat rates information
var FiatRates *models.FiatRates

// RateSevice is the main wrapper for all different exchanges and fiat rates data
type RateSevice struct {
	FiatRates           *models.FiatRates
	BittrexService      *bittrex.Service
	BinanceService      *binance.Service
	CryptoBridgeService *cryptobridge.Service
	Crex24Service       *crex24.Service
	StexService         *stex.Service
	SouthXChangeService *southxhcange.Service
	NovaExchangeService *novaexchange.Service
	KuCoinService       *kucoin.Service
}

// GetCoinRates is the main function to get the rates of a coin using the OpenRates structure
func (rs *RateSevice) GetCoinRates(coin *coins.Coin, buyWall bool) (rates []models.Rate, err error) {
	btcRates, err := rs.GetBtcRates()
	if err != nil {
		return rates, err
	}
	if coin.Tag == "BTC" {
		return btcRates, nil
	}
	ratesWall, err := rs.GetCoinOrdersWall(coin)
	if err != nil {
		return rates, err
	}
	var orders []models.MarketOrder
	if buyWall {
		orders = ratesWall["buy"]
	} else {
		orders = ratesWall["sell"]
	}
	for _, singleRate := range btcRates {
		rate := models.Rate{
			Code: singleRate.Code,
			Name: singleRate.Name,
		}
		if singleRate.Code == "BTC" {
			rate.Rate = math.Floor((orders[0].Price*singleRate.Rate)*1e8) / 1e8
		} else {
			rate.Rate = math.Floor((orders[0].Price*singleRate.Rate)*10000) / 10000
		}
		rates = append(rates, rate)
	}
	return rates, err
}

// GetCoinToCoinRates will return the rates from a crypto to a crypto using the exchanges data
func (rs *RateSevice) GetCoinToCoinRates(coinFrom *coins.Coin, coinTo *coins.Coin) (rate float64, err error) {
	if coinFrom.Tag == "BTC" {
		coinRates, err := rs.GetCoinRates(coinTo, true)
		for _, rate := range coinRates {
			if rate.Code == "BTC" {
				return rate.Rate, err
			}
		}
	}
	if coinTo.Tag == "BTC" {
		coinRates, err := rs.GetCoinRates(coinFrom, false)
		for _, rate := range coinRates {
			if rate.Code == "BTC" {
				return 1 / rate.Rate, err
			}
		}
	}
	if coinFrom.Tag == coinTo.Tag {
		return rate, config.ErrorNoC2CWithSameCoin
	}
	coinFromRates, err := rs.GetCoinRates(coinFrom, false)
	coinToRates, err := rs.GetCoinRates(coinTo, false)
	var coinFromCommonRate float64
	var coinToCommonRate float64
	for _, rate := range coinFromRates {
		if rate.Code == "BTC" {
			coinFromCommonRate = rate.Rate
		}
	}
	for _, rate := range coinToRates {
		if rate.Code == "BTC" {
			coinToCommonRate = rate.Rate
		}
	}
	return coinToCommonRate / coinFromCommonRate, err
}

// GetCoinToCoinRatesWithAmount is used to get the rates from crypto to crypto using a specified amount to convert
func (rs *RateSevice) GetCoinToCoinRatesWithAmount(coinFrom *coins.Coin, coinTo *coins.Coin, amount float64) (rate float64, err error) {
	if coinFrom.Tag == coinTo.Tag {
		return rate, config.ErrorNoC2CWithSameCoin
	}
	var coinMarkets map[string][]models.MarketOrder
	var coinRates []models.Rate
	// First get the orders wall from the coin we are converting
	if coinFrom.Tag == "BTC" {
		coinMarkets, err = rs.GetCoinOrdersWall(coinTo)
		if err != nil {
			return 0, err
		}
		coinRates, err = rs.GetCoinRates(coinTo, true)
	} else {
		coinMarkets, err = rs.GetCoinOrdersWall(coinFrom)
		if err != nil {
			return 0, err
		}
		coinRates, err = rs.GetCoinRates(coinTo, false)
	}
	// Get BTC rate of the coin.
	var coinToBTCRate float64
	for _, rate := range coinRates {
		if rate.Code == "BTC" {
			coinToBTCRate = rate.Rate
		}
	}
	// Init vars for loop
	var countedAmount float64
	var pricesSum float64
	var orders []models.MarketOrder
	if coinFrom.Tag == "BTC" {
		orders = coinMarkets["buy"]
	} else {
		orders = coinMarkets["sell"]
	}
	// Looping against values on exchange to make a approachable rate based on the amount.
	for _, order := range orders {
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
	var finaleRate float64
	if coinFrom.Tag == "BTC" {
		finaleRate = priceTrunk
	} else {
		finaleRate = coinToBTCRate / priceTrunk
	}
	return finaleRate, err
}

// GetCoinOrdersWall will return the buy/sell orders from selected or fallback exchange
func (rs *RateSevice) GetCoinOrdersWall(coin *coins.Coin) (orders map[string][]models.MarketOrder, err error) {
	var service Exchange
	switch coin.Rates.Exchange {
	case "binance":
		service = rs.BinanceService
	case "bittrex":
		service = rs.BittrexService
	case "cryptobridge":
		service = rs.CryptoBridgeService
	case "crex24":
		service = rs.Crex24Service
	case "kucoin":
		service = rs.KuCoinService
	case "novaexchange":
		service = rs.NovaExchangeService
	case "stex":
		service = rs.StexService
	case "southxchange":
		service = rs.SouthXChangeService
	}
	if service == nil {
		return nil, config.ErrorNoServiceForCoin
	}
	orders, err = service.CoinMarketOrders(coin.Tag)
	if err != nil {
		var fallBackService Exchange
		switch coin.Rates.FallBackExchange {
		case "binance":
			fallBackService = rs.BinanceService
		case "bittrex":
			fallBackService = rs.BittrexService
		case "cryptobridge":
			fallBackService = rs.CryptoBridgeService
		case "kucoin":
			fallBackService = rs.KuCoinService
		case "crex24":
			fallBackService = rs.Crex24Service
		case "novaexchange":
			fallBackService = rs.NovaExchangeService
		case "stex":
			fallBackService = rs.StexService
		case "southxchange":
			fallBackService = rs.SouthXChangeService
		}
		if fallBackService == nil {
			return nil, config.ErrorNoFallBackServiceForCoin
		}
		fallBackOrders, err := fallBackService.CoinMarketOrders(coin.Tag)
		return fallBackOrders, err
	}
	return orders, err
}

// GetBtcMxnRate will return the price of BTC on MXN
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
	if err != nil {
		return 0, err
	}
	rate, err := strconv.ParseFloat(bitsoRates.Payload.Last, 64)
	return rate, err
}

// GetBtcRates will return the Bitcoin rates using the OpenRates structure
func (rs *RateSevice) GetBtcRates() (rates []models.Rate, err error) {
	if rs.FiatRates.LastUpdated.Unix()+UpdateFiatRatesTimeFrame > time.Now().Unix() {
		loadFiatRates()
	}
	btcMxnRate, err := rs.GetBtcMxnRate()
	newRate := btcMxnRate / FiatRates.Rates["MXN"]
	for key, rate := range FiatRates.Rates {
		rate := models.Rate{
			Code: key,
			Name: models.FixerRatesNames[key],
			Rate: rate * newRate,
		}
		rates = append(rates, rate)
	}
	btcRate := models.Rate{
		Code: "BTC",
		Name: "Bitcoin",
		Rate: 1,
	}
	rates = append(rates, btcRate)
	return rates, err
}

func loadFiatRates() {
	res, err := config.HttpClient.Get(config.FixerRatesURL + "?access_key=" + os.Getenv("FIXER_RATES_TOKEN"))
	if err != nil {
		fmt.Println("unable to load fiat rates")
	} else {
		defer func() {
			_ = res.Body.Close()
		}()
		contents, _ := ioutil.ReadAll(res.Body)
		var fiatRates models.FixerRates
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

// InitRateService is a safe to use function to init the rate service.
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
		KuCoinService:       kucoin.InitService(),
	}
	return rs
}

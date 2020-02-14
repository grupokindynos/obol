package services

import (
	"encoding/json"
	"errors"
	coinfactory "github.com/grupokindynos/common/coin-factory"
	"github.com/grupokindynos/common/coin-factory/coins"
	"github.com/grupokindynos/common/obol"
	"github.com/grupokindynos/obol/config"
	"github.com/grupokindynos/obol/models"
	"github.com/grupokindynos/obol/models/exchanges"
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
	"github.com/olympus-protocol/ogen/utils/amount"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"time"
)

func init() {
	_ = godotenv.Load("../.env")
}

const (
	// UpdateFiatRatesTimeFrame is the time frame to update fiat rates
	UpdateFiatRatesTimeFrame = 60 * 60 * 24 // 24 Hour timeframe
	UpdateBtcRatesTimeFrame  = 60 * 15      // 15 minutes
)

//Exchange is the interface to make sure all exchange services have the same properties
type Exchange interface {
	CoinMarketOrders(coin string) (orders map[string][]models.MarketOrder, err error)
}

type BtcRates struct {
	LastUpdated int64
	Rates       []models.Rate
}

// RateSevice is the main wrapper for all different exchanges and fiat rates data
type RateSevice struct {
	FiatRates           *models.FiatRates
	FiatRatesToken      string
	BtcRates            BtcRates
	BittrexService      *bittrex.Service
	BinanceService      *binance.Service
	Crex24Service       *crex24.Service
	StexService         *stex.Service
	SouthXChangeService *southxhcange.Service
	NovaExchangeService *novaexchange.Service
	KuCoinService       *kucoin.Service
	GraviexService      *graviex.Service
	BitrueService       *bitrue.Service
}

// GetCoinRates is the main function to get the rates of a coin using the OpenRates structure
func (rs *RateSevice) GetCoinRates(coin *coins.Coin, buyWall bool) (rates []models.Rate, err error) {
	btcRates, err := rs.GetBtcRates()
	if err != nil {
		return rates, err
	}
	if coin.Info.Tag == "BTC" {
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
	orderPrice := orders[0].Price
	for _, singleRate := range btcRates {
		singleRateConv, err := amount.NewAmount(singleRate.Rate)
		if err != nil {
			return nil, err
		}
		rate := models.Rate{
			Code: singleRate.Code,
			Name: singleRate.Name,
		}
		var rateNum float64
		if coin.Info.Tag == "USDC" || coin.Info.Tag == "TUSD" || coin.Info.Tag == "USDT" {
			rateNum = singleRateConv.ToNormalUnit() / orderPrice.ToNormalUnit()
		} else {
			rateNum = orderPrice.ToNormalUnit() * singleRateConv.ToNormalUnit()
		}
		if singleRate.Code == "BTC" {
			rate.Rate = toFixed(rateNum, 8)
		} else {
			rate.Rate = toFixed(rateNum, 4)
		}
		rates = append(rates, rate)
	}
	return rates, err
}

// GetCoinToCoinRates will return the rates from a crypto to a crypto using the exchanges data
func (rs *RateSevice) GetCoinToCoinRates(coinFrom *coins.Coin, coinTo *coins.Coin) (rate float64, err error) {
	if coinFrom.Info.Tag == coinTo.Info.Tag {
		return rate, config.ErrorNoC2CWithSameCoin
	}
	if coinTo.Info.Tag == "BTC" {
		coinRates, err := rs.GetCoinRates(coinFrom, false)
		for _, rate := range coinRates {
			if rate.Code == "BTC" {
				return rate.Rate, err
			}
		}
	}
	coinFromRates, err := rs.GetCoinRates(coinFrom, false)
	if err != nil {
		return 0, err
	}
	coinToRates, err := rs.GetCoinRates(coinTo, false)
	if err != nil {
		return 0, err
	}
	var coinFromUSDRate float64
	for _, rate := range coinFromRates {
		if rate.Code == "USD" {
			coinFromUSDRate = rate.Rate
		}
	}
	var coinToUSDRate float64
	for _, rate := range coinToRates {
		if rate.Code == "USD" {
			coinToUSDRate = rate.Rate
		}
	}
	return toFixed(coinFromUSDRate/coinToUSDRate, 6), nil
}

func (rs *RateSevice) GetCoinLiquidity(coin *coins.Coin) (float64, error) {
	coinWalls, err := rs.GetCoinOrdersWall(coin)
	if err != nil {
		return 0, err
	}
	orderWall := coinWalls["sell"]
	var liquidity float64
	for _, order := range orderWall {
		liquidity += order.Amount * order.Price.ToNormalUnit()
	}
	btcRates, err := rs.GetBtcRates()
	if err != nil {
		return 0, err
	}
	var btcUSDRate float64
	for _, rate := range btcRates {
		if rate.Code == "USD" {
			btcUSDRate = rate.Rate
		}
	}
	return toFixed(liquidity*btcUSDRate, 8), err
}

// GetCoinToCoinRatesWithAmount is used to get the rates from crypto to crypto using a specified amount to convert
func (rs *RateSevice) GetCoinToCoinRatesWithAmount(coinFrom *coins.Coin, coinTo *coins.Coin, amountReq float64) (obol.CoinToCoinWithAmountResponse, error) {
	if coinFrom.Info.Tag == coinTo.Info.Tag {
		return obol.CoinToCoinWithAmountResponse{}, config.ErrorNoC2CWithSameCoin
	}
	amountRequested := toFixed(amountReq, 6)
	if amountRequested <= 0 {
		return obol.CoinToCoinWithAmountResponse{}, errors.New("amount must be greater than 0")
	}
	var coinWall []models.MarketOrder
	if coinFrom.Info.Tag == "BTC" {
		coinToWalls, err := rs.GetCoinOrdersWall(coinTo)
		if err != nil {
			return obol.CoinToCoinWithAmountResponse{}, err
		}
		coinWall = coinToWalls["buy"]
	} else {
		coinFromWalls, err := rs.GetCoinOrdersWall(coinFrom)
		if err != nil {
			return obol.CoinToCoinWithAmountResponse{}, err
		}
		coinWall = coinFromWalls["sell"]
	}
	amountParsed := amountRequested
	btcData, err := coinfactory.GetCoin("BTC")
	if err != nil {
		return obol.CoinToCoinWithAmountResponse{}, err
	}
	if amountRequested <= coinWall[0].Amount {
		return obol.CoinToCoinWithAmountResponse{
			AveragePrice: toFixed(coinWall[0].Price.ToNormalUnit(), 8),
		}, nil
	}

	var rates [][]float64
	var percentageSum float64
	for _, order := range coinWall {
		percentage := toFixed(order.Amount/amountReq, 6)
		percentageSum += percentage
		var orderArr []float64
		if percentageSum > 1 {
			exceed := percentageSum - 1
			rest := percentage - exceed
			orderArr = []float64{order.Amount, order.Price.ToNormalUnit(), toFixed(rest, 6)}
			percentageSum -= exceed
		} else {
			orderArr = []float64{order.Amount, order.Price.ToNormalUnit(), percentage}
		}
		rates = append(rates, orderArr)
		amountParsed -= order.Amount
		if amountParsed <= 0 {
			break
		}
	}
	amountParsed = amountRequested
	var AvrPrice float64
	for _, rateFloat := range rates {
		AvrPrice += rateFloat[1] * rateFloat[2]
		amountParsed -= rateFloat[0]
		if amountParsed <= 0 {
			break
		}
	}
	var rate obol.CoinToCoinWithAmountResponse
	if coinTo.Info.Tag == "BTC" || coinFrom.Info.Tag == "BTC" {
		rate.AveragePrice = toFixed(AvrPrice, 8)
	} else {
		rateConv, err := rs.GetCoinToCoinRates(coinTo, btcData)
		if err != nil {
			return rate, err
		}
		rate.AveragePrice = toFixed(AvrPrice/rateConv, 8)
	}
	return rate, err
}

// GetCoinOrdersWall will return the buy/sell orders from selected or fallback exchange
func (rs *RateSevice) GetCoinOrdersWall(coin *coins.Coin) (orders map[string][]models.MarketOrder, err error) {
	var service Exchange
	switch coin.Rates.Exchange {
	case "binance":
		service = rs.BinanceService
	case "bittrex":
		service = rs.BittrexService
	case "bitrue":
		service = rs.BitrueService
	case "crex24":
		service = rs.Crex24Service
	case "kucoin":
		service = rs.KuCoinService
	case "graviex":
		service = rs.GraviexService
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
	orders, err = service.CoinMarketOrders(coin.Info.Tag)
	if err != nil {
		var fallBackService Exchange
		switch coin.Rates.FallBackExchange {
		case "binance":
			fallBackService = rs.BinanceService
		case "bittrex":
			fallBackService = rs.BittrexService
		case "bitrue":
			fallBackService = rs.BitrueService
		case "kucoin":
			fallBackService = rs.KuCoinService
		case "graviex":
			fallBackService = rs.GraviexService
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
		fallBackOrders, err := fallBackService.CoinMarketOrders(coin.Info.Tag)
		return fallBackOrders, err
	}
	return orders, err
}

// GetBtcEURRate will return the price of BTC on EUR
func (rs *RateSevice) GetBtcEURRate() (float64, error) {
	res, err := config.HttpClient.Get("https://bitstamp.net/api/v2/ticker/btceur")
	if err != nil {
		return 0, config.ErrorRequestTimeout
	}
	defer func() {
		_ = res.Body.Close()
	}()
	contents, _ := ioutil.ReadAll(res.Body)
	var rate exchanges.BitstampRate
	err = json.Unmarshal(contents, &rate)
	if err != nil {
		return 0, err
	}
	rateNum, err := strconv.ParseFloat(rate.Last, 64)
	if err != nil {
		return 0, err
	}
	return rateNum, err
}

// GetBtcRates will return the Bitcoin rates using the OpenRates structure
func (rs *RateSevice) GetBtcRates() (rates []models.Rate, err error) {
	if rs.FiatRates.LastUpdated.Unix()+UpdateFiatRatesTimeFrame < time.Now().Unix() {
		err = rs.LoadFiatRates()
		if err != nil {
			return nil, err
		}
	}
	if rs.BtcRates.LastUpdated+UpdateBtcRatesTimeFrame > time.Now().Unix() {
		return rs.BtcRates.Rates, nil
	}
	btcRate, err := rs.GetBtcEURRate()
	for key, rate := range rs.FiatRates.Rates {
		newRate, err := amount.NewAmount(rate * btcRate)
		if err != nil {
			return nil, err
		}
		rate := models.Rate{
			Code: key,
			Name: models.FixerRatesNames[key],
			Rate: newRate.ToNormalUnit(),
		}
		rates = append(rates, rate)
	}
	return rates, err
}

func (rs *RateSevice) LoadFiatRates() error {
	res, err := config.HttpClient.Get(config.FixerRatesURL + "?access_key=" + os.Getenv("FIXER_RATES_TOKEN"))
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	contents, _ := ioutil.ReadAll(res.Body)
	var fiatRates models.FixerRates
	err = json.Unmarshal(contents, &fiatRates)
	if err != nil {
		return err
	}
	if fiatRates.Error.Code != 0 {
		panic("unable to load fiat rates")
	}
	rateBytes, err := json.Marshal(fiatRates.Rates)
	if err != nil {
		return err
	}
	ratesMap := make(map[string]float64)
	err = json.Unmarshal(rateBytes, &ratesMap)
	if err != nil {
		return err
	}
	rs.FiatRates = &models.FiatRates{
		Rates:       ratesMap,
		LastUpdated: time.Now(),
	}
	return nil
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

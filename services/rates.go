package services

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/grupokindynos/obol/services/exchanges/folgory"
	"github.com/grupokindynos/obol/services/exchanges/hitbtc"
	"github.com/grupokindynos/obol/services/exchanges/lukki"
	"github.com/shopspring/decimal"

	coinFactory "github.com/grupokindynos/common/coin-factory"
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
	Rates       map[string]models.RateV2
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
	FolgoryService      *folgory.Service
	HitBTCService       *hitbtc.Service
	LukkiService        *lukki.Service
}

// GetCoinRates is the main function to get the rates of a coin using the OpenRates structure
func (rs *RateSevice) GetCoinRates(coin *coins.Coin, buyWall bool) (map[string]models.RateV2, error) {
	rates := make(map[string]models.RateV2)
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
	for code, singleRate := range btcRates {
		rate := models.RateV2{
			Name: singleRate.Name,
		}
		var rateNum decimal.Decimal
		if coin.Info.Tag == "USDC" || coin.Info.Tag == "TUSD" || coin.Info.Tag == "USDT" {
			rateNum = singleRate.Rate.Div(orderPrice)
		} else {
			rateNum = orderPrice.Mul(singleRate.Rate)
		}
		if code == "BTC" {
			rate.Rate = rateNum.Round(8)
		} else {
			rate.Rate = rateNum.Round(6)
		}
		rates[code] = rate
	}
	return rates, err
}

// GetCoinToCoinRates will return the rates from a crypto to a crypto using the exchanges data
func (rs *RateSevice) GetCoinToCoinRates(coinFrom *coins.Coin, coinTo *coins.Coin) (rate decimal.Decimal, err error) {
	if coinFrom.Info.Tag == coinTo.Info.Tag {
		return rate, config.ErrorNoC2CWithSameCoin
	}
	if coinTo.Info.Tag == "BTC" {
		coinRates, err := rs.GetCoinRates(coinFrom, false)
		for code, rate := range coinRates {
			if code == "BTC" {
				return rate.Rate, err
			}
		}
	}
	coinFromRates, err := rs.GetCoinRates(coinFrom, false)
	if err != nil {
		return decimal.Zero, err
	}
	coinToRates, err := rs.GetCoinRates(coinTo, false)
	if err != nil {
		return decimal.Zero, err
	}
	var coinFromUSDRate decimal.Decimal
	for code, rate := range coinFromRates {
		if code == "USD" {
			coinFromUSDRate = rate.Rate
		}
	}
	var coinToUSDRate decimal.Decimal
	for code, rate := range coinToRates {
		if code == "USD" {
			coinToUSDRate = rate.Rate
		}
	}
	return coinFromUSDRate.DivRound(coinToUSDRate, 6), nil
}

func (rs *RateSevice) GetCoinLiquidity(coin *coins.Coin) (decimal.Decimal, error) {
	coinWalls, err := rs.GetCoinOrdersWall(coin)
	if err != nil {
		return decimal.Zero, err
	}
	orderWall := coinWalls["sell"]
	var liquidity decimal.Decimal
	for _, order := range orderWall {
		liquidity.Add(order.Amount.Mul(order.Price))
	}
	btcRates, err := rs.GetBtcRates()
	if err != nil {
		return decimal.Zero, err
	}
	var btcUSDRate decimal.Decimal
	for code, rate := range btcRates {
		if code == "USD" {
			btcUSDRate = rate.Rate
		}
	}
	return liquidity.Mul(btcUSDRate).Round(8), err
}

// GetCoinToCoinRatesWithAmount is used to get the rates from crypto to crypto using a specified amount to convert
func (rs *RateSevice) GetCoinToCoinRatesWithAmount(coinFrom *coins.Coin, coinTo *coins.Coin, amountReq float64) (obol.CoinToCoinWithAmountResponse, error) {
	if coinFrom.Info.Tag == coinTo.Info.Tag {
		return obol.CoinToCoinWithAmountResponse{}, config.ErrorNoC2CWithSameCoin
	}
	amountRequested := decimal.NewFromFloat(amountReq).Round(6)
	if amountRequested.LessThanOrEqual(decimal.Zero) {
		return obol.CoinToCoinWithAmountResponse{}, errors.New("amount must be greater than 0")
	}
	var coinWall []models.MarketOrder
	if coinFrom.Info.Tag == "BTC" {
		coinToWalls, err := rs.GetCoinOrdersWall(coinTo)
		if err != nil {
			return obol.CoinToCoinWithAmountResponse{}, err
		}
		coinWall = coinToWalls["buy"]
		for i := 0; i < len(coinWall); i++ {
			coinWall[i].Amount.Mul(coinWall[i].Price)
		}
	} else {
		coinFromWalls, err := rs.GetCoinOrdersWall(coinFrom)
		if err != nil {
			return obol.CoinToCoinWithAmountResponse{}, err
		}
		coinWall = coinFromWalls["sell"]
	}
	amountParsed := amountRequested
	btcData, err := coinFactory.GetCoin("BTC")
	if err != nil {
		return obol.CoinToCoinWithAmountResponse{}, err
	}

	var rates [][]float64
	var percentageSum decimal.Decimal
	for _, order := range coinWall {
		percentage := order.Amount.DivRound(amountRequested, 6)
		percentageSum.Add(percentage)
		var orderArr []float64
		if percentageSum.GreaterThan(decimal.NewFromInt(1)) {
			exceed := percentageSum.Sub(decimal.NewFromInt(1))
			rest := percentage.Sub(exceed)
			orderArr = []float64{order.Amount, order.Price, toFixed(rest, 6)}
			percentageSum.Sub(exceed)
		} else {
			orderArr = []float64{order.Amount, order.Price, percentage}
		}
		rates = append(rates, orderArr)
		amountParsed.Sub(order.Amount)
		if amountParsed.LessThanOrEqual(decimal.Zero) {
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
	if coinTo.Info.Tag == "BTC" {
		rate.AveragePrice = toFixed(AvrPrice, 8)
	} else if coinFrom.Info.Tag == "BTC" {
		rate.AveragePrice = toFixed(1.0/AvrPrice, 8)
	} else {
		rateConv, err := rs.GetCoinToCoinRates(coinTo, btcData)
		if err != nil {
			return rate, err
		}
		rate.AveragePrice = toFixed(AvrPrice/rateConv, 8)
	}
	rate.Amount = rate.AveragePrice * amountRequested
	return rate, err
}

// GetCoinOrdersWall will return the buy/sell orders from selected or fallback exchange
func (rs *RateSevice) GetCoinOrdersWall(coin *coins.Coin) (orders map[string][]models.MarketOrder, err error) {
	var service Exchange
	switch coin.Rates.Exchange {
	case "binance":
		service = rs.BinanceService
	case "folgory":
		service = rs.FolgoryService
	case "lukki":
		service = rs.LukkiService
	case "hitbtc":
		service = rs.HitBTCService
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
		case "folgory":
			fallBackService = rs.FolgoryService
		case "lukki":
			fallBackService = rs.LukkiService
		case "hitbtc":
			fallBackService = rs.HitBTCService
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
func (rs *RateSevice) GetBtcRates() (map[string]models.RateV2, error) {
	rates := make(map[string]models.RateV2)
	if rs.FiatRates.LastUpdated.Unix()+UpdateFiatRatesTimeFrame < time.Now().Unix() {
		err := rs.LoadFiatRates()
		if err != nil {
			return nil, err
		}
	}
	if rs.BtcRates.LastUpdated+UpdateBtcRatesTimeFrame > time.Now().Unix() {
		return rs.BtcRates.Rates, nil
	}
	btcRate, err := rs.GetBtcEURRate()
	for code, rate := range rs.FiatRates.Rates {
		newRate := decimal.NewFromFloat(rate * btcRate)
		rate := models.RateV2{
			Name: models.FixerRatesNames[code],
			Rate: newRate,
		}
		rates[code] = rate
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

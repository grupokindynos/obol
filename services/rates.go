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
		rateDec := decimal.NewFromFloat(singleRate.Rate)
		rate := models.RateV2{
			Name: singleRate.Name,
		}
		var rateNum decimal.Decimal
		if coin.Info.Tag == "USDC" || coin.Info.Tag == "TUSD" || coin.Info.Tag == "USDT" {
			rateNum = rateDec.Div(orderPrice)
		} else {
			rateNum = orderPrice.Mul(rateDec)
		}
		if code == "BTC" {
			rate.Rate, _ = rateNum.Round(8).Float64()
		} else {
			rate.Rate, _ = rateNum.Round(6).Float64()
		}
		rates[code] = rate
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
		for code, rate := range coinRates {
			if code == "BTC" {
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
	var coinFromUSDRate decimal.Decimal
	for code, rate := range coinFromRates {
		if code == "USD" {
			coinFromUSDRate = decimal.NewFromFloat(rate.Rate)
		}
	}
	var coinToUSDRate decimal.Decimal
	for code, rate := range coinToRates {
		if code == "USD" {
			coinToUSDRate = decimal.NewFromFloat(rate.Rate)
		}
	}
	floatConvert, _ := coinFromUSDRate.DivRound(coinToUSDRate, 6).Float64()
	return floatConvert, nil
}

func (rs *RateSevice) GetCoinLiquidity(coin *coins.Coin) (float64, error) {
	coinWalls, err := rs.GetCoinOrdersWall(coin)
	if err != nil {
		return 0, err
	}
	orderWall := coinWalls["sell"]
	var liquidity decimal.Decimal
	for _, order := range orderWall {
		orderLiquidity := order.Amount.Mul(order.Price)
		liquidity = liquidity.Add(orderLiquidity)
	}
	btcRates, err := rs.GetBtcRates()
	if err != nil {
		return 0, err
	}
	var btcUSDRate decimal.Decimal
	for code, rate := range btcRates {
		if code == "USD" {
			btcUSDRate = decimal.NewFromFloat(rate.Rate)
		}
	}
	floatConvert, _ := liquidity.Mul(btcUSDRate).Round(8).Float64()
	return floatConvert, err
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
			coinWall[i].Amount = coinWall[i].Amount.Mul(coinWall[i].Price)
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

	var rates [][]decimal.Decimal
	var percentageSum decimal.Decimal
	for _, order := range coinWall {
		percentage := order.Amount.DivRound(amountRequested, 6)
		percentageSum = percentageSum.Add(percentage)
		var orderArr []decimal.Decimal
		if percentageSum.GreaterThan(decimal.NewFromInt(1)) {
			exceed := percentageSum.Sub(decimal.NewFromInt(1))
			rest := percentage.Sub(exceed)
			orderArr = []decimal.Decimal{order.Amount, order.Price, rest.Round(6)}
			percentageSum = percentageSum.Sub(exceed)
		} else {
			orderArr = []decimal.Decimal{order.Amount, order.Price, percentage}
		}
		rates = append(rates, orderArr)
		amountParsed = amountParsed.Sub(order.Amount)
		if amountParsed.LessThanOrEqual(decimal.Zero) {
			break
		}
	}
	amountParsed = amountRequested
	var AvrPrice decimal.Decimal
	for _, rateFloat := range rates {
		AvrPrice = AvrPrice.Add(rateFloat[1].Mul(rateFloat[2]))
		amountParsed = amountParsed.Sub(rateFloat[0])
		if amountParsed.LessThanOrEqual(decimal.Zero) {
			break
		}
	}
	var rate obol.CoinToCoinWithAmountResponse
	if coinTo.Info.Tag == "BTC" {
		rate.AveragePrice, _ = AvrPrice.Round(8).Float64()
	} else if coinFrom.Info.Tag == "BTC" {
		rate.AveragePrice, _ = decimal.NewFromInt(1).DivRound(AvrPrice, 8).Float64()
	} else {
		rateConv, err := rs.GetCoinToCoinRates(coinTo, btcData)
		if err != nil {
			return rate, err
		}
		rate.AveragePrice, _ = AvrPrice.DivRound(decimal.NewFromFloat(rateConv), 8).Float64()
	}
	amount := decimal.NewFromFloat(rate.AveragePrice)
	amount = amount.Mul(amountRequested)
	rate.Amount, _ = amount.Float64()
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
	for code, r := range rs.FiatRates.Rates {
		var rate models.RateV2
		if code == "BTC" {
			rate = models.RateV2{
				Name: models.FixerRatesNames[code],
				Rate: 1,
			}
		} else {
			newRate := decimal.NewFromFloat(r * btcRate)
			float, _ := newRate.Float64()
			rate = models.RateV2{
				Name: models.FixerRatesNames[code],
				Rate: float,
			}
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

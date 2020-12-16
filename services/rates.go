package services

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/grupokindynos/obol/services/exchanges/birake"
	"github.com/grupokindynos/obol/services/exchanges/bithumb"

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
	KuCoinService       *kucoin.Service
	GraviexService      *graviex.Service
	BitrueService       *bitrue.Service
	HitBTCService       *hitbtc.Service
	LukkiService        *lukki.Service
	BithumbService      *bithumb.Service
	BirakeService       *birake.Service
}

// GetCoinRates is the main function to get the rates of a coin using the OpenRates structure
func (rs *RateSevice) GetCoinRates(coin *coins.Coin, exchange string, buyWall bool) (map[string]models.RateV2, error) {

	rates := make(map[string]models.RateV2)
	btcRates, err := rs.GetBtcRates()
	if err != nil {
		return rates, err
	}
	if coin.Info.Tag == "BTC" {
		return btcRates, nil
	}
	// TODO Remove when exchanges update to FIRO Ticker
	if coin.Info.Tag == "FIRO" {
		coin.Info.Tag = "XZC"
	}
	ratesWall, err := rs.GetCoinOrdersWall(coin, exchange)
	if err != nil {
		return rates, err
	}
	var orders []models.MarketOrder
	var orderPrice decimal.Decimal
	if buyWall {
		orders = ratesWall["buy"]
		orderPrice = orders[len(orders)-1].Price
	} else {
		orders = ratesWall["sell"]
		orderPrice = orders[0].Price
	}
	if strings.ToUpper(coin.Info.Tag) == "" {
		// This handles coins with no BTC market, only stablecoin markets. Tested for USDT
		usdRates, err := rs.GetUsdRates()
		if err != nil {
			return rates, err
		}
		for code, singleRate := range usdRates {
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

			// temporal solution
			if coin.Rates.Exchange == "mock" {
				rate.Rate = 0
			}
			rates[code] = rate
		}
		return rates, err
	} else {
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

			// temporal solution
			if coin.Rates.Exchange == "mock" {
				rate.Rate = 0
			}
			rates[code] = rate
		}
		return rates, err
	}
}

// GetCoinToCoinRates will return the rates from a crypto to a crypto using the exchanges data
func (rs *RateSevice) GetCoinToCoinRates(coinFrom *coins.Coin, coinTo *coins.Coin, exchange string) (rate float64, err error) {
	if coinFrom.Info.Tag == coinTo.Info.Tag {
		return rate, config.ErrorNoC2CWithSameCoin
	}
	if coinTo.Info.Tag == "BTC" {
		coinRates, err := rs.GetCoinRates(coinFrom, exchange, false)
		for code, rate := range coinRates {
			if code == "BTC" {
				return rate.Rate, err
			}
		}
	}
	coinFromRates, err := rs.GetCoinRates(coinFrom, exchange, false)
	if err != nil {
		return 0, err
	}
	coinToRates, err := rs.GetCoinRates(coinTo, exchange, false)
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

func (rs *RateSevice) GetCoinLiquidity(coin *coins.Coin, exchange string) (float64, error) {
	coinWalls, err := rs.GetCoinOrdersWall(coin, exchange)
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
func (rs *RateSevice) GetCoinToCoinRatesWithAmount(coinFrom *coins.Coin, coinTo *coins.Coin, amountReq float64, exchange string) (obol.CoinToCoinWithAmountResponse, error) {
	if coinFrom.Info.Tag == coinTo.Info.Tag {
		return obol.CoinToCoinWithAmountResponse{}, config.ErrorNoC2CWithSameCoin
	}
	amountRequested := decimal.NewFromFloat(amountReq).Round(6)
	var coinWall []models.MarketOrder
	if coinFrom.Info.Tag == "BTC" {
		coinToWalls, err := rs.GetCoinOrdersWall(coinTo, exchange)
		if err != nil {
			return obol.CoinToCoinWithAmountResponse{}, err
		}
		coinWall = coinToWalls["buy"]
		for i := 0; i < len(coinWall); i++ {
			coinWall[i].Amount = coinWall[i].Amount.Mul(coinWall[i].Price)
		}
	} else {
		coinFromWalls, err := rs.GetCoinOrdersWall(coinFrom, exchange)
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
		// Calculates the percentage of the order's total amount an order from the coin wall represents and continues
		// looping through the coinWall until the added percentage is 100%. Stops when parsedAmount is zero or less.
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
	var avrPrice decimal.Decimal
	for _, rateFloat := range rates {
		avrPrice = avrPrice.Add(rateFloat[1].Mul(rateFloat[2]))
		amountParsed = amountParsed.Sub(rateFloat[0])
		if amountParsed.LessThanOrEqual(decimal.Zero) {
			break
		}
	}
	var rate obol.CoinToCoinWithAmountResponse
	if coinTo.Info.Tag == "BTC" {
		rate.AveragePrice, _ = avrPrice.Round(8).Float64()
	} else if coinFrom.Info.Tag == "BTC" {
		rate.AveragePrice, _ = decimal.NewFromInt(1).DivRound(avrPrice, 8).Float64()
	} else if coinFrom.Info.Token && coinFrom.Info.Tag != "ETH" {
		rateConv, err := rs.GetCoinToCoinRates(coinFrom, coinTo, exchange)
		if err != nil {
			return rate, err
		}
		rate.AveragePrice = rateConv
	} else {
		rateConv, err := rs.GetCoinToCoinRates(coinTo, btcData, exchange)
		if err != nil {
			return rate, err
		}
		rate.AveragePrice, _ = avrPrice.DivRound(decimal.NewFromFloat(rateConv), 8).Float64()
	}
	amount := decimal.NewFromFloat(rate.AveragePrice).Mul(amountRequested)
	rate.Amount, _ = amount.Float64()
	return rate, err
}

// GetCoinOrdersWall will return the buy/sell orders from selected or fallback exchange
func (rs *RateSevice) GetCoinOrdersWall(coin *coins.Coin, exchange string) (orders map[string][]models.MarketOrder, err error) {
	var service Exchange
	coinTag := coin.Info.Tag
	preferredExchange := ""
	if exchange != "" {
		preferredExchange = exchange
	} else {
		preferredExchange = coin.Rates.Exchange
	}

	// TODO Remove when exchanges update FIROs info
	if coin.Info.Tag == "FIRO" {
		coin.Info.Tag = "XZC"
	}
	switch preferredExchange {
	case "binance":
		service = rs.BinanceService
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
	case "stex":
		service = rs.StexService
	case "southxchange":
		service = rs.SouthXChangeService
	case "mock":
		service = rs.BinanceService
		coinTag = "ETH"
	case "bithumb":
		service = rs.BithumbService
	case "birake":
		service = rs.BirakeService
	}

	if service == nil {
		return nil, config.ErrorNoServiceForCoin
	}
	orders, err = service.CoinMarketOrders(coinTag)
	if err != nil {
		var fallBackService Exchange
		switch coin.Rates.FallBackExchange {
		case "binance":
			fallBackService = rs.BinanceService
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
		case "stex":
			fallBackService = rs.StexService
		case "southxchange":
			fallBackService = rs.SouthXChangeService
		case "birake":
			service = rs.BirakeService
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

// GetUSD-URRate will return the price of BTC on EUR
func (rs *RateSevice) GetUsdEurRate() (float64, error) {
	res, err := config.HttpClient.Get("https://bitstamp.net/api/v2/ticker/eurusd")
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

// TODO RATES FOR USDT GTH
// GetUsdRates will return the USD rates using the OpenRates structure
func (rs *RateSevice) GetUsdRates() (map[string]models.RateV2, error) {
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
	usdRate, err := rs.GetUsdEurRate()
	for code, r := range rs.FiatRates.Rates {
		var rate models.RateV2
		if code == "BTC" {
			rate = models.RateV2{
				Name: models.FixerRatesNames[code],
				Rate: 1,
			}
		} else {
			newRate := decimal.NewFromFloat(r * usdRate)
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

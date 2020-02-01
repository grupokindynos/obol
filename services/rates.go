package services

import (
	"encoding/json"
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
	if coinFrom.Info.Tag == "BTC" {
		coinRates, err := rs.GetCoinRates(coinTo, true)
		for _, rate := range coinRates {
			if rate.Code == "BTC" {
				return rate.Rate, err
			}
		}
	}
	if coinTo.Info.Tag == "BTC" {
		coinRates, err := rs.GetCoinRates(coinFrom, false)
		for _, rate := range coinRates {
			if rate.Code == "BTC" {
				return 1 / rate.Rate, err
			}
		}
	}
	if coinFrom.Info.Tag == coinTo.Info.Tag {
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
func (rs *RateSevice) GetCoinToCoinRatesWithAmount(coinFrom *coins.Coin, coinTo *coins.Coin, amountReq float64, wall string) (rate obol.CoinToCoinWithAmountResponse, err error) {
	if coinFrom.Info.Tag == coinTo.Info.Tag {
		return rate, config.ErrorNoC2CWithSameCoin
	}
	var coinMarkets map[string][]models.MarketOrder
	var coinRates []models.Rate
	// First get the orders wall from the coin we are converting
	if coinFrom.Info.Tag == "BTC" {
		coinMarkets, err = rs.GetCoinOrdersWall(coinTo)
		if err != nil {
			return obol.CoinToCoinWithAmountResponse{}, err
		}
		coinRates, err = rs.GetCoinRates(coinTo, true)
	} else {
		coinMarkets, err = rs.GetCoinOrdersWall(coinFrom)
		if err != nil {
			return obol.CoinToCoinWithAmountResponse{}, err
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
	var orderWalls string
	if wall != "" {
		orderWalls = wall
	} else {
		orderWalls = "sell"
	}
	// Init vars for loop
	var countedAmount float64
	var pricesSum amount.AmountType
	var orders []models.MarketOrder
	if coinFrom.Info.Tag == "BTC" {
		orders = coinMarkets["buy"]
	} else {
		orders = coinMarkets[orderWalls]
	}
	ratesResponse := obol.CoinToCoinWithAmountResponse{
		Rates:        [][]float64{},
		AveragePrice: 0,
	}
	// Looping against values on exchange to make an approachable rate based on the amount.
	for _, order := range orders {
		if countedAmount+order.Amount >= amountReq {
			diff := math.Abs((countedAmount + order.Amount) - amountReq)
			newAmount := order.Amount - diff
			countedAmount += newAmount
			percentage, err := amount.NewAmount(newAmount / amountReq)
			if err != nil {
				return obol.CoinToCoinWithAmountResponse{}, err
			}
			pricesSum += order.Price * percentage
			if coinFrom.Info.Tag == "BTC" {
				orderAmount, err := amount.NewAmount(order.Amount)
				if err != nil {
					return obol.CoinToCoinWithAmountResponse{}, err
				}
				var floatArray []float64
				floatArray = append(floatArray, orderAmount.ToNormalUnit(), order.Price.ToNormalUnit())
				ratesResponse.Rates = append(ratesResponse.Rates, floatArray)
			} else {
				orderAmount, err := amount.NewAmount(order.Amount)
				if err != nil {
					return obol.CoinToCoinWithAmountResponse{}, err
				}
				var floatArray []float64
				floatArray = append(floatArray, orderAmount.ToNormalUnit(), coinToBTCRate/order.Price.ToNormalUnit())
				ratesResponse.Rates = append(ratesResponse.Rates, floatArray)
			}
		} else {
			countedAmount += order.Amount
			percentage, err := amount.NewAmount(order.Amount / amountReq)
			if err != nil {
				return obol.CoinToCoinWithAmountResponse{}, err
			}
			pricesSum += order.Price * percentage
			if coinFrom.Info.Tag == "BTC" {
				orderAmount, err := amount.NewAmount(order.Amount)
				if err != nil {
					return obol.CoinToCoinWithAmountResponse{}, err
				}
				var floatArray []float64
				floatArray = append(floatArray, orderAmount.ToNormalUnit(), order.Price.ToNormalUnit())
				ratesResponse.Rates = append(ratesResponse.Rates, floatArray)
			} else {
				orderAmount, err := amount.NewAmount(order.Amount)
				if err != nil {
					return obol.CoinToCoinWithAmountResponse{}, err
				}
				var floatArray []float64
				floatArray = append(floatArray, orderAmount.ToNormalUnit(), coinToBTCRate/order.Price.ToNormalUnit())
				ratesResponse.Rates = append(ratesResponse.Rates, floatArray)
			}
		}
		if countedAmount >= amountReq {
			break
		}
	}
	var finalRate float64
	if coinFrom.Info.Tag == "BTC" {
		finalRate = pricesSum.ToNormalUnit()
	} else {
		finalRate = coinToBTCRate / pricesSum.ToNormalUnit()
	}
	ratesResponse.AveragePrice = finalRate
	return ratesResponse, err
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

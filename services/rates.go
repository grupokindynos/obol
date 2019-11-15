package services

import (
	"encoding/json"
	"github.com/grupokindynos/common/coin-factory/coins"
	"github.com/grupokindynos/common/obol"
	"github.com/grupokindynos/obol/config"
	"github.com/grupokindynos/obol/models"
	"github.com/grupokindynos/obol/models/exchanges"
	"github.com/grupokindynos/obol/services/exchanges/binance"
	"github.com/grupokindynos/obol/services/exchanges/bittrex"
	"github.com/grupokindynos/obol/services/exchanges/crex24"
	"github.com/grupokindynos/obol/services/exchanges/cryptobridge"
	"github.com/grupokindynos/obol/services/exchanges/graviex"
	"github.com/grupokindynos/obol/services/exchanges/kucoin"
	"github.com/grupokindynos/obol/services/exchanges/novaexchange"
	"github.com/grupokindynos/obol/services/exchanges/southxhcange"
	"github.com/grupokindynos/obol/services/exchanges/stex"
	"github.com/grupokindynos/olympus-utils/amount"
	"github.com/joho/godotenv"
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
	CryptoBridgeService *cryptobridge.Service
	Crex24Service       *crex24.Service
	StexService         *stex.Service
	SouthXChangeService *southxhcange.Service
	NovaExchangeService *novaexchange.Service
	KuCoinService       *kucoin.Service
	GraviexService      *graviex.Service
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
			if coin.Tag == "USDC" || coin.Tag == "TUSD" || coin.Tag == "USDT" {
				rate.Rate = math.Floor((singleRate.Rate/orders[0].Price)*1e8) / 1e8
			} else {
				rate.Rate = math.Floor((orders[0].Price*singleRate.Rate)*1e8) / 1e8
			}
		} else {
			if coin.Tag == "USDC" || coin.Tag == "TUSD" || coin.Tag == "USDT" {
				rate.Rate = math.Floor((orders[0].Price/singleRate.Rate)*10000) / 10000
			} else {
				rate.Rate = math.Floor((orders[0].Price*singleRate.Rate)*10000) / 10000
			}
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
				return 1 / rate.Rate, err
			}
		}
	}
	if coinTo.Tag == "BTC" {
		coinRates, err := rs.GetCoinRates(coinFrom, false)
		for _, rate := range coinRates {
			if rate.Code == "BTC" {
				return rate.Rate, err
			}
		}
	}
	if coinFrom.Tag == coinTo.Tag {
		return rate, config.ErrorNoC2CWithSameCoin
	}
	coinFromRates, err := rs.GetCoinRates(coinFrom, true)
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
	return coinFromCommonRate/coinToCommonRate, err
}

// GetCoinToCoinRatesWithAmount is used to get the rates from crypto to crypto using a specified amount to convert
func (rs *RateSevice) GetCoinToCoinRatesWithAmount(coinFrom *coins.Coin, coinTo *coins.Coin, amountReq float64, wall string) (rate obol.CoinToCoinWithAmountResponse, err error) {
	if coinFrom.Tag == coinTo.Tag {
		return rate, config.ErrorNoC2CWithSameCoin
	}
	var coinMarkets map[string][]models.MarketOrder
	var coinRates []models.Rate
	// First get the orders wall from the coin we are converting to
	if coinFrom.Tag == "BTC" {
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
	var pricesSum float64
	var orders []models.MarketOrder
	if coinFrom.Tag == "BTC" {
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
			percentage := newAmount / amountReq
			pricesSum += order.Price * percentage
			if coinFrom.Tag == "BTC" {
				orderAmount, err := amount.NewAmount(order.Amount)
				if err != nil {
					return obol.CoinToCoinWithAmountResponse{}, err
				}
				orderPrice, err := amount.NewAmount(order.Price)
				if err != nil {
					return obol.CoinToCoinWithAmountResponse{}, err
				}
				var floatArray []float64
				floatArray = append(floatArray, orderAmount.ToNormalUnit(), orderPrice.ToNormalUnit())
				ratesResponse.Rates = append(ratesResponse.Rates, floatArray)
			} else {
				orderAmount, err := amount.NewAmount(order.Amount)
				if err != nil {
					return obol.CoinToCoinWithAmountResponse{}, err
				}
				orderPrice, err := amount.NewAmount(coinToBTCRate / order.Price)
				if err != nil {
					return obol.CoinToCoinWithAmountResponse{}, err
				}
				var floatArray []float64
				floatArray = append(floatArray, orderAmount.ToNormalUnit(), orderPrice.ToNormalUnit())
				ratesResponse.Rates = append(ratesResponse.Rates, floatArray)
			}
		} else {
			countedAmount += order.Amount
			percentage := order.Amount / amountReq
			pricesSum += order.Price * percentage
			if coinFrom.Tag == "BTC" {
				orderAmount, err := amount.NewAmount(order.Amount)
				if err != nil {
					return obol.CoinToCoinWithAmountResponse{}, err
				}
				orderPrice, err := amount.NewAmount(order.Price)
				if err != nil {
					return obol.CoinToCoinWithAmountResponse{}, err
				}
				var floatArray []float64
				floatArray = append(floatArray, orderAmount.ToNormalUnit(), orderPrice.ToNormalUnit())
				ratesResponse.Rates = append(ratesResponse.Rates, floatArray)
			} else {
				orderAmount, err := amount.NewAmount(order.Amount)
				if err != nil {
					return obol.CoinToCoinWithAmountResponse{}, err
				}
				orderPrice, err := amount.NewAmount(coinToBTCRate / order.Price)
				if err != nil {
					return obol.CoinToCoinWithAmountResponse{}, err
				}
				var floatArray []float64
				floatArray = append(floatArray, orderAmount.ToNormalUnit(), orderPrice.ToNormalUnit())
				ratesResponse.Rates = append(ratesResponse.Rates, floatArray)
			}
		}
		if countedAmount >= amountReq {
			break
		}
	}
	var finalRate float64
	if coinFrom.Tag == "BTC" {
		finalRate = pricesSum
	} else {
		finalRate = coinToBTCRate / pricesSum
	}
	finalRateHand, err := amount.NewAmount(finalRate)
	if err != nil {
		return obol.CoinToCoinWithAmountResponse{}, err
	}
	ratesResponse.AveragePrice = finalRateHand.ToNormalUnit()
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
	case "cryptobridge":
		service = rs.CryptoBridgeService
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
	if rs.FiatRates.LastUpdated.Unix()+UpdateFiatRatesTimeFrame < time.Now().Unix() {
		err = rs.LoadFiatRates()
		if err != nil {
			return nil, err
		}
	}
	if rs.BtcRates.LastUpdated+UpdateBtcRatesTimeFrame > time.Now().Unix() {
		return rs.BtcRates.Rates, nil
	}
	btcMxnRate, err := rs.GetBtcMxnRate()
	newRate := btcMxnRate / rs.FiatRates.Rates["MXN"]
	for key, rate := range rs.FiatRates.Rates {
		rate := models.Rate{
			Code: key,
			Name: models.FixerRatesNames[key],
			Rate: rate * newRate,
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

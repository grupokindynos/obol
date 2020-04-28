package stex

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"

	"github.com/grupokindynos/obol/config"
	"github.com/grupokindynos/obol/models"
	"github.com/grupokindynos/obol/models/exchanges"
	"github.com/shopspring/decimal"
)

// Service is a common structure for a exchange
type Service struct {
	BaseRateURL   string
	MarketRateURL string
	TickerID      map[string]string
}

// CoinMarketOrders is used to get the market sell and buy wall from a coin
func (s *Service) CoinMarketOrders(coin string) (orders map[string][]models.MarketOrder, err error) {
	orders = make(map[string][]models.MarketOrder)
	// Instead of using the ticker, this one uses an ID
	// A map is created on the Init Service with known coins and ticker ID for this exchange.
	// First get the ID
	value, exist := s.TickerID[strings.ToUpper(coin)]
	if !exist {
		return orders, config.ErrorUnknownIdForCoin
	}
	res, err := config.HttpClient.Get(s.MarketRateURL + value)
	if err != nil {
		return orders, config.ErrorRequestTimeout
	}
	defer func() {
		_ = res.Body.Close()
	}()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return orders, config.ErrorRequestTimeout
	}
	var Response exchanges.StexMarkets
	err = json.Unmarshal(contents, &Response)
	if err != nil {
		return orders, config.ErrorRequestTimeout
	}
	// For DAPS we use the ETH pair, this means we need to convert prices to BTC factor.
	// First we get the ETH to BTC price.
	var ethPrice float64
	if value == "819" {
		price, err := s.GetETHPrice()
		if err != nil {
			return nil, err
		}
		ethPrice, err = strconv.ParseFloat(price, 64)
		if err != nil {
			return nil, err
		}
	}
	var buyOrders []models.MarketOrder
	var sellOrders []models.MarketOrder
	for _, order := range Response.Data.Ask {
		var price float64
		if value == "819" {
			price, err = strconv.ParseFloat(order.Price, 64)
			if err != nil {
				return nil, err
			}
			price = price * ethPrice
		} else {
			price, err = strconv.ParseFloat(order.Price, 64)
			if err != nil {
				return nil, err
			}
		}
		am, err := strconv.ParseFloat(order.Amount, 64)
		if err != nil {
			return nil, err
		}
		newOrder := models.MarketOrder{
			Price:  decimal.NewFromFloat(price),
			Amount: decimal.NewFromFloat(am),
		}
		sellOrders = append(sellOrders, newOrder)
	}
	for _, order := range Response.Data.Bid {
		var price float64
		if value == "819" {
			price, err = strconv.ParseFloat(order.Price, 64)
			if err != nil {
				return nil, err
			}
			price = price * ethPrice
		} else {
			price, err = strconv.ParseFloat(order.Price, 64)
			if err != nil {
				return nil, err
			}
		}
		am, _ := strconv.ParseFloat(order.Amount, 64)
		newOrder := models.MarketOrder{
			Price:  decimal.NewFromFloat(price),
			Amount: decimal.NewFromFloat(am),
		}
		buyOrders = append(buyOrders, newOrder)
	}
	sort.Slice(buyOrders, func(i, j int) bool {
		return buyOrders[i].Price.GreaterThan(buyOrders[j].Price)
	})
	sort.Slice(sellOrders, func(i, j int) bool {
		return sellOrders[i].Price.LessThan(sellOrders[j].Price)
	})
	orders["buy"] = buyOrders
	orders["sell"] = sellOrders
	return orders, err
}

func (s *Service) GetETHPrice() (string, error) {
	res, err := config.HttpClient.Get(s.BaseRateURL)
	if err != nil {
		return "", config.ErrorRequestTimeout
	}
	defer func() {
		_ = res.Body.Close()
	}()
	contents, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", config.ErrorRequestTimeout
	}
	var response exchanges.StexTickers
	err = json.Unmarshal(contents, &response)
	if err != nil {
		return "", config.ErrorRequestTimeout
	}
	for _, ticker := range response.Data {
		if ticker.Symbol == "ETH_BTC" {
			return ticker.Last, nil
		}
	}
	return "", errors.New("no information")
}

// InitService is used to safely start a new service reference.
func InitService() *Service {
	tickerID := make(map[string]string)

	// Populate with known ID and Tickers
	tickerID["XSG"] = "250"
	tickerID["DIVI"] = "1119"
	tickerID["DAPS"] = "819"

	s := &Service{
		BaseRateURL:   "https://api3.stex.com/public/ticker/",
		MarketRateURL: "https://api3.stex.com/public/orderbook/",
		TickerID:      tickerID,
	}
	return s
}

package stex

import (
	"encoding/json"
	"github.com/grupokindynos/obol/config"
	"github.com/grupokindynos/obol/models"
	"github.com/grupokindynos/obol/models/exchanges"
	"io/ioutil"
	"strconv"
	"strings"
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
	var Response exchanges.StexMarkets
	err = json.Unmarshal(contents, &Response)
	var buyOrders []models.MarketOrder
	var sellOrders []models.MarketOrder
	for _, order := range Response.Data.Ask {
		price, _ := strconv.ParseFloat(order.Price, 64)
		amount, _ := strconv.ParseFloat(order.Amount, 64)
		newOrder := models.MarketOrder{
			Price:  price,
			Amount: amount,
		}
		sellOrders = append(sellOrders, newOrder)
	}
	for _, order := range Response.Data.Bid {
		price, _ := strconv.ParseFloat(order.Price, 64)
		amount, _ := strconv.ParseFloat(order.Amount, 64)
		newOrder := models.MarketOrder{
			Price:  price,
			Amount: amount,
		}
		buyOrders = append(buyOrders, newOrder)
	}
	orders["buy"] = buyOrders
	orders["sell"] = sellOrders
	return orders, err
}

// InitService is used to safely start a new service reference.
func InitService() *Service {
	tickerID := make(map[string]string)

	// Populate with known ID and Tickers
	tickerID["XSG"] = "250"

	s := &Service{
		BaseRateURL:   "https://api3.stex.com/public/ticker/",
		MarketRateURL: "https://api3.stex.com/public/orderbook/",
		TickerID:      tickerID,
	}
	return s
}

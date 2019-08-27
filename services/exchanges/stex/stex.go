package stex

import (
	"encoding/json"
	"github.com/grupokindynos/obol/config"
	"github.com/grupokindynos/obol/models"
	"github.com/grupokindynos/obol/models/exchanges"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Service struct {
	BaseRateURL   string
	MarketRateURL string
	TickerID      map[string]string
}

func (s Service) CoinRate(coin string) (rate float64, err error) {
	// Instead of using the ticker, this one uses an ID
	// A map is created on the Init Service with known coins and ticker ID for this exchange.
	// First get the ID
	value, exist := s.TickerID[strings.ToUpper(coin)]
	if !exist {
		return rate, config.ErrorUnknownIdForCoin
	}
	res, err := http.Get(s.BaseRateURL + value)
	if err != nil {
		return rate, err
	} else {
		defer func() {
			_ = res.Body.Close()
		}()
		contents, err := ioutil.ReadAll(res.Body)
		var Response exchanges.StexRate
		err = json.Unmarshal(contents, &Response)
		if err != nil {
			return rate, err
		}
		rate, err := strconv.ParseFloat(Response.Data.Last, 64)
		return rate, err
	}
}

// CoinMarketOrders is used to get the market sell and buy wall from a coin
func (s *Service) CoinMarketOrders(coin string) (orders []models.MarketOrder, err error) {
	// Instead of using the ticker, this one uses an ID
	// A map is created on the Init Service with known coins and ticker ID for this exchange.
	// First get the ID
	value, exist := s.TickerID[strings.ToUpper(coin)]
	if !exist {
		return orders, config.ErrorUnknownIdForCoin
	}
	res, err := http.Get(s.MarketRateURL + value)
	if err != nil {
		return orders, err
	} else {
		defer func() {
			_ = res.Body.Close()
		}()
		contents, err := ioutil.ReadAll(res.Body)
		var Response exchanges.StexMarkets
		err = json.Unmarshal(contents, &Response)
		for _, ask := range Response.Data.Ask {
			price, _ := strconv.ParseFloat(ask.Price, 64)
			amount, _ := strconv.ParseFloat(ask.Amount, 64)
			newOrder := models.MarketOrder{
				Price:  price,
				Amount: amount,
			}
			orders = append(orders, newOrder)
		}
		return orders, err
	}
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

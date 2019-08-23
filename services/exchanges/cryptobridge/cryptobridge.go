package cryptobridge

import (
	"encoding/json"
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
}

func (s Service) CoinRate(coin string) (rate float64, err error) {
	res, err := http.Get(s.BaseRateURL + strings.ToUpper(coin) + "_BTC")
	if err != nil {
		return rate, err
	} else {
		defer func() {
			_ = res.Body.Close()
		}()
		contents, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return rate, err
		}
		var Response exchanges.CryptoBridgeRate
		err = json.Unmarshal(contents, &Response)
		if err != nil {
			return rate, err
		}
		rate, err := strconv.ParseFloat(Response.Last, 64)
		return rate, err
	}
}

func (s *Service) CoinMarketOrders(coin string) (orders []models.MarketOrder, err error) {
	res, err := http.Get(s.MarketRateURL + strings.ToUpper(coin) + "_BTC")
	if err != nil {
		return orders, err
	} else {
		defer func() {
			_ = res.Body.Close()
		}()
		contents, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return orders, err
		}
		var Response exchanges.CryptoBridgeMarkets
		err = json.Unmarshal(contents, &Response)
		if err != nil {
			return orders, err
		}
		for _, ask := range Response.Bids {
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

func InitService() *Service {
	s := &Service{
		BaseRateURL:   "https://api.crypto-bridge.org/v1/ticker/",
		MarketRateURL: "https://api.crypto-bridge.org/v2/market/book/",
	}
	return s
}

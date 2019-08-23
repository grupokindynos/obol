package binance

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

func (s *Service) CoinRate(coin string) (rate float64, err error) {
	res, err := http.Get(s.BaseRateURL + strings.ToUpper(coin) + "BTC")
	if err != nil {
		return rate, err
	} else {
		defer func() {
			_ = res.Body.Close()
		}()
		contents, err := ioutil.ReadAll(res.Body)
		var Response exchanges.BinanceRate
		err = json.Unmarshal(contents, &Response)
		rate, err := strconv.ParseFloat(Response.LastPrice, 64)
		return rate, err
	}
}

func (s *Service) CoinMarketOrders(coin string) (orders []models.MarketOrder, err error) {
	res, err := http.Get(s.MarketRateURL + strings.ToUpper(coin) + "BTC")
	if err != nil {
		return orders, err
	} else {
		defer func() {
			_ = res.Body.Close()
		}()
		contents, err := ioutil.ReadAll(res.Body)
		var Response exchanges.BinanceMarkets
		err = json.Unmarshal(contents, &Response)
		for _, ask := range Response.Asks {
			price, _ := strconv.ParseFloat(ask[0], 64)
			amount, _ := strconv.ParseFloat(ask[1], 64)
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
		BaseRateURL:   "https://api.binance.com/api/v1/ticker/24hr?symbol=",
		MarketRateURL: "https://api.binance.com/api/v1/depth?symbol=",
	}
	return s
}

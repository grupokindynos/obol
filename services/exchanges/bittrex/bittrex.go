package bittrex

import (
	"encoding/json"
	"github.com/grupokindynos/obol/models"
	"github.com/grupokindynos/obol/models/exchanges"
	"io/ioutil"
	"net/http"
	"strings"
)

// Service is a common structure for a exchange
type Service struct {
	BaseRateURL   string
	MarketRateURL string
}

func (s *Service) CoinRate(coin string) (rate float64, err error) {
	res, err := http.Get(s.BaseRateURL + strings.ToUpper(coin))
	if err != nil {
		return rate, err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	contents, err := ioutil.ReadAll(res.Body)
	var Response exchanges.BittrexRate
	err = json.Unmarshal(contents, &Response)
	return Response.Result.Last, err
}

// CoinMarketOrders is used to get the market sell and buy wall from a coin
func (s *Service) CoinMarketOrders(coin string) (orders []models.MarketOrder, err error) {
	res, err := http.Get(s.MarketRateURL + strings.ToUpper(coin) + "&type=both")
	if err != nil {
		return orders, err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	contents, err := ioutil.ReadAll(res.Body)
	var Response exchanges.BittrexMarkets
	err = json.Unmarshal(contents, &Response)
	for _, ask := range Response.Result.Buy {
		price := ask.Rate
		amount := ask.Quantity
		newOrder := models.MarketOrder{
			Price:  price,
			Amount: amount,
		}
		orders = append(orders, newOrder)
	}
	return orders, err
}

// InitService is used to safely start a new service reference.
func InitService() *Service {
	s := &Service{
		BaseRateURL:   "https://api.bittrex.com/api/v1.1/public/getticker?market=BTC-",
		MarketRateURL: "https://api.bittrex.com/api/v1.1/public/getorderbook?market=BTC-",
	}
	return s
}

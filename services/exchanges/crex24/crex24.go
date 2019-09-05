package crex24

import (
	"encoding/json"
	"github.com/grupokindynos/obol/config"
	"github.com/grupokindynos/obol/models"
	"github.com/grupokindynos/obol/models/exchanges"
	"io/ioutil"
	"strings"
)

// Service is a common structure for a exchange
type Service struct {
	MarketRateURL string
}

// CoinMarketOrders is used to get the market sell and buy wall from a coin
func (s *Service) CoinMarketOrders(coin string) (orders []models.MarketOrder, err error) {
	res, err := config.HttpClient.Get(s.MarketRateURL + strings.ToUpper(coin) + "-BTC")
	if err != nil {
		return orders, config.ErrorRequestTimeout
	}
	defer func() {
		_ = res.Body.Close()
	}()
	contents, err := ioutil.ReadAll(res.Body)
	var Response exchanges.Crex24Markets
	err = json.Unmarshal(contents, &Response)
	for _, ask := range Response.BuyLevels {
		newOrder := models.MarketOrder{
			Price:  ask.Price,
			Amount: ask.Volume,
		}
		orders = append(orders, newOrder)
	}
	return orders, err
}

// InitService is used to safely start a new service reference.
func InitService() *Service {
	s := &Service{
		MarketRateURL: "https://api.crex24.com/v2/public/orderBook?instrument=",
	}
	return s
}

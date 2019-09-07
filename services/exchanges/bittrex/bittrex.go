package bittrex

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
func (s *Service) CoinMarketOrders(coin string) (orders map[string][]models.MarketOrder, err error) {
	orders = make(map[string][]models.MarketOrder)
	res, err := config.HttpClient.Get(s.MarketRateURL + "BTC-" + strings.ToUpper(coin) + "&type=both")
	if err != nil {
		return orders, config.ErrorRequestTimeout
	}
	defer func() {
		_ = res.Body.Close()
	}()
	contents, err := ioutil.ReadAll(res.Body)
	var Response exchanges.BittrexMarkets
	err = json.Unmarshal(contents, &Response)
	var buyOrders []models.MarketOrder
	var sellOrders []models.MarketOrder
	for _, order := range Response.Result.Sell {
		price := order.Rate
		amount := order.Quantity
		newOrder := models.MarketOrder{
			Price:  price,
			Amount: amount,
		}
		sellOrders = append(sellOrders, newOrder)
	}
	for _, order := range Response.Result.Buy {
		price := order.Rate
		amount := order.Quantity
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
	s := &Service{
		MarketRateURL: "https://api.bittrex.com/api/v1.1/public/getorderbook?market=",
	}
	return s
}

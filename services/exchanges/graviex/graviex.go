package graviex

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
	MarketRateURL string
}

// CoinMarketOrders is used to get the market sell and buy wall from a coin
func (s *Service) CoinMarketOrders(coin string) (orders map[string][]models.MarketOrder, err error) {
	orders = make(map[string][]models.MarketOrder)
	res, err := config.HttpClient.Get(s.MarketRateURL + strings.ToLower(coin) + "btc")
	if err != nil {
		return orders, config.ErrorRequestTimeout
	}
	defer func() {
		_ = res.Body.Close()
	}()
	contents, err := ioutil.ReadAll(res.Body)
	var Response exchanges.GraviexMarkets
	err = json.Unmarshal(contents, &Response)
	var buyOrders []models.MarketOrder
	var sellOrders []models.MarketOrder
	for _, order := range Response.Asks {
		price, _ := strconv.ParseFloat(order.Price, 64)
		amount, _ := strconv.ParseFloat(order.Volume, 64)
		newOrder := models.MarketOrder{
			Price:  price,
			Amount: amount,
		}
		sellOrders = append(sellOrders, newOrder)
	}
	for _, order := range Response.Bids {
		price, _ := strconv.ParseFloat(order.Price, 64)
		amount, _ := strconv.ParseFloat(order.Volume, 64)
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
		MarketRateURL: "https://graviex.net/api/v3/order_book?market=",
	}
	return s
}

package probit

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
	res, err := config.HttpClient.Get(s.MarketRateURL + strings.ToUpper(coin) + "_BTC")
	if err != nil {
		return orders, config.ErrorRequestTimeout
	}
	defer func() {
		_ = res.Body.Close()
	}()
	contents, err := ioutil.ReadAll(res.Body)
	var Response exchanges.ProbitMarkets
	err = json.Unmarshal(contents, &Response)
	var buyOrders []models.MarketOrder
	var sellOrders []models.MarketOrder
	for _, order := range Response.Data {
		price, _ := strconv.ParseFloat(order.Price, 64)
		amount, _ := strconv.ParseFloat(order.Quantity, 64)
		newOrder := models.MarketOrder{
			Price:  price,
			Amount: amount,
		}
		if order.Side == "buy" {
			buyOrders = append(buyOrders, newOrder)
		} else if order.Side == "sell" {
			sellOrders = append(sellOrders, newOrder)
		}
	}
	orders["buy"] = buyOrders
	orders["sell"] = sellOrders
	return orders, err
}

// InitService is used to safely start a new service reference.
func InitService() *Service {
	s := &Service{
		MarketRateURL: "https://api.probit.com/api/exchange/v1/order_book?market_id=",
	}
	return s
}

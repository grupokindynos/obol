package kucoin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/grupokindynos/obol/config"
	"github.com/grupokindynos/obol/models"
	"github.com/grupokindynos/obol/models/exchanges"
	"github.com/shopspring/decimal"
)

// Service is a common structure for a exchange
type Service struct {
	MarketRateURL string
}

// CoinMarketOrders is used to get the market sell and buy wall from a coin
func (s *Service) CoinMarketOrders(coin string) (orders map[string][]models.MarketOrder, err error) {
	orders = make(map[string][]models.MarketOrder)
	res, err := config.HttpClient.Get(s.MarketRateURL + strings.ToUpper(coin) + "-BTC")
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
	var Response exchanges.KuCoinMarkets
	err = json.Unmarshal(contents, &Response)
	if err != nil {
		return orders, config.ErrorRequestTimeout
	}
	fmt.Println(Response)
	var buyOrders []models.MarketOrder
	var sellOrders []models.MarketOrder
	for _, order := range Response.Data.Asks {
		price, err := strconv.ParseFloat(order[0], 64)
		if err != nil {
			return orders, config.ErrorRequestTimeout
		}
		am, err := strconv.ParseFloat(order[1], 64)
		if err != nil {
			return orders, config.ErrorRequestTimeout
		}
		newOrder := models.MarketOrder{
			Price:  decimal.NewFromFloat(price),
			Amount: decimal.NewFromFloat(am),
		}
		sellOrders = append(sellOrders, newOrder)
	}
	for _, order := range Response.Data.Bids {
		price, err := strconv.ParseFloat(order[0], 64)
		if err != nil {
			return orders, config.ErrorRequestTimeout
		}
		am, err := strconv.ParseFloat(order[1], 64)
		if err != nil {
			return orders, config.ErrorRequestTimeout
		}
		newOrder := models.MarketOrder{
			Price:  decimal.NewFromFloat(price),
			Amount: decimal.NewFromFloat(am),
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
		MarketRateURL: "https://api.kucoin.com/api/v1/market/orderbook/level2_100?symbol=",
	}
	return s
}

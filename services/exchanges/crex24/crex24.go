package crex24

import (
	"encoding/json"
	"github.com/grupokindynos/obol/config"
	"github.com/grupokindynos/obol/models"
	"github.com/grupokindynos/obol/models/exchanges"
	"github.com/olympus-protocol/ogen/utils/amount"
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
	var buyOrders []models.MarketOrder
	var sellOrders []models.MarketOrder
	for _, order := range Response.SellLevels {
		priceConv, err := amount.NewAmount(order.Price)
		if err != nil {
			return nil, err
		}
		newOrder := models.MarketOrder{
			Price:  priceConv,
			Amount: order.Volume,
		}
		sellOrders = append(sellOrders, newOrder)
	}
	for _, order := range Response.BuyLevels {
		priceConv, err := amount.NewAmount(order.Price)
		if err != nil {
			return nil, err
		}
		newOrder := models.MarketOrder{
			Price:  priceConv,
			Amount: order.Volume,
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
		MarketRateURL: "https://api.crex24.com/v2/public/orderBook?instrument=",
	}
	return s
}

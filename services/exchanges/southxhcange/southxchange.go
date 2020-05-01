package southxhcange

import (
	"encoding/json"
	"io/ioutil"
	"sort"
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
	res, err := config.HttpClient.Get(s.MarketRateURL + strings.ToUpper(coin) + "/BTC")
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
	var Response exchanges.SouthXChangeMarkets
	err = json.Unmarshal(contents, &Response)
	if err != nil {
		return orders, config.ErrorRequestTimeout
	}
	var buyOrders []models.MarketOrder
	var sellOrders []models.MarketOrder
	for _, order := range Response.BuyOrders {
		newOrder := models.MarketOrder{
			Price:  decimal.NewFromFloat(order.Price),
			Amount: decimal.NewFromFloat(order.Amount),
		}
		sellOrders = append(sellOrders, newOrder)
	}
	for _, order := range Response.SellOrders {
		newOrder := models.MarketOrder{
			Price:  decimal.NewFromFloat(order.Price),
			Amount: decimal.NewFromFloat(order.Amount),
		}
		buyOrders = append(buyOrders, newOrder)
	}
	sort.Slice(buyOrders, func(i, j int) bool {
		return buyOrders[i].Price.LessThan(buyOrders[j].Price)
	})
	sort.Slice(sellOrders, func(i, j int) bool {
		return sellOrders[i].Price.GreaterThan(sellOrders[j].Price)
	})
	orders["buy"] = buyOrders
	orders["sell"] = sellOrders
	return orders, err
}

// InitService is used to safely start a new service reference.
func InitService() *Service {
	s := &Service{
		MarketRateURL: "https://www.southxchange.com/api/book/",
	}
	return s
}

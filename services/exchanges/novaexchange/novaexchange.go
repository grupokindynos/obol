package novaexchange

import (
	"encoding/json"
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
	res, err := config.HttpClient.Get(s.MarketRateURL + "BTC_" + strings.ToUpper(coin))
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
	var Response exchanges.NovaExchangeMarkets
	err = json.Unmarshal(contents, &Response)
	if err != nil {
		return orders, config.ErrorRequestTimeout
	}
	var buyOrders []models.MarketOrder
	var sellOrders []models.MarketOrder
	for _, order := range Response.Items {
		if order.Tradetype == "BUY" {
			price, err := strconv.ParseFloat(order.Price, 64)
			if err != nil {
				return orders, config.ErrorRequestTimeout
			}
			am, err := strconv.ParseFloat(order.Amount, 64)
			if err != nil {
				return orders, config.ErrorRequestTimeout
			}
			newOrder := models.MarketOrder{
				Price:  decimal.NewFromFloat(price),
				Amount: decimal.NewFromFloat(am),
			}
			buyOrders = append(buyOrders, newOrder)
		} else if order.Tradetype == "SELL" {
			price, err := strconv.ParseFloat(order.Price, 64)
			if err != nil {
				return orders, config.ErrorRequestTimeout
			}
			am, err := strconv.ParseFloat(order.Amount, 64)
			if err != nil {
				return orders, config.ErrorRequestTimeout
			}
			newOrder := models.MarketOrder{
				Price:  decimal.NewFromFloat(price),
				Amount: decimal.NewFromFloat(am),
			}
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
		MarketRateURL: "https://novaexchange.com/remote/v2/market/orderhistory/",
	}
	return s
}

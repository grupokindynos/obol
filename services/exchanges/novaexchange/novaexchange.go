package novaexchange

import (
	"encoding/json"
	"github.com/grupokindynos/obol/config"
	"github.com/grupokindynos/obol/models"
	"github.com/grupokindynos/obol/models/exchanges"
	"github.com/olympus-protocol/ogen/utils/amount"
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
	res, err := config.HttpClient.Get(s.MarketRateURL + "BTC_" + strings.ToUpper(coin))
	if err != nil {
		return orders, config.ErrorRequestTimeout
	}
	defer func() {
		_ = res.Body.Close()
	}()
	contents, err := ioutil.ReadAll(res.Body)
	var Response exchanges.NovaExchangeMarkets
	err = json.Unmarshal(contents, &Response)
	var buyOrders []models.MarketOrder
	var sellOrders []models.MarketOrder
	for _, order := range Response.Items {
		if order.Tradetype == "BUY" {
			price, _ := strconv.ParseFloat(order.Price, 64)
			priceConv, err := amount.NewAmount(price)
			if err != nil {
				return nil, err
			}
			am, _ := strconv.ParseFloat(order.Amount, 64)
			newOrder := models.MarketOrder{
				Price:  priceConv,
				Amount: am,
			}
			buyOrders = append(buyOrders, newOrder)
		} else if order.Tradetype == "SELL" {
			price, _ := strconv.ParseFloat(order.Price, 64)
			priceConv, err := amount.NewAmount(price)
			if err != nil {
				return nil, err
			}
			am, _ := strconv.ParseFloat(order.Amount, 64)
			newOrder := models.MarketOrder{
				Price:  priceConv,
				Amount: am,
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

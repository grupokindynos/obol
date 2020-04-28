package southxhcange

import (
	"encoding/json"
	"io/ioutil"
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
		price := order.Price
		am := order.Amount
		priceConv := decimal.NewFromFloat(price)
		newOrder := models.MarketOrder{
			Price:  priceConv,
			Amount: am,
		}
		sellOrders = append(sellOrders, newOrder)
	}
	for _, order := range Response.SellOrders {
		price := order.Price
		priceConv := decimal.NewFromFloat(price)
		am := order.Amount
		newOrder := models.MarketOrder{
			Price:  priceConv,
			Amount: am,
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
		MarketRateURL: "https://www.southxchange.com/api/book/",
	}
	return s
}

package lukki

import (
	"encoding/json"
	"io/ioutil"
	"sort"
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
	marketStr := strings.ToLower(coin) + "_btc"
	res, err := config.HttpClient.Get(s.MarketRateURL + marketStr)
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
	var Response exchanges.LukkiMarkets
	err = json.Unmarshal(contents, &Response)
	if err != nil {
		return orders, config.ErrorRequestTimeout
	}
	var buyOrders []models.MarketOrder
	var sellOrders []models.MarketOrder
	for _, order := range Response.Data {
		if order.Direction == 0 {
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
		} else {
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
	// Since lukki send the information unordered we should order the walls.
	// For buy order is from lowest to biggest and for sell orders is from biggest to lowest based on the price.

	sort.Slice(buyOrders, func(i, j int) bool {
		return buyOrders[i].Price.GreaterThan(buyOrders[j].Price)
	})
	sort.Slice(sellOrders, func(i, j int) bool {
		return sellOrders[i].Price.LessThan(sellOrders[j].Price)
	})
	orders["buy"] = buyOrders
	orders["sell"] = sellOrders
	return orders, err
}

// InitService is used to safely start a new service reference.
func InitService() *Service {
	s := &Service{
		MarketRateURL: "https://tva.lukki.io/trading/books?&ticker=",
	}
	return s
}

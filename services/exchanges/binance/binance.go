package binance

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
	var marketStr string
	if coin == "TUSD" || coin == "USDC" || coin == "USDT" {
		marketStr = "BTC" + strings.ToUpper(coin)
	} else {
		marketStr = strings.ToUpper(coin) + "BTC"
	}
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
	var Response exchanges.BinanceMarkets
	err = json.Unmarshal(contents, &Response)
	if err != nil {
		return orders, config.ErrorRequestTimeout
	}
	var buyOrders []models.MarketOrder
	var sellOrders []models.MarketOrder
	for _, order := range Response.Asks {
		price, err := strconv.ParseFloat(order[0], 64)
		if err != nil {
			return orders, config.ErrorRequestTimeout
		}
		priceConv := decimal.NewFromFloat(price)
		am, err := strconv.ParseFloat(order[1], 64)
		if err != nil {
			return orders, config.ErrorRequestTimeout
		}
		newOrder := models.MarketOrder{
			Price:  priceConv,
			Amount: am,
		}
		sellOrders = append(sellOrders, newOrder)
	}
	for _, order := range Response.Bids {
		price, err := strconv.ParseFloat(order[0], 64)
		if err != nil {
			return orders, config.ErrorRequestTimeout
		}
		priceConv := decimal.NewFromFloat(price)
		if err != nil {
			return nil, err
		}
		am, err := strconv.ParseFloat(order[1], 64)
		if err != nil {
			return orders, config.ErrorRequestTimeout
		}
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
		MarketRateURL: "https://api.binance.com/api/v1/depth?symbol=",
	}
	return s
}

package bitrue

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
	res, err := config.HttpClient.Get(s.MarketRateURL + strings.ToUpper(coin) + "BTC")
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
	var Response exchanges.BitrueMarkets
	err = json.Unmarshal(contents, &Response)
	if err != nil {
		return orders, config.ErrorRequestTimeout
	}
	var buyOrders []models.MarketOrder
	var sellOrders []models.MarketOrder
	for _, order := range Response.Asks {
		strMarhshalPrice, err := json.Marshal(order[0])
		if err != nil {
			return nil, err
		}
		strMarhshalAmount, err := json.Marshal(order[1])
		if err != nil {
			return nil, err
		}
		var strPrice, strAmount string
		err = json.Unmarshal(strMarhshalPrice, &strPrice)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(strMarhshalAmount, &strAmount)
		if err != nil {
			return nil, err
		}
		price, err := strconv.ParseFloat(strPrice, 64)
		if err != nil {
			return orders, config.ErrorRequestTimeout
		}
		priceConv := decimal.NewFromFloat(price)
		if err != nil {
			return nil, err
		}
		am, err := strconv.ParseFloat(strAmount, 64)
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
		strMarhshalPrice, err := json.Marshal(order[0])
		if err != nil {
			return nil, err
		}
		strMarhshalAmount, err := json.Marshal(order[1])
		if err != nil {
			return nil, err
		}
		var strPrice, strAmount string
		err = json.Unmarshal(strMarhshalPrice, &strPrice)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(strMarhshalAmount, &strAmount)
		if err != nil {
			return nil, err
		}
		price, err := strconv.ParseFloat(strPrice, 64)
		if err != nil {
			return orders, config.ErrorRequestTimeout
		}
		priceConv := decimal.NewFromFloat(price)
		am, err := strconv.ParseFloat(strAmount, 64)
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
		MarketRateURL: "https://www.bitrue.com/api/v1/depth?symbol=",
	}
	return s
}

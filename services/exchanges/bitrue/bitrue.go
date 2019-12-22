package bitrue

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
	res, err := config.HttpClient.Get(s.MarketRateURL + strings.ToUpper(coin) + "BTC")
	if err != nil {
		return orders, config.ErrorRequestTimeout
	}
	defer func() {
		_ = res.Body.Close()
	}()
	contents, err := ioutil.ReadAll(res.Body)
	var Response exchanges.BitrueMarkets
	err = json.Unmarshal(contents, &Response)
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
		price, _ := strconv.ParseFloat(strPrice, 64)
		amount, _ := strconv.ParseFloat(strAmount, 64)
		newOrder := models.MarketOrder{
			Price:  price,
			Amount: amount,
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
		price, _ := strconv.ParseFloat(strPrice, 64)
		amount, _ := strconv.ParseFloat(strAmount, 64)
		newOrder := models.MarketOrder{
			Price:  price,
			Amount: amount,
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

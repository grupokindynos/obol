package bithumb

import (
	"encoding/json"
	"github.com/grupokindynos/obol/config"
	"github.com/grupokindynos/obol/models"
	"github.com/grupokindynos/obol/models/exchanges"
	exchanges2 "github.com/grupokindynos/obol/services/exchanges"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"sort"
	"strconv"
)

// Service is a common structure for a exchange
type Service struct {
	MarketRateURL string
	coinGecko *exchanges2.CoinGecko
}

// CoinMarketOrders is used to get the market sell and buy wall from a coin
func (s *Service) CoinMarketOrders(coin string) (orders map[string][]models.MarketOrder, err error) {
	return s.coinGecko.GetSimplePriceToBtcAsRate(coin)
	/*orders, err = s.getMarketInfo(strings.ToUpper(coin) + "-USDT")
	if err != nil {
		return orders, err
	}
	return orders, err*/
}

func (s *Service) getMarketInfo(market string) (orders map[string][]models.MarketOrder, err error){
	// Retrieves BTC price as Bithumb markets for GTH have no GTH/BTC market
	orders = make(map[string][]models.MarketOrder)
	res, err := config.HttpClient.Get(s.MarketRateURL + "?symbol=" + market)
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
	var Response exchanges.BithumbMarkets
	err = json.Unmarshal(contents, &Response)
	if err != nil {
		return orders, config.ErrorRequestTimeout
	}
	var buyOrders []models.MarketOrder
	var sellOrders []models.MarketOrder
	for _, order := range Response.Data.Asks {
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
		am, err := strconv.ParseFloat(strAmount, 64)
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
		am, err := strconv.ParseFloat(strAmount, 64)
		if err != nil {
			return orders, config.ErrorRequestTimeout
		}
		newOrder := models.MarketOrder{
			Price:  decimal.NewFromFloat(price),
			Amount: decimal.NewFromFloat(am),
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
	return orders, nil
}

// InitService is used to safely start a new service reference.
func InitService() *Service {
	s := &Service{
		MarketRateURL: "https://global-openapi.bithumb.pro/openapi/v1/spot/orderBook",
		coinGecko: exchanges2.NewCoinGecko(),
	}
	return s
}

package crex24

import (
	"encoding/json"
	"github.com/grupokindynos/obol/models"
	"github.com/grupokindynos/obol/models/exchanges"
	"io/ioutil"
	"net/http"
	"strings"
)

type Service struct {
	BaseRateURL   string
	MarketRateURL string
}

func (s Service) CoinRate(coin string) (rate float64, err error) {
	res, err := http.Get(s.BaseRateURL + "[NamePairs=BTC_" + strings.ToUpper(coin) + "]")
	if err != nil {
		return rate, err
	} else {
		defer func() {
			_ = res.Body.Close()
		}()
		contents, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return rate, err
		}
		var Response exchanges.Crex24Rates
		err = json.Unmarshal(contents, &Response)
		if err != nil {
			return rate, err
		}
		return Response.Tickers[0].Last, err
	}
}

func (s *Service) CoinMarketOrders(coin string) (orders []models.MarketOrder, err error) {
	res, err := http.Get(s.MarketRateURL + strings.ToUpper(coin) + "-BTC")
	if err != nil {
		return orders, err
	} else {
		defer func() {
			_ = res.Body.Close()
		}()
		contents, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return orders, err
		}
		var Response exchanges.Crex24Markets
		err = json.Unmarshal(contents, &Response)
		if err != nil {
			return orders, err
		}
		for _, ask := range Response.BuyLevels {
			newOrder := models.MarketOrder{
				Price:  ask.Price,
				Amount: ask.Volume,
			}
			orders = append(orders, newOrder)
		}
		return orders, err
	}
}

func InitService() *Service {
	s := &Service{
		BaseRateURL:   "https://api.crex24.com/CryptoExchangeService/BotPublic/ReturnTicker?request=",
		MarketRateURL: "https://api.crex24.com/v2/public/orderBook?instrument=",
	}
	return s
}

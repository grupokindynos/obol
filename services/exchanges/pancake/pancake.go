package pancake

import (
	"bytes"
	"encoding/json"
	"github.com/grupokindynos/obol/config"
	"github.com/grupokindynos/obol/models"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"strings"
)

// Service is a common structure for a exchange
type Service struct {
	MarketRateURL string
}

type priceParams struct {
	Amount float64 `json:"amount"`
	Asset string `json:"asset"`
}

// CoinMarketOrders is used to get the market sell and buy wall from a coin
func (s *Service) CoinMarketOrders(coin string) (orders map[string][]models.MarketOrder, err error) {
	orders = make(map[string][]models.MarketOrder)
	buyOrders := []models.MarketOrder{
		{
			Price: decimal.NewFromFloat(1),
			Amount: decimal.NewFromFloat(1000),
		},
	}
	sellOrders := []models.MarketOrder{
		{
			Price: decimal.NewFromFloat(1),
			Amount: decimal.NewFromFloat(1000),
		},
	}

	if strings.ToUpper(coin) == "BUSD" {
		// TODO Return hardcoded orders
		orders["buy"] = buyOrders
		orders["sell"] = sellOrders
		return orders, err
	}
	priceData, _ := json.Marshal(priceParams{
		Amount: 1000,
		Asset:  coin,
	})


	res, err := config.HttpClient.Post(s.MarketRateURL, "application/json", bytes.NewReader(priceData))
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
	var Response BSCPriceResponse
	err = json.Unmarshal(contents, &Response)
	if err != nil {
		return orders, config.ErrorRequestTimeout
	}

	buyOrders = []models.MarketOrder{
		{
			Price: decimal.NewFromFloat(Response.Data.AveragePrice),
			Amount: decimal.NewFromFloat(1000),
		},
	}
	sellOrders = []models.MarketOrder{
		{
			Price: decimal.NewFromFloat(Response.Data.AveragePrice),
			Amount: decimal.NewFromFloat(1000),
		},
	}

	orders["buy"] = buyOrders
	orders["sell"] = sellOrders
	return orders, err
}

// InitService is used to safely start a new service reference.
func InitService() *Service {
	s := &Service{
		MarketRateURL: "https://pp-bsc-api.herokuapp.com/api/v1/price",
	}
	return s
}

type BSCBase struct {
	Status int32 `json:"status"`
	Error string `json:"error"`
}

type BSCPriceResponse struct {
	BSCBase
	Data struct{
		AveragePrice float64 `json:"average_price"`
		ReceivedAmount float64 `json:"received_amount"`
	} `json:"data"`
}


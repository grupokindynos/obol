package exchanges

import (
	"encoding/json"
	"fmt"
	coinfactory "github.com/grupokindynos/common/coin-factory"
	"github.com/grupokindynos/obol/models"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type CoinGecko struct {
	baseUrl string
}

func NewCoinGecko() *CoinGecko {
	c := new(CoinGecko)
	c.baseUrl =  "https://api.coingecko.com/api/v3/simple/price"
	return c
}

/*
	Meant to be used as a provisional way of retrieving fees for coins as their BTC values. The only case it should be
	used indefinitely is if a listing does not require PolisPay services (a.k.a. Shift and Gift Cards).
 */
func (c *CoinGecko) GetSimplePriceToBtcAsRate(asset string) (orders map[string][]models.MarketOrder, err error){
	orders = make(map[string][]models.MarketOrder)
	coinfInfo, err := coinfactory.GetCoin(strings.ToUpper(asset))
	if err != nil {
		return
	}
	rateUrl := fmt.Sprintf("%s?ids=%s&vs_currencies=btc,usd", c.baseUrl, coinfInfo.Rates.CoinGeckoId)

	res, err := http.Get(rateUrl)
	if err != nil {
		return
	}

	var price map[string]Price

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	defer res.Body.Close()
	err = json.Unmarshal(data, &price)
	if err != nil {
		log.Println("problem unmarshalling: ", err)
		log.Println("original data: ", string(data))
		return
	}
	orderSample := models.MarketOrder{
		Amount: decimal.NewFromFloat(100),
		Price:  price[coinfInfo.Rates.CoinGeckoId].Btc,
	}
	orders["buy"] = append(orders["buy"], orderSample)
	orders["sell"] = append(orders["sell"], orderSample)
	fmt.Println(orders)
	return
}

type Price struct {
	Usd decimal.Decimal `json:"usd"`
	Btc decimal.Decimal `json:"btc"`
}
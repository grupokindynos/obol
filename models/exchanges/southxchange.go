package exchanges

// SouthXChangeMarkets is the response of the market depth query on SouthXChange Exchange
type SouthXChangeMarkets struct {
	BuyOrders []struct {
		Index  int     `json:"Index"`
		Amount float64 `json:"Amount"`
		Price  float64 `json:"Price"`
	} `json:"BuyOrders"`
	SellOrders []struct {
		Index  int     `json:"Index"`
		Amount float64 `json:"Amount"`
		Price  float64 `json:"Price"`
	} `json:"SellOrders"`
}

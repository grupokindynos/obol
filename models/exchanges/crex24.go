package exchanges

// Crex24Markets is the response of the market depth query on Crex24 Exchange
type Crex24Markets struct {
	BuyLevels []struct {
		Price  float64 `json:"price"`
		Volume float64 `json:"volume"`
	} `json:"buyLevels"`
	SellLevels []struct {
		Price  float64 `json:"price"`
		Volume float64 `json:"volume"`
	} `json:"sellLevels"`
}

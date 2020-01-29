package exchanges

// ProbitMarkets is the response of the market depth query on Binance Exchange
type ProbitMarkets struct {
	Data []struct {
		Side     string `json:"side"`
		Price    string `json:"price"`
		Quantity string `json:"quantity"`
	} `json:"data"`
}

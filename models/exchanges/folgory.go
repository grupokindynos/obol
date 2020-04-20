package exchanges

// FolgoryMarkets is the response of the market depth query on Binance Exchange
type FolgoryMarkets struct {
	Message string `json:"message"`
	Data    struct {
		Timestamp int        `json:"timestamp"`
		Asks      [][]string `json:"asks"`
		Bids      [][]string `json:"bids"`
	} `json:"data"`
}

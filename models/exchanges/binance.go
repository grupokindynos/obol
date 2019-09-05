package exchanges

// BinanceMarkets is the response of the market depth query on Binance Exchange
type BinanceMarkets struct {
	LastUpdateID int        `json:"lastUpdateId"`
	Bids         [][]string `json:"bids"`
	Asks         [][]string `json:"asks"`
}

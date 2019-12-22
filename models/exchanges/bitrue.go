package exchanges

// BitrueMarkets is the response of the market depth query on Binance Exchange
type BitrueMarkets struct {
	LastUpdateID int             `json:"lastUpdateId"`
	Bids         [][]interface{} `json:"bids"`
	Asks         [][]interface{} `json:"asks"`
}

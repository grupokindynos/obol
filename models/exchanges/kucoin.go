package exchanges

// KuCoinMarkets is the response of the market depth query on KuCoin Exchange
type KuCoinMarkets struct {
	Code string `json:"code"`
	Data struct {
		Sequence string     `json:"sequence"`
		Asks     [][]string `json:"asks"`
		Bids     [][]string `json:"bids"`
		Time     int64      `json:"time"`
	} `json:"data"`
}

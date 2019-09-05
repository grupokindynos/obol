package exchanges

// BittrexMarkets is the response of the market depth query on Bittrex Exchange
type BittrexMarkets struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Result  struct {
		Buy []struct {
			Quantity float64 `json:"Quantity"`
			Rate     float64 `json:"Rate"`
		} `json:"buy"`
		Sell []struct {
			Quantity float64 `json:"Quantity"`
			Rate     float64 `json:"Rate"`
		} `json:"sell"`
	} `json:"result"`
}

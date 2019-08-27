package exchanges

// BittrexRate is the response of the rates query on Bittrex Exchange
type BittrexRate struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Result  struct {
		Bid  float64 `json:"Bid"`
		Ask  float64 `json:"Ask"`
		Last float64 `json:"Last"`
	} `json:"result"`
}

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

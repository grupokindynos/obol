package exchanges

// CryptoBridgeRate is the response of the rates query on CryptoBridge Exchange
type CryptoBridgeRate struct {
	ID            string  `json:"id"`
	Last          string  `json:"last"`
	Volume        string  `json:"volume"`
	Ask           string  `json:"ask"`
	Bid           string  `json:"bid"`
	PercentChange float64 `json:"percentChange"`
}

// CryptoBridgeMarkets is the response of the market depth query on CryptoBridge Exchange
type CryptoBridgeMarkets struct {
	Bids []struct {
		Price  string `json:"price"`
		Amount string `json:"amount"`
	} `json:"bids"`
	Asks []struct {
		Price  string `json:"price"`
		Amount string `json:"amount"`
	} `json:"asks"`
}

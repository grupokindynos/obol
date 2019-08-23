package exchanges

type CryptoBridgeRate struct {
	ID            string  `json:"id"`
	Last          string  `json:"last"`
	Volume        string  `json:"volume"`
	Ask           string  `json:"ask"`
	Bid           string  `json:"bid"`
	PercentChange float64 `json:"percentChange"`
}

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

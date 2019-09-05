package exchanges

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

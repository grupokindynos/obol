package exchanges

// GraviexMarkets is the response of the market depth query on Graviex Exchange
type GraviexMarkets struct {
	Timestamp int        `json:"timestamp"`
	Asks      [][]string `json:"asks"`
	Bids      [][]string `json:"bids"`
}

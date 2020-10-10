package exchanges

type BithumbMarkets struct {
	Data      BithumbMarketData `json:"data"`
	Code      string            `json:"code"`
	Msg       string            `json:"msg"`
	Timestamp int64             `json:"timestamp"`
	StartTime interface{}       `json:"startTime"`
}

type BithumbMarketData struct {
	Symbol string     `json:"symbol"`
	Bids   [][]string `json:"b"`
	Ver    string     `json:"ver"`
	Asks   [][]string `json:"s"`
}
